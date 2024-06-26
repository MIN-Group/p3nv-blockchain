// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/logger"
	"github.com/wooyang2018/ppov-blockchain/storage"
)

type genesis struct {
	resources *Resources
	chainID   int64

	done chan struct{}

	// collect votes from all validators instead of majority for genesis block
	votes map[string]*core.Vote

	mtxVote    sync.Mutex
	mtxNewView sync.Mutex

	b0 *core.Block
	q0 *core.QuorumCert

	mtxB0 sync.RWMutex
	mtxQ0 sync.RWMutex
}

func (gns *genesis) run() (*core.Block, *core.QuorumCert) {
	logger.I().Infow("creating genesis block...")
	gns.done = make(chan struct{})

	go gns.proposalLoop()
	go gns.voteLoop()
	go gns.newViewLoop()
	gns.propose()

	<-gns.done
	logger.I().Info("got genesis block and qc")
	gns.commit()
	return gns.getB0(), gns.getQ0()
}

func (gns *genesis) commit() {
	data := &storage.CommitData{
		Block: gns.getB0(),
		QC:    gns.getQ0(),
	}
	data.BlockCommit = core.NewBlockCommit().SetHash(data.Block.Hash())
	err := gns.resources.Storage.Commit(data)
	if err != nil {
		logger.I().Fatalf("commit storage error: %+v", err)
	}
	logger.I().Debugw("committed genesis bock")
}

func (gns *genesis) propose() {
	if !gns.isLeader(gns.resources.Signer.PublicKey()) {
		return
	}
	gns.votes = make(map[string]*core.Vote, gns.resources.VldStore.MajorityValidatorCount())
	b0 := gns.createGenesisBlock()
	gns.setB0(b0)
	logger.I().Infow("created genesis block, broadcasting...")
	go gns.broadcastProposalLoop()
	gns.onReceiveVote(b0.ProposerVote())
}

func (gns *genesis) createGenesisBlock() *core.Block {
	return core.NewBlock().
		SetHeight(0).
		SetParentHash(hashChainID(gns.chainID)).
		SetTimestamp(time.Now().UnixNano()).
		Sign(gns.resources.Signer)
}

