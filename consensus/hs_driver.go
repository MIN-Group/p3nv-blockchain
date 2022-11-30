// Copyright (C) 2021 Aung Maw
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"time"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/hotstuff"
	"github.com/wooyang2018/ppov-blockchain/logger"
	"github.com/wooyang2018/ppov-blockchain/storage"
)

type hsDriver struct {
	resources *Resources
	config    Config

	checkTxDelay time.Duration //检测TxPool交易数量的延迟

	state *state

	leaderState *leaderState
}

// 验证hsDriver实现了hotstuff的Driver
var _ hotstuff.Driver = (*hsDriver)(nil)

func (hsd *hsDriver) MajorityValidatorCount() int {
	return hsd.resources.VldStore.MajorityValidatorCount()
}

func (hsd *hsDriver) CreateLeaf(parent hotstuff.Block, qc hotstuff.QC, height uint64) hotstuff.Block {
	batchs := hsd.leaderState.popReadyBatch()
	txs := hsd.getBatchTxs(batchs)
	//core.Block的链式调用
	blk := core.NewBlock().
		SetParentHash(parent.(*hsBlock).block.Hash()).
		SetQuorumCert(qc.(*hsQC).qc).
		SetHeight(height).
		SetBatchs(batchs).
		SetTransactions(txs).
		SetExecHeight(hsd.resources.Storage.GetBlockHeight()).
		SetMerkleRoot(hsd.resources.Storage.GetMerkleRoot()).
		SetTimestamp(time.Now().UnixNano()).
		Sign(hsd.resources.Signer)
	if err := hsd.resources.TxPool.StoreTxs(blk.TxList()); err != nil {
		logger.I().Fatalf("store transactions of block error: %+v", err)
	}
	hsd.state.setBlock(blk)
	return newHsBlock(blk, hsd.state)
}

func (hsd *hsDriver) getBatchTxs(val []*core.Batch) [][]byte {
	txSet := make(map[string]struct{})
	txList := make([][]byte, 0)
	for _, batch := range val {
		for _, hash := range batch.Transactions() {
			if _, ok := txSet[string(hash)]; ok {
				continue //重复交易则跳过
			}
			if hsd.resources.Storage.HasTx(hash) {
				continue //已提交交易则跳过
			}
			txSet[string(hash)] = struct{}{} //集合去重
			txList = append(txList, hash)
		}
	}
	return txList
}

func (hsd *hsDriver) CreateQC(hsVotes []hotstuff.Vote) hotstuff.QC {
	votes := make([]*core.Vote, len(hsVotes))
	for i, hsv := range hsVotes {
		votes[i] = hsv.(*hsVote).vote
	}
	qc := core.NewQuorumCert().Build(votes)
	return newHsQC(qc, hsd.state)
}

func (hsd *hsDriver) BroadcastProposal(hsBlk hotstuff.Block) {
	blk := hsBlk.(*hsBlock).block
	hsd.resources.MsgSvc.BroadcastProposal(blk)
}

func (hsd *hsDriver) VoteBlock(hsBlk hotstuff.Block) {
	blk := hsBlk.(*hsBlock).block
	vote := blk.Vote(hsd.resources.Signer)
	hsd.resources.TxPool.SetTxsPending(blk.Transactions())
	hsd.delayVoteWhenNoTxs()
	proposer := hsd.resources.VldStore.GetWorkerIndex(blk.Proposer())
	if proposer != hsd.state.getLeaderIndex() {
		return // view changed happened
	}
	hsd.resources.MsgSvc.SendVote(blk.Proposer(), vote)
	logger.I().Debugw("voted block",
		"proposer", proposer,
		"height", hsBlk.Height(),
		"qc", qcRefHeight(hsBlk.Justify()),
	)
}

func (hsd *hsDriver) delayVoteWhenNoTxs() {
	timer := time.NewTimer(hsd.config.TxWaitTime)
	defer timer.Stop()
	for hsd.resources.TxPool.GetStatus().Total == 0 {
		select {
		case <-timer.C:
			return
		case <-time.After(hsd.checkTxDelay):
		}
	}
}

func (hsd *hsDriver) Commit(hsBlk hotstuff.Block) {
	bexe := hsBlk.(*hsBlock).block
	start := time.Now()
	txs, old := hsd.resources.TxPool.GetTxsToExecute(bexe.Transactions())
	logger.I().Debugw("commiting block", "height", bexe.Height(), "txs", len(txs))
	bcm, txcs := hsd.resources.Execution.Execute(bexe, txs)
	bcm.SetOldBlockTxs(old)
	data := &storage.CommitData{
		Block:        bexe,
		QC:           hsd.state.getQC(bexe.Hash()),
		Transactions: txs,
		BlockCommit:  bcm,
		TxCommits:    txcs,
	}
	err := hsd.resources.Storage.Commit(data)
	if err != nil {
		logger.I().Fatalf("commit storage error: %+v", err)
	}
	hsd.state.addCommitedTxCount(len(txs))
	hsd.cleanStateOnCommited(bexe)
	logger.I().Debugw("commited bock",
		"height", bexe.Height(),
		"batchs", len(bexe.Batchs()),
		"txs", len(txs),
		"elapsed", time.Since(start))
}

func (hsd *hsDriver) cleanStateOnCommited(bexec *core.Block) {
	// qc for bexe is no longer needed here after commited to storage
	hsd.state.deleteQC(bexec.Hash())
	hsd.resources.TxPool.RemoveTxs(bexec.Transactions())
	hsd.state.setCommitedBlock(bexec)

	folks := hsd.state.getUncommitedOlderBlocks(bexec)
	for _, blk := range folks {
		// put transactions from folked block back to queue
		hsd.resources.TxPool.PutTxsToQueue(blk.Transactions())
		hsd.state.deleteBlock(blk.Hash())
		hsd.state.deleteQC(blk.Hash())
	}
	hsd.deleteCommitedOlderBlocks(bexec)
}

func (hsd *hsDriver) deleteCommitedOlderBlocks(bexec *core.Block) {
	height := int64(bexec.Height()) - 20
	if height < 0 {
		return
	}
	blks := hsd.state.getOlderBlocks(uint64(height))
	for _, blk := range blks {
		hsd.state.deleteBlock(blk.Hash())
		hsd.state.deleteCommited(blk.Hash())
	}
}
