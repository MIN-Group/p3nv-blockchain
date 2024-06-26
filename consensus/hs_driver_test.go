// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/hotstuff"
	"github.com/wooyang2018/ppov-blockchain/storage"
	"github.com/wooyang2018/ppov-blockchain/txpool"
)

func setupTestHsDriver() *hsDriver {
	resources := &Resources{
		Signer: core.GenerateKey(nil),
	}
	state := newState(resources)
	return &hsDriver{
		resources: resources,
		config:    DefaultConfig,
		state:     state,
	}
}

func TestHsDriver_TestMajorityCount(t *testing.T) {
	hsd := setupTestHsDriver()
	validators := []string{
		core.GenerateKey(nil).PublicKey().String(),
		core.GenerateKey(nil).PublicKey().String(),
		core.GenerateKey(nil).PublicKey().String(),
		core.GenerateKey(nil).PublicKey().String(),
	}
	hsd.resources.VldStore = core.NewValidatorStore(validators, validators)

	res := hsd.MajorityValidatorCount()

	assert := assert.New(t)
	assert.Equal(hsd.resources.VldStore.MajorityValidatorCount(), res)
}

func TestHsDriver_CreateLeaf(t *testing.T) {
	hsd := setupTestHsDriver()
	parent := newHsBlock(core.NewBlock().Sign(hsd.resources.Signer), hsd.state)
	hsd.state.setBlock(parent.(*hsBlock).block)
	qc := newHsQC(core.NewQuorumCert(), hsd.state)
	height := uint64(5)

	storage := new(MockStorage)
	storage.On("GetBlockHeight").Return(2) // driver should get bexec height from storage
	storage.On("GetMerkleRoot").Return([]byte("merkle-root"))
	hsd.resources.Storage = storage

	signer := core.GenerateKey(nil)
	hsd.resources.VldStore = core.NewValidatorStore([]string{signer.PublicKey().String()}, []string{signer.PublicKey().String()})

	tx1, tx2 := []byte("tx1"), []byte("tx2")
	txsInQ := [][]byte{tx1, tx2}
	if ExecuteTxFlag {
		storage.On("HasTx", tx1).Return(false)
		storage.On("HasTx", tx2).Return(false)
	}

	batch := core.NewBatch().Header().SetTransactions(txsInQ).Sign(signer)
	if VoteBatchFlag {
		hsd.leaderState = newLeaderState().setBatchWaitTime(3 * time.Second).setBatchSignLimit(1).setBlockBatchLimit(1)
		batchVote := core.NewBatchVote().Build([]*core.BatchHeader{batch}, signer)
		hsd.leaderState.addBatchVote(batchVote)
	} else {
		hsd.voterState = newVoterState().setVoteBatchLimit(1)
		hsd.voterState.addBatch(batch, 0, len(batch.Transactions()))
	}

	txpool := new(MockTxPool)
	if ExecuteTxFlag {
		txpool.On("SyncTxs", batch.Proposer(), batch.Transactions()).Return(nil)
	}
	hsd.resources.TxPool = txpool

	leaf := hsd.CreateLeaf(parent, qc, height)

	storage.AssertExpectations(t)

	assert := assert.New(t)
	assert.NotNil(leaf)
	assert.True(parent.Equal(leaf.Parent()), "should link to parent")
	assert.Equal(qc, leaf.Justify(), "should add qc")
	assert.Equal(height, leaf.Height())

	blk := leaf.(*hsBlock).block
	assert.Equal(txsInQ, blk.Transactions())
	assert.EqualValues(2, blk.ExecHeight())
	assert.Equal([]byte("merkle-root"), blk.MerkleRoot())
	assert.NotEmpty(blk.Timestamp(), "should add timestamp")

	assert.NotNil(hsd.state.getBlock(blk.Hash()), "should store leaf block in state")
}

func TestHsDriver_VoteBlock(t *testing.T) {
	hsd := setupTestHsDriver()
	hsd.checkTxDelay = time.Millisecond
	hsd.config.TxWaitTime = 20 * time.Millisecond

	proposer := core.GenerateKey(nil)
	blk := core.NewBlock().Sign(proposer)

	validators := []string{blk.Proposer().String()}
	hsd.resources.VldStore = core.NewValidatorStore(validators, validators)

	txPool := new(MockTxPool)
	txPool.On("GetStatus").Return(txpool.Status{}) // no txs in the pool
	if !PreserveTxFlag {
		txPool.On("SetTxsPending", blk.Transactions())
	}
	hsd.resources.TxPool = txPool

	// should sign block and send vote
	msgSvc := new(MockMsgService)
	msgSvc.On("SendVote", proposer.PublicKey(), blk.Vote(hsd.resources.Signer)).Return(nil)
	hsd.resources.MsgSvc = msgSvc

	start := time.Now()
	hsd.VoteBlock(newHsBlock(blk, hsd.state))
	elapsed := time.Since(start)

	txPool.AssertExpectations(t)
	msgSvc.AssertExpectations(t)

	assert := assert.New(t)
	assert.GreaterOrEqual(elapsed, hsd.config.TxWaitTime, "should delay if no txs in the pool")

	txPool = new(MockTxPool)
	txPool.On("GetStatus").Return(txpool.Status{Total: 1}) // one txs in the pool
	if !PreserveTxFlag {
		txPool.On("SetTxsPending", blk.Transactions())
	}
	hsd.resources.TxPool = txPool

	start = time.Now()
	hsd.VoteBlock(newHsBlock(blk, hsd.state))
	elapsed = time.Since(start)

	txPool.AssertExpectations(t)
	msgSvc.AssertExpectations(t)

	assert.Less(elapsed, hsd.config.TxWaitTime, "should not delay if txs in the pool")
}

