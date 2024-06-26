// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/emitter"
	"github.com/wooyang2018/ppov-blockchain/p2p"
	"github.com/wooyang2018/ppov-blockchain/storage"
	"github.com/wooyang2018/ppov-blockchain/txpool"
)

type TxPool interface {
	SubmitTx(tx *core.Transaction) error
	StoreTxs(txs *core.TxList) error
	StorePendingTxs(txs *core.TxList) error
	PopTxsFromQueue(max int) []*core.Transaction
	GetTxsFromQueue(max int) []*core.Transaction
	SetTxsPending(hashes [][]byte)
	GetTxsToExecute(hashes [][]byte) ([]*core.Transaction, [][]byte)
	RemoveTxs(hashes [][]byte)
	PutTxsToQueue(hashes [][]byte)
	SyncTxs(peer *core.PublicKey, hashes [][]byte) error
	GetTx(hash []byte) *core.Transaction
	GetStatus() txpool.Status
	GetTxStatus(hash []byte) txpool.TxStatus
}

type Storage interface {
	GetMerkleRoot() []byte
	Commit(data *storage.CommitData) error
	GetBlock(hash []byte) (*core.Block, error)
	GetLastBlock() (*core.Block, error)
	GetLastQC() (*core.QuorumCert, error)
	GetBlockHeight() uint64
	HasTx(hash []byte) bool
}

type MsgService interface {
	BroadcastProposal(blk *core.Block) error
	BroadcastBatch(batch *core.Batch) error
	BroadcastNewView(qc *core.QuorumCert) error
	SendBatch(pubKey *core.PublicKey, batch *core.Batch) error
	SendVote(pubKey *core.PublicKey, vote *core.Vote) error
	SendBatchVote(pubKey *core.PublicKey, vote *core.BatchVote) error
	RequestBlock(pubKey *core.PublicKey, hash []byte) (*core.Block, error)
	RequestBlockByHeight(pubKey *core.PublicKey, height uint64) (*core.Block, error)
	SendNewView(pubKey *core.PublicKey, qc *core.QuorumCert) error
	SubscribeBatch(buffer int) *emitter.Subscription
	SubscribeProposal(buffer int) *emitter.Subscription
	SubscribeVote(buffer int) *emitter.Subscription
	SubscribeBatchVote(buffer int) *emitter.Subscription
	SubscribeNewView(buffer int) *emitter.Subscription
}

type Execution interface {
	Execute(blk *core.Block, txs []*core.Transaction) (*core.BlockCommit, []*core.TxCommit)
	MockExecute(blk *core.Block) (*core.BlockCommit, []*core.TxCommit)
}

type Resources struct {
	Signer    core.Signer
	VldStore  core.ValidatorStore
	Storage   Storage
	MsgSvc    MsgService
	Host      *p2p.Host
	TxPool    TxPool
	Execution Execution
}
