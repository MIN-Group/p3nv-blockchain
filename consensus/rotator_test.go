// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/hotstuff"
)

func setupRotator() (*rotator, *core.Block) {
	key1 := core.GenerateKey(nil)
	key2 := core.GenerateKey(nil)
	keys := []string{
		key1.PublicKey().String(),
		key2.PublicKey().String(),
	}

	resources := &Resources{VldStore: core.NewValidatorStore(keys, keys)}

	b0 := core.NewBlock().Sign(key1)
	q0 := core.NewQuorumCert().Build([]*core.Vote{b0.ProposerVote()})
	b0.SetQuorumCert(q0)

	state := newState(resources)
	state.setBlock(b0)
	hsDriver := &hsDriver{
		resources: resources,
		state:     state,
	}
	hotstuff := hotstuff.New(hsDriver, nil, newHsBlock(b0, state), newHsQC(q0, state))
	return &rotator{
		resources: resources,
		config:    DefaultConfig,
		state:     state,
		hotstuff:  hotstuff,
	}, b0
}

func Test_rotator_isNewViewApproval(t *testing.T) {
	assert := assert.New(t)

	rot1, _ := setupRotator()
	rot2, _ := setupRotator()

	rot1.setPendingViewChange(true)
	rot2.setPendingViewChange(false)

	tests := []struct {
		name        string
		rot         *rotator
		proposerIdx int
		want        bool
	}{
		{"pending and same leader", rot1, 0, true},
		{"not pending and different leader", rot2, 1, true},
		{"pending and different leader", rot1, 1, false},
		{"not pending and same leader", rot2, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(tt.want, tt.rot.isNewViewApproval(tt.proposerIdx))
		})
	}
}

func TestRotator_resetViewTimer(t *testing.T) {
	assert := assert.New(t)

	rot, _ := setupRotator()
	rot.setPendingViewChange(true)

	rot.approveViewLeader(1)

	assert.False(rot.getPendingViewChange())
	assert.EqualValues(rot.state.getLeaderIndex(), 1)
}
