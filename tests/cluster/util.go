// Copyright (C) 2021 Aung Maw
// Licensed under the GNU General Public License v3.0

package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/multiformats/go-multiaddr"
	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/node"
)

func ReadRemoteHosts(hostsPath string, nodeCount int) ([]string, error) {
	raw, err := os.ReadFile(hostsPath)
	if err != nil {
		return nil, err
	}
	hosts := strings.Split(string(raw), "\n")
	if len(hosts) < nodeCount {
		return nil, fmt.Errorf("not enough hosts, %d | %d",
			len(hosts), nodeCount)
	}
	return hosts[:nodeCount], nil
}

func WriteNodeKey(datadir string, key *core.PrivateKey) error {
	f, err := os.Create(path.Join(datadir, node.NodekeyFile))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(key.Bytes())
	return err
}

func WriteGenesisFile(datadir string, genesis *node.Genesis) error {
	f, err := os.Create(path.Join(datadir, node.GenesisFile))
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	return e.Encode(genesis)
}

func WritePeersFile(datadir string, peers []node.Peer) error {
	f, err := os.Create(path.Join(datadir, node.PeersFile))
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	return e.Encode(peers)
}

func MakeRandomKeys(count int) []*core.PrivateKey {
	keys := make([]*core.PrivateKey, count)
	for i := 0; i < count; i++ {
		keys[i] = core.GenerateKey(nil)
	}
	return keys
}

func MakePeers(keys []*core.PrivateKey, addrs []multiaddr.Multiaddr) []node.Peer {
	vlds := make([]node.Peer, len(addrs))
	// create validator infos (pubkey + addr)
	for i, addr := range addrs {
		vlds[i] = node.Peer{
			PubKey: keys[i].PublicKey().Bytes(),
			Addr:   addr.String(),
		}
	}
	return vlds
}

func SetupTemplateDir(dir string, keys []*core.PrivateKey, vlds []node.Peer) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}
	genesis := &node.Genesis{
		Workers: make([]string, 0, 0),
		Voters:  make([]string, 0, 0),
		Weights: make([]int, 0, 0),
	}
	for i, v := range keys {
		if i%3 == 0 {
			genesis.Workers = append(genesis.Workers, v.PublicKey().String())
			genesis.Weights = append(genesis.Weights, 1)
			genesis.Voters = append(genesis.Voters, v.PublicKey().String())
		}
		if i%3 == 1 {
			genesis.Workers = append(genesis.Workers, v.PublicKey().String())
			genesis.Weights = append(genesis.Weights, 1)
		}
		if i%3 == 2 {
			genesis.Voters = append(genesis.Voters, v.PublicKey().String())
		}
	}
	for i, key := range keys {
		dir := path.Join(dir, strconv.Itoa(i))
		os.Mkdir(dir, 0755)
		if err := WriteNodeKey(dir, key); err != nil {
			return err
		}
		if err := WriteGenesisFile(dir, genesis); err != nil {
			return err
		}
		if err := WritePeersFile(dir, vlds); err != nil {
			return err
		}
	}
	return nil
}

func RunCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	fmt.Printf(" $ %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func AddJuriaFlags(cmd *exec.Cmd, config *node.Config) {
	cmd.Args = append(cmd.Args, "-d", config.Datadir)
	cmd.Args = append(cmd.Args, "-p", strconv.Itoa(config.Port))
	cmd.Args = append(cmd.Args, "-P", strconv.Itoa(config.APIPort))
	cmd.Args = append(cmd.Args, "--debug", strconv.FormatBool(config.Debug))

	cmd.Args = append(cmd.Args, "--storage-merkleBranchFactor",
		strconv.Itoa(int(config.StorageConfig.MerkleBranchFactor)))

	cmd.Args = append(cmd.Args, "--execution-txExecTimeout",
		config.ExecutionConfig.TxExecTimeout.String(),
	)
	cmd.Args = append(cmd.Args, "--execution-concurrentLimit",
		strconv.Itoa(config.ExecutionConfig.ConcurrentLimit))

	cmd.Args = append(cmd.Args, "--chainid",
		strconv.Itoa(int(config.ConsensusConfig.ChainID)))

	cmd.Args = append(cmd.Args, "--consensus-blockTxLimit",
		strconv.Itoa(config.ConsensusConfig.BatchTxLimit))

	cmd.Args = append(cmd.Args, "--consensus-txWaitTime",
		config.ConsensusConfig.TxWaitTime.String())

	cmd.Args = append(cmd.Args, "--consensus-proposeTimeout",
		config.ConsensusConfig.ProposeTimeout.String())

	cmd.Args = append(cmd.Args, "--consensus-blockDelay",
		config.ConsensusConfig.BlockDelay.String())

	cmd.Args = append(cmd.Args, "--consensus-viewWidth",
		config.ConsensusConfig.ViewWidth.String())

	cmd.Args = append(cmd.Args, "--consensus-leaderTimeout",
		config.ConsensusConfig.LeaderTimeout.String())
}