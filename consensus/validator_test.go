// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wooyang2018/ppov-blockchain/core"
)

func TestValidator_verifyProposalToVote(t *testing.T) {
	priv0 := core.GenerateKey(nil)
	priv1 := core.GenerateKey(nil)
	keys := []string{
		priv0.PublicKey().String(),
		priv1.PublicKey().String(),
	}
	resources := &Resources{
		VldStore: core.NewValidatorStore(keys, keys),
	}
	mStrg := new(MockStorage)
	mTxPool := new(MockTxPool)

	resources.Storage = mStrg
	resources.TxPool = mTxPool

	mRoot := []byte("merkle-root")
	mStrg.On("GetBlockHeight").Return(10)
	mStrg.On("GetMerkleRoot").Return(mRoot)

	// valid tx
	tx1 := core.NewTransaction().SetExpiry(15).Sign(core.GenerateKey(nil))
	// committed tx
	tx2 := core.NewTransaction().SetExpiry(9).Sign(core.GenerateKey(nil))
	// expired tx
	tx3 := core.NewTransaction().SetExpiry(13).Sign(core.GenerateKey(nil))
	// no expiry tx (should only used for test)
	tx4 := core.NewTransaction().Sign(core.GenerateKey(nil))
	// not found tx
	// This should not happen at run time.
	// Not found tx means sync txs failed. If sync failed, cannot vote already
	tx5 := core.NewTransaction().SetExpiry(15).Sign(core.GenerateKey(nil))

	header1 := core.NewBatch().Header().SetTransactions([][]byte{tx1.Hash(), tx4.Hash()}).Sign(priv0)
	header2 := core.NewBatch().Header().SetTransactions([][]byte{tx1.Hash(), tx2.Hash(), tx4.Hash()}).Sign(priv0)
	header3 := core.NewBatch().Header().SetTransactions([][]byte{tx1.Hash(), tx3.Hash(), tx4.Hash()}).Sign(priv0)
	header4 := core.NewBatch().Header().SetTransactions([][]byte{tx1.Hash(), tx5.Hash(), tx4.Hash()}).Sign(priv0)

	mStrg.On("HasTx", tx1.Hash()).Return(false)
	mStrg.On("HasTx", tx2.Hash()).Return(true)
	mStrg.On("HasTx", tx3.Hash()).Return(false)
	mStrg.On("HasTx", tx4.Hash()).Return(false)
	mStrg.On("HasTx", tx5.Hash()).Return(false)

	mTxPool.On("GetTx", tx1.Hash()).Return(tx1)
	mTxPool.On("GetTx", tx3.Hash()).Return(tx3)
	mTxPool.On("GetTx", tx4.Hash()).Return(tx4)
	mTxPool.On("GetTx", tx5.Hash()).Return(nil)

	vld := &validator{
		resources: resources,
		state:     newState(resources),
	}
	vld.state.committedHeight = mStrg.GetBlockHeight()
	vld.state.setLeaderIndex(1)
	type testCase struct {
		name     string
		valid    bool
		proposal *core.Block
	}
	tests := []testCase{
		{"valid", true, core.NewBlock().
			SetHeight(14).SetExecHeight(10).SetMerkleRoot(mRoot).
			SetBatchHeaders([]*core.BatchHeader{header1}, true).
			Sign(priv1),
		},
		{"proposer is not leader", false, core.NewBlock().
			SetHeight(14).SetExecHeight(10).SetMerkleRoot(mRoot).
			SetBatchHeaders([]*core.BatchHeader{header1}, true).
			Sign(priv0),
		},
		{"different exec height", false, core.NewBlock().
			SetHeight(14).SetExecHeight(9).SetMerkleRoot(mRoot).
			SetBatchHeaders([]*core.BatchHeader{header1}, true).
			Sign(priv1),
		},
	}
	if ExecuteTxFlag {
		tests = append(tests, []testCase{
			{"different merkle root", false, core.NewBlock().
				SetHeight(14).SetExecHeight(10).SetMerkleRoot([]byte("different")).
				SetBatchHeaders([]*core.BatchHeader{header1}, true).
				Sign(priv1),
			},
			{"committed tx", false, core.NewBlock().
				SetHeight(14).SetExecHeight(10).SetMerkleRoot(mRoot).
				SetBatchHeaders([]*core.BatchHeader{header2}, true).
				Sign(priv1),
			},
			{"expired tx", false, core.NewBlock().
				SetHeight(14).SetExecHeight(10).SetMerkleRoot(mRoot).
				SetBatchHeaders([]*core.BatchHeader{header3}, true).
				Sign(priv1),
			},
			{"not found tx", false, core.NewBlock().
				SetHeight(14).SetExecHeight(10).SetMerkleRoot(mRoot).
				SetBatchHeaders([]*core.BatchHeader{header4}, true).
				Sign(priv1),
			},
		}...)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if tt.valid {
				assert.NoError(vld.verifyProposalToVote(tt.proposal))
			} else {
				assert.Error(vld.verifyProposalToVote(tt.proposal))
			}
		})
	}
}
