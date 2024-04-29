// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package p2p

import (
	"strconv"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/wooyang2018/ppov-blockchain/core"
)

type PeerStore struct {
	pub2peer map[string]*Peer
	id2name  map[peer.ID]string
	mtx      sync.RWMutex
}

func NewPeerStore() *PeerStore {
	return &PeerStore{
		pub2peer: make(map[string]*Peer),
		id2name:  make(map[peer.ID]string),
	}
}

func (s *PeerStore) Load(pubKey *core.PublicKey) *Peer {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.pub2peer[pubKey.String()]
}

func (s *PeerStore) Store(p *Peer) *Peer {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.pub2peer[p.PublicKey().String()] = p
	id, err := getIDFromPublicKey(p.PublicKey())
	if err != nil {
		panic(nil)
	}
	s.id2name[id] = p.name
	return p
}

func (s *PeerStore) Delete(pubKey *core.PublicKey) *Peer {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	p := s.pub2peer[pubKey.String()]
	delete(s.pub2peer, pubKey.String())
	id, err := getIDFromPublicKey(p.PublicKey())
	if err != nil {
		panic(nil)
	}
	delete(s.id2name, id)
	return p
}

func (s *PeerStore) List() []*Peer {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	peers := make([]*Peer, 0, len(s.pub2peer))
	for _, p := range s.pub2peer {
		peers = append(peers, p)
	}
	return peers
}

func (s *PeerStore) LoadOrStore(p *Peer) (actual *Peer, loaded bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if actual, loaded = s.pub2peer[p.PublicKey().String()]; loaded {
		return actual, loaded
	}
	s.pub2peer[p.PublicKey().String()] = p
	id, err := getIDFromPublicKey(p.PublicKey())
	if err != nil {
		panic(nil)
	}
	s.id2name[id] = p.name
	return p, false
}

func (s *PeerStore) IsValidID(id peer.ID) bool {
	_, ok := s.id2name[id]
	return ok
}

func (s *PeerStore) GetPeerDist(src string) map[peer.ID]int {
	res := make(map[peer.ID]int)
	for id, name := range s.id2name {
		res[id] = hammingWeight(src, name)
	}
	return res
}

func hammingWeight(name1 string, name2 string) (cnt int) {
	n1, _ := strconv.Atoi(name1)
	n2, _ := strconv.Atoi(name2)
	num := n1 ^ n2
	for num > 0 {
		num = num & (num - 1) // 消除最后一位1
		cnt++
	}
	return
}
