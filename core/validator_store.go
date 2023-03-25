// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package core

import (
	"encoding/base64"
	"math"

	"github.com/wooyang2018/ppov-blockchain/logger"
)

// ValidatorStore godoc
type ValidatorStore interface {
	VoterCount() int                      //返回投票节点数量
	MajorityVoterCount() int              //返回大多数投票节点数量
	WorkerCount() int                     //返回记账节点数量
	EnoughWorkerCount() int               //返回足够记账节点数量
	ValidatorCount() int                  //返回验证节点数量
	MajorityValidatorCount() int          //返回大多数验证节点数量
	IsVoter(pubKey *PublicKey) bool       //返回指定公钥是否投票节点
	IsWorker(pubKey *PublicKey) bool      //返回指定公钥是否记账节点
	GetVoter(idx int) *PublicKey          //获取指定索引的投票节点
	GetWorker(idx int) *PublicKey         //获取指定索引的记账节点
	GetVoterIndex(pubKey *PublicKey) int  //获取指定公钥的投票节点的索引
	GetWorkerIndex(pubKey *PublicKey) int //获取指定公钥的记账节点的索引
	GetWorkerWeight(idx int) int          //获取指定索引的记账节点权重
}

type ppovValidatorStore struct {
	voters  []*PublicKey //投票节点列表
	workers []*PublicKey //记账节点列表
	weights []int        //记账节点权重列表

	vcount    int //投票节点收集区块的阈值
	voterMap  map[string]int
	workerMap map[string]int

	validators []*PublicKey //投票节点和记账节点的集合
}

var _ ValidatorStore = (*ppovValidatorStore)(nil)

func StringToPubKey(v string) *PublicKey {
	key, err := base64.StdEncoding.DecodeString(v)
	pubKey, err := NewPublicKey(key)
	if err != nil {
		logger.I().Fatalw("parse voter failed", "error", err)
	}
	return pubKey
}

func NewValidatorStore(workers []string, weights []int, voters []string) ValidatorStore {
	store := &ppovValidatorStore{
		weights: weights,
	}

	set := make(map[string]*PublicKey)
	for _, v := range workers {
		set[v] = StringToPubKey(v)
	}
	for _, v := range voters {
		if _, ok := set[v]; !ok {
			set[v] = StringToPubKey(v)
		}
	}

	store.validators = make([]*PublicKey, 0, len(set))
	for _, v := range set {
		store.validators = append(store.validators, v)
	}

	store.voters = make([]*PublicKey, len(voters))
	for i, v := range voters {
		store.voters[i] = set[v]
	}
	store.workers = make([]*PublicKey, len(workers))
	for i, v := range workers {
		store.workers[i] = set[v]
	}

	store.vcount = len(workers) //TODO 投票节点收集区块的阈值

	store.voterMap = make(map[string]int, len(store.voters))
	for i, v := range store.voters {
		store.voterMap[v.String()] = i
	}
	store.workerMap = make(map[string]int, len(store.workers))
	for i, v := range store.workers {
		store.workerMap[v.String()] = i
	}

	return store
}

func (store *ppovValidatorStore) VoterCount() int {
	return len(store.voters)
}

func (store *ppovValidatorStore) MajorityVoterCount() int {
	return MajorityCount(store.VoterCount())
}

func (store *ppovValidatorStore) WorkerCount() int {
	return len(store.workers)
}

func (store *ppovValidatorStore) EnoughWorkerCount() int {
	return len(store.workers)
}

func (store *ppovValidatorStore) ValidatorCount() int {
	return len(store.validators)
}

func (store *ppovValidatorStore) MajorityValidatorCount() int {
	return MajorityCount(len(store.validators))
}

func (store *ppovValidatorStore) IsVoter(pubKey *PublicKey) bool {
	if pubKey == nil {
		return false
	}
	_, ok := store.voterMap[pubKey.String()]
	return ok
}

func (store *ppovValidatorStore) IsWorker(pubKey *PublicKey) bool {
	if pubKey == nil {
		return false
	}
	_, ok := store.workerMap[pubKey.String()]
	return ok
}

func (store *ppovValidatorStore) GetVoter(idx int) *PublicKey {
	if idx >= len(store.voters) || idx < 0 {
		return nil
	}
	return store.voters[idx]
}

func (store *ppovValidatorStore) GetWorker(idx int) *PublicKey {
	if idx >= len(store.workers) || idx < 0 {
		return nil
	}
	return store.workers[idx]
}

func (store *ppovValidatorStore) GetVoterIndex(pubKey *PublicKey) int {
	if pubKey == nil {
		return 0
	}
	return store.voterMap[pubKey.String()]
}

func (store *ppovValidatorStore) GetWorkerIndex(pubKey *PublicKey) int {
	if pubKey == nil {
		return 0
	}
	return store.workerMap[pubKey.String()]
}

func (store *ppovValidatorStore) GetWorkerWeight(idx int) int {
	if idx >= len(store.weights) || idx < 0 {
		return -1
	}
	return store.weights[idx]
}

// MajorityCount returns 2f + 1 members
func MajorityCount(validatorCount int) int {
	// n=3f+1 -> f=floor((n-1)3) -> m=n-f -> m=ceil((2n+1)/3)
	return int(math.Ceil(float64(2*validatorCount+1) / 3))
}
