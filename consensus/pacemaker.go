// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"time"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/hotstuff"
	"github.com/wooyang2018/ppov-blockchain/logger"
)

type pacemaker struct {
	resources *Resources
	config    Config
	hotstuff  *hotstuff.Hotstuff

	state       *state
	voterState  *voterState
	leaderState *leaderState

	stopCh chan struct{}
}

func (pm *pacemaker) start() {
	if pm.stopCh != nil {
		return
	}
	pm.stopCh = make(chan struct{})
	go pm.batchRun()
	go pm.run()
	logger.I().Info("started pacemaker")
}

func (pm *pacemaker) stop() {
	if pm.stopCh == nil {
		return // not started yet
	}
	select {
	case <-pm.stopCh: // already stopped
		return
	default:
	}
	close(pm.stopCh)
	logger.I().Info("stopped pacemaker")
	pm.stopCh = nil
}

func (pm *pacemaker) batchRun() {
	subQC := pm.hotstuff.SubscribeNewQCHigh()
	defer subQC.Unsubscribe()

	for {
		pm.newBatch()

		select {
		case <-pm.stopCh:
			return
		case <-subQC.Events():
		}
	}
}

func (pm *pacemaker) run() {
	subQC := pm.hotstuff.SubscribeNewQCHigh()
	defer subQC.Unsubscribe()

	for {
		blkDelayT := pm.nextBlockDelay()
		pm.newBlock()
		beatT := pm.nextProposeTimeout()

		select {
		case <-pm.stopCh:
			return

		// either beatdelay timeout or I'm able to create qc
		case <-beatT.C:
		case <-subQC.Events():
		}
		beatT.Stop()

		select {
		case <-pm.stopCh:
			return
		case <-blkDelayT.C:
		}
		blkDelayT.Stop()
	}
}

func (pm *pacemaker) nextBlockDelay() *time.Timer {
	delay := pm.config.BlockDelay
	if VoteBatchFlag {
		if pm.leaderState.getBatchReadyNum() < pm.config.BlockBatchLimit {
			delay += time.Duration(pm.config.BlockBatchLimit-
				pm.leaderState.getBatchReadyNum()) * (pm.config.TxWaitTime / 2)
		}
	} else {
		if pm.voterState.getBatchNum() < pm.config.BlockBatchLimit {
			delay += time.Duration(pm.config.BlockBatchLimit-
				pm.voterState.getBatchNum()) * (pm.config.TxWaitTime / 2)
		}
	}
	return time.NewTimer(delay)
}

func (pm *pacemaker) nextProposeTimeout() *time.Timer {
	proposeWait := pm.config.ProposeTimeout
	return time.NewTimer(proposeWait)
}

func (pm *pacemaker) newBlock() {
	pm.state.mtxUpdate.Lock()
	defer pm.state.mtxUpdate.Unlock()

	select {
	case <-pm.stopCh:
		return
	default:
	}

	if !pm.state.isThisNodeLeader() {
		return
	}

	blk := pm.hotstuff.OnPropose()
	logger.I().Debugw("proposed block", "height", blk.Height(),
		"qc", qcRefHeight(blk.Justify()), "txs", len(blk.Transactions()))
	vote := blk.(*hsBlock).block.ProposerVote()
	pm.hotstuff.OnReceiveVote(newHsVote(vote, pm.state))
	pm.hotstuff.Update(blk)
}

func (pm *pacemaker) newBatch() {
	pm.state.mtxUpdate.Lock()
	defer pm.state.mtxUpdate.Unlock()

	select {
	case <-pm.stopCh:
		return
	default:
	}

	if !pm.state.isThisNodeWorker() {
		return
	}

	var txs []*core.Transaction
	if PreserveTxFlag {
		txs = pm.resources.TxPool.GetTxsFromQueue(pm.config.BatchTxLimit)
	} else {
		txs = pm.resources.TxPool.PopTxsFromQueue(pm.config.BatchTxLimit)
	}
	if len(txs) == 0 { // 忽略打包空Batch
		return
	}

	signer := pm.resources.Signer
	batch := core.NewBatch().SetTransactions(txs).SetTimestamp(time.Now().UnixNano()).Sign(signer)
	if VoteBatchFlag {
		if pm.state.isThisNodeVoter() {
			pm.voterState.addBatch(batch.Header())
		}
		pm.resources.MsgSvc.BroadcastBatch(batch)
	} else {
		if pm.state.isThisNodeLeader() {
			pm.voterState.addBatch(batch.Header())
		}
		leader := pm.resources.VldStore.GetWorker(pm.state.getLeaderIndex())
		pm.resources.MsgSvc.SendBatch(leader, batch)
	}

	widx := pm.resources.VldStore.GetWorkerIndex(signer.PublicKey())
	logger.I().Debugw("generated batch", "worker", widx,
		"txs", len(batch.Header().Transactions()))
}
