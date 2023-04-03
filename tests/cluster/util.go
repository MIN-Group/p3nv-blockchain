// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package cluster

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/multiformats/go-multiaddr"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/node"
)

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

func SetupTemplateDir(dir string, keys []*core.PrivateKey, vlds []node.Peer, WorkerProportion, VoterProportion float32) error {
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

	workers := PickUniqueRandoms(len(keys), int(float32(len(keys))*WorkerProportion))
	fmt.Printf("Setup workers: %v\n", workers)
	for _, v := range workers {
		genesis.Workers = append(genesis.Workers, keys[v].PublicKey().String())
		genesis.Weights = append(genesis.Weights, 1)
	}

	// Ensure that the node is either a Worker or a Voter
	var voters []int
	unselectedIndexes := GetUnselectedIndexes(len(keys), workers)
	if len(unselectedIndexes) <= int(float32(len(keys))*VoterProportion) {
		voters = append(voters, unselectedIndexes...)
		indexes := PickUniqueRandoms(len(workers), int(float32(len(keys))*VoterProportion)-len(unselectedIndexes))
		for _, v := range indexes {
			voters = append(voters, workers[v])
		}
	} else {
		indexes := PickUniqueRandoms(len(unselectedIndexes), int(float32(len(keys))*VoterProportion))
		for _, v := range indexes {
			voters = append(voters, unselectedIndexes[v])
		}
	}
	fmt.Printf("Setup voters: %v\n", voters)
	for _, v := range voters {
		genesis.Voters = append(genesis.Voters, keys[v].PublicKey().String())
	}

	for i, key := range keys {
		d := path.Join(dir, strconv.Itoa(i))
		os.Mkdir(d, 0755)
		if err := WriteNodeKey(d, key); err != nil {
			return err
		}
		if err := WriteGenesisFile(d, genesis); err != nil {
			return err
		}
		if err := WritePeersFile(d, vlds); err != nil {
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

func AddPPoVFlags(cmd *exec.Cmd, config *node.Config) {
	cmd.Args = append(cmd.Args, "-d", config.Datadir)
	cmd.Args = append(cmd.Args, "-p", strconv.Itoa(config.Port))
	cmd.Args = append(cmd.Args, "-P", strconv.Itoa(config.APIPort))
	if config.Debug {
		cmd.Args = append(cmd.Args, "--debug")
	}
	if config.BroadcastTx {
		cmd.Args = append(cmd.Args, "--broadcast-tx")
	}

	cmd.Args = append(cmd.Args, "--storage-merkleBranchFactor",
		strconv.Itoa(int(config.StorageConfig.MerkleBranchFactor)))

	cmd.Args = append(cmd.Args, "--execution-txExecTimeout",
		config.ExecutionConfig.TxExecTimeout.String(),
	)

	cmd.Args = append(cmd.Args, "--execution-concurrentLimit",
		strconv.Itoa(config.ExecutionConfig.ConcurrentLimit))

	cmd.Args = append(cmd.Args, "--chainID",
		strconv.Itoa(int(config.ConsensusConfig.ChainID)))

	cmd.Args = append(cmd.Args, "--consensus-batchTxLimit",
		strconv.Itoa(config.ConsensusConfig.BatchTxLimit))

	cmd.Args = append(cmd.Args, "--consensus-blockBatchLimit",
		strconv.Itoa(config.ConsensusConfig.BlockBatchLimit))

	cmd.Args = append(cmd.Args, "--consensus-voteBatchLimit",
		strconv.Itoa(config.ConsensusConfig.VoteBatchLimit))

	cmd.Args = append(cmd.Args, "--consensus-txWaitTime",
		config.ConsensusConfig.TxWaitTime.String())

	cmd.Args = append(cmd.Args, "--consensus-batchWaitTime",
		config.ConsensusConfig.BatchWaitTime.String())

	cmd.Args = append(cmd.Args, "--consensus-proposeTimeout",
		config.ConsensusConfig.ProposeTimeout.String())

	cmd.Args = append(cmd.Args, "--consensus-batchTimeout",
		config.ConsensusConfig.BatchTimeout.String())

	cmd.Args = append(cmd.Args, "--consensus-blockDelay",
		config.ConsensusConfig.BlockDelay.String())

	cmd.Args = append(cmd.Args, "--consensus-viewWidth",
		config.ConsensusConfig.ViewWidth.String())

	cmd.Args = append(cmd.Args, "--consensus-leaderTimeout",
		config.ConsensusConfig.LeaderTimeout.String())

	cmd.Args = append(cmd.Args, "--consensus-benchmarkPath",
		config.ConsensusConfig.BenchmarkPath)
}

func PickUniqueRandoms(total, count int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	unique := make(map[int]struct{}, count)
	for len(unique) < count {
		unique[r.Intn(total)] = struct{}{}
	}
	ret := make([]int, 0, count)
	for v := range unique {
		ret = append(ret, v)
	}
	return ret
}

func GetUnselectedIndexes(total int, selected []int) []int {
	smap := make(map[int]struct{}, len(selected))
	for _, idx := range selected {
		smap[idx] = struct{}{}
	}
	ret := make([]int, 0, total-len(selected))
	for i := 0; i < total; i++ {
		if _, found := smap[i]; !found {
			ret = append(ret, i)
		}
	}
	return ret
}

// Sleep print duration and call time.Sleep
func Sleep(d time.Duration) {
	fmt.Printf("Wait for %s\n", d)
	time.Sleep(d)
}