func TestHsDriver_Commit(t *testing.T) {
	hsd := setupTestHsDriver()
	parent := core.NewBlock().SetHeight(10).Sign(hsd.resources.Signer)
	batch := core.NewBatch().Header().SetTransactions([][]byte{[]byte("txfromfolk")}).Sign(hsd.resources.Signer)
	bfolk := core.NewBlock().SetBatchHeaders([]*core.BatchHeader{batch}, true).SetHeight(10).Sign(hsd.resources.Signer)

	tx := core.NewTransaction().Sign(hsd.resources.Signer)
	batch2 := core.NewBatch().SetTransactions([]*core.Transaction{tx}).Header().Sign(hsd.resources.Signer)
	bexec := core.NewBlock().SetBatchHeaders([]*core.BatchHeader{batch2}, true).SetParentHash(parent.Hash()).SetHeight(11).Sign(hsd.resources.Signer)
	hsd.state.setBlock(parent)
	hsd.state.setCommittedBlock(parent)
	hsd.state.setBlock(bfolk)
	hsd.state.setBlock(bexec)

	txs := []*core.Transaction{tx}
	txPool := new(MockTxPool)
	if ExecuteTxFlag {
		txPool.On("GetTxsToExecute", bexec.Transactions()).Return(txs, nil)
	}
	if !PreserveTxFlag {
		txPool.On("RemoveTxs", bexec.Transactions()).Once() // should remove txs from pool after commit
	}
	txPool.On("PutTxsToQueue", bfolk.Transactions()).Once() // should put txs of folked block back to queue from pending
	hsd.resources.TxPool = txPool

	bcm := core.NewBlockCommit().SetHash(bexec.Hash())
	txcs := []*core.TxCommit{core.NewTxCommit().SetHash(tx.Hash())}
	cdata := &storage.CommitData{
		Block:        bexec,
		Transactions: txs,
		BlockCommit:  bcm,
		TxCommits:    txcs,
	}

	execution := new(MockExecution)
	if ExecuteTxFlag {
		execution.On("Execute", bexec, txs).Return(bcm, txcs)
	} else {
		cdata.Transactions = nil
		execution.On("MockExecute", bexec).Return(bcm, txcs)
	}
	hsd.resources.Execution = execution

	storage := new(MockStorage)
	storage.On("Commit", cdata).Return(nil)
	hsd.resources.Storage = storage

	hsd.Commit(newHsBlock(bexec, hsd.state))

	txPool.AssertExpectations(t)
	execution.AssertExpectations(t)
	storage.AssertExpectations(t)

	assert := assert.New(t)
	assert.NotNil(hsd.state.getBlockFromState(bexec.Hash()),
		"should not delete bexec from state")
	assert.Nil(hsd.state.getBlockFromState(bfolk.Hash()),
		"should delete folked block from state")
}

func TestHsDriver_CreateQC(t *testing.T) {
	hsd := setupTestHsDriver()
	blk := core.NewBlock().Sign(hsd.resources.Signer)
	hsd.state.setBlock(blk)
	votes := []hotstuff.Vote{
		newHsVote(blk.ProposerVote(), hsd.state),
		newHsVote(blk.Vote(core.GenerateKey(nil)), hsd.state),
	}
	qc := hsd.CreateQC(votes)

	assert := assert.New(t)
	assert.Equal(blk, qc.Block().(*hsBlock).block, "should get qc reference block")
}

func TestHsDriver_BroadcastProposal(t *testing.T) {
	hsd := setupTestHsDriver()
	blk := core.NewBlock().Sign(hsd.resources.Signer)
	hsd.state.setBlock(blk)

	msgSvc := new(MockMsgService)
	msgSvc.On("BroadcastProposal", blk).Return(nil)
	hsd.resources.MsgSvc = msgSvc

	hsd.BroadcastProposal(newHsBlock(blk, hsd.state))

	msgSvc.AssertExpectations(t)
}
