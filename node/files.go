// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/multiformats/go-multiaddr"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/p2p"
)

type Peer struct {
	PubKey    []byte
	PointAddr string
	TopicAddr string
	Name      string
}

type Genesis struct {
	Workers []string // 记账节点列表
	Voters  []string // 投票节点列表
}

const (
	NodekeyFile = "nodekey"
	GenesisFile = "genesis.json"
	PeersFile   = "peers.json"
)

func readNodeKey(datadir string) (*core.PrivateKey, error) {
	b, err := os.ReadFile(path.Join(datadir, NodekeyFile))
	if err != nil {
		return nil, fmt.Errorf("cannot read %s, %w", NodekeyFile, err)
	}
	return core.NewPrivateKey(b)
}

func readGenesis(datadir string) (*Genesis, error) {
	f, err := os.Open(path.Join(datadir, GenesisFile))
	if err != nil {
		return nil, fmt.Errorf("cannot read %s, %w", GenesisFile, err)
	}
	defer f.Close()
	genesis := new(Genesis)
	if err := json.NewDecoder(f).Decode(&genesis); err != nil {
		return nil, fmt.Errorf("cannot parse %s, %w", GenesisFile, err)
	}
	return genesis, nil
}

func readPeers(datadir string) ([]*p2p.Peer, error) {
	f, err := os.Open(path.Join(datadir, PeersFile))
	if err != nil {
		return nil, fmt.Errorf("cannot read %s, %w", PeersFile, err)
	}
	defer f.Close()

	var raws []Peer
	if err := json.NewDecoder(f).Decode(&raws); err != nil {
		return nil, fmt.Errorf("cannot parse %s, %w", PeersFile, err)
	}

	peers := make([]*p2p.Peer, len(raws))

	for i, r := range raws {
		pubKey, err := core.NewPublicKey(r.PubKey)
		if err != nil {
			return nil, fmt.Errorf("invalid public key %w", err)
		}
		pointAddr, err := multiaddr.NewMultiaddr(r.PointAddr)
		topicAddr, err := multiaddr.NewMultiaddr(r.TopicAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid multiaddr %w", err)
		}
		if r.Name == "" {
			return nil, errors.New("name can not be empty")
		}
		peers[i] = p2p.NewPeer(pubKey, pointAddr, topicAddr)
		peers[i].SetName(r.Name)
	}
	return peers, nil
}