func (gns *genesis) broadcastProposalLoop() {
	for {
		select {
		case <-gns.done:
			return
		default:
		}
		if gns.getQ0() == nil {
			if err := gns.resources.MsgSvc.BroadcastProposal(gns.getB0()); err != nil {
				logger.I().Errorw("broadcast proposal failed", "error", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (gns *genesis) isLeader(pubKey *core.PublicKey) bool {
	if !gns.resources.VldStore.IsWorker(pubKey) {
		return false
	}
	return gns.resources.VldStore.GetWorkerIndex(pubKey) == 0
}

func hashChainID(chainID int64) []byte {
	h := sha3.New256()
	binary.Write(h, binary.BigEndian, chainID)
	return h.Sum(nil)
}

func (gns *genesis) proposalLoop() {
	sub := gns.resources.MsgSvc.SubscribeProposal(1)
	defer sub.Unsubscribe()

	for {
		select {
		case <-gns.done:
			return

		case e := <-sub.Events():
			if err := gns.onReceiveProposal(e.(*core.Block)); err != nil {
				logger.I().Warnf("receive proposal failed, %+v", err.Error())
			}
		}
	}
}

func (gns *genesis) voteLoop() {
	sub := gns.resources.MsgSvc.SubscribeVote(10)
	defer sub.Unsubscribe()

	for {
		select {
		case <-gns.done:
			return

		case e := <-sub.Events():
			if err := gns.onReceiveVote(e.(*core.Vote)); err != nil {
				logger.I().Warnf("receive vote failed, %+v", err.Error())
			}
		}
	}
}

func (gns *genesis) newViewLoop() {
	sub := gns.resources.MsgSvc.SubscribeNewView(1)
	defer sub.Unsubscribe()

	for {
		select {
		case <-gns.done:
			return

		case e := <-sub.Events():
			if err := gns.onReceiveNewView(e.(*core.QuorumCert)); err != nil {
				logger.I().Warnf("receive new view failed, %+v", err.Error())
			}
		}
	}
}

func (gns *genesis) onReceiveProposal(proposal *core.Block) error {
	if err := proposal.Validate(gns.resources.VldStore); err != nil {
		return err
	}
	if !proposal.IsGenesis() {
		logger.I().Info("left behind, fetching genesis block...")
		return gns.fetchGenesisBlockAndQC(proposal.Proposer())
	}
	if !bytes.Equal(hashChainID(gns.chainID), proposal.ParentHash()) {
		return fmt.Errorf("different chain id genesis")
	}
	if !gns.isLeader(proposal.Proposer()) {
		return fmt.Errorf("proposer is not leader")
	}
	if len(proposal.Transactions()) != 0 {
		return fmt.Errorf("genesis block with txs")
	}
	gns.setB0(proposal)
	logger.I().Infow("got genesis block, voting...", "proposer", gns.resources.VldStore.GetWorkerIndex(proposal.Proposer()))
	return gns.resources.MsgSvc.SendVote(proposal.Proposer(), proposal.Vote(gns.resources.Signer))
}

func (gns *genesis) fetchGenesisBlockAndQC(peer *core.PublicKey) error {
	b0, err := gns.requestBlockByHeight(peer, 0)
	if err != nil {
		return err
	}
	if !b0.IsGenesis() {
		return fmt.Errorf("not genesis block")
	}
	b1, err := gns.requestBlockByHeight(peer, 1)
	if err != nil {
		return err
	}
	if !bytes.Equal(b0.Hash(), b1.QuorumCert().BlockHash()) {
		return fmt.Errorf("b1 qc ref is not b0")
	}
	gns.setB0(b0)
	gns.setQ0(b1.QuorumCert())
	close(gns.done)
	return nil
}

func (gns *genesis) requestBlockByHeight(peer *core.PublicKey, height uint64) (*core.Block, error) {
	blk, err := gns.resources.MsgSvc.RequestBlockByHeight(peer, height)
	if err != nil {
		return nil, fmt.Errorf("cannot get block by height %d, %w", height, err)
	}
	if err := blk.Validate(gns.resources.VldStore); err != nil {
		return nil, fmt.Errorf("validate block %d error %w", height, err)
	}
	return blk, nil
}

func (gns *genesis) onReceiveVote(vote *core.Vote) error {
	if gns.votes == nil {
		return errors.New("not accepting votes")
	}
	if err := vote.Validate(gns.resources.VldStore); err != nil {
		return err
	}
	gns.acceptVote(vote)
	return nil
}

func (gns *genesis) acceptVote(vote *core.Vote) {
	gns.mtxVote.Lock()
	defer gns.mtxVote.Unlock()

	gns.votes[vote.Voter().String()] = vote
	if len(gns.votes) < gns.resources.VldStore.MajorityValidatorCount() {
		return
	}
	vlist := make([]*core.Vote, 0, len(gns.votes))
	for _, vote := range gns.votes {
		vlist = append(vlist, vote)
	}
	gns.setQ0(core.NewQuorumCert().Build(vlist))
	logger.I().Infow("created qc, broadcasting...")
	gns.broadcastQC()
}

func (gns *genesis) broadcastQC() {
	for {
		select {
		case <-gns.done:
			return
		default:
		}
		if err := gns.resources.MsgSvc.BroadcastNewView(gns.getQ0()); err != nil {
			logger.I().Errorw("broadcast proposal failed", "error", err)
		}
		time.Sleep(time.Second)
	}
}

func (gns *genesis) onReceiveNewView(qc *core.QuorumCert) error {
	if err := qc.Validate(gns.resources.VldStore); err != nil {
		return err
	}
	b0 := gns.getB0()
	if b0 == nil {
		return fmt.Errorf("no received genesis block yet")
	}
	if !bytes.Equal(b0.Hash(), qc.BlockHash()) {
		return fmt.Errorf("invalid qc reference")
	}
	gns.acceptQC(qc)
	return nil
}

func (gns *genesis) acceptQC(qc *core.QuorumCert) {
	gns.mtxNewView.Lock()
	defer gns.mtxNewView.Unlock()

	select {
	case <-gns.done: // already done genesis
		return
	default:
	}
	b0 := gns.getB0()
	gns.setQ0(qc)
	if !gns.isLeader(gns.resources.Signer.PublicKey()) {
		gns.resources.MsgSvc.SendNewView(b0.Proposer(), qc)
	}
	close(gns.done) // when qc is accepted, genesis creation is done
}

func (gns *genesis) setB0(val *core.Block) {
	gns.mtxB0.Lock()
	defer gns.mtxB0.Unlock()
	gns.b0 = val
}

func (gns *genesis) setQ0(val *core.QuorumCert) {
	gns.mtxQ0.Lock()
	defer gns.mtxQ0.Unlock()
	gns.q0 = val
}

func (gns *genesis) getB0() *core.Block {
	gns.mtxB0.RLock()
	defer gns.mtxB0.RUnlock()
	return gns.b0
}

func (gns *genesis) getQ0() *core.QuorumCert {
	gns.mtxQ0.RLock()
	defer gns.mtxQ0.RUnlock()
	return gns.q0
}
