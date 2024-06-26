// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/wooyang2018/ppov-blockchain/consensus"
	"github.com/wooyang2018/ppov-blockchain/node"
	"github.com/wooyang2018/ppov-blockchain/tests/cluster"
	"github.com/wooyang2018/ppov-blockchain/tests/experiments"
	"github.com/wooyang2018/ppov-blockchain/tests/testutil"
)

var (
	WorkDir     = "./workdir"
	NodeCount   = 64
	WorkerCount = NodeCount / 2
	VoterCount  = NodeCount

	LoadTxPerSec    = 10   // tps for client to submit tx during functional testing
	LoadJobPerTick  = 1000 // num of tasks to be completed per tick
	LoadSubmitNodes = []int{}
	LoadBatchSubmit = true // whether to enable batch transaction submission

	// chaincode priority: empty > ppovcoin bincc > ppovcoin
	EmptyChainCode = true  // deploy empty chaincode instead of ppovcoin
	PPoVCoinBinCC  = false // deploy ppovcoin chaincode as bincc type (not embeded in ppov node)
	CheckRotation  = false
	BroadcastTx    = false

	// run tests in remote linux cluster
	RemoteLinuxCluster    = false // if false it'll use local cluster (running multiple nodes on single local machine)
	RemoteSetupRequired   = true
	RemoteInstallRequired = false // if false it will not try to install dstat on remote machine
	RemoteRunRequired     = false // if false it will not run dstat on remote machine
	RemoteKeySSH          = "~/.ssh/id_rsa"
	RemoteHostsPath       = "hosts"

	// run benchmark, otherwise run experiments
	RunBenchmark  = false
	BenchDuration = max(5*time.Minute, time.Duration(NodeCount/2))
	BenchLoads    = []int{10000}

	OnlySetupDocker  = true
	OnlySetupCluster = false
	OnlyRunCluster   = false
)

func getNodeConfig() node.Config {
	config := node.DefaultConfig
	config.Debug = true
	config.BroadcastTx = BroadcastTx
	if !CheckRotation {
		config.ConsensusConfig.ViewWidth = 24 * time.Hour
		config.ConsensusConfig.LeaderTimeout = 24 * time.Hour
	}
	return config
}

func setupExperiments() []Experiment {
	expms := make([]Experiment, 0)
	if RemoteLinuxCluster {
		expms = append(expms, &experiments.NetworkDelay{
			Delay: 100 * time.Millisecond,
		})
		expms = append(expms, &experiments.NetworkPacketLoss{
			Percent: 10,
		})
	}
	expms = append(expms, &experiments.MajorityKeepRunning{})
	expms = append(expms, &experiments.CorrectExecution{})
	expms = append(expms, &experiments.RestartCluster{})
	return expms
}

func main() {
	printAndCheckVars()
	os.Mkdir(WorkDir, 0755)
	buildPPoV()
	setupTransport()

	if RunBenchmark {
		runBenchmark()
	} else {
		var cfactory cluster.ClusterFactory
		if RemoteLinuxCluster {
			cfactory = makeRemoteClusterFactory()
		} else {
			cfactory = makeLocalClusterFactory()
		}

		if OnlySetupDocker {
			setupRapidDocker(cfactory)
			return
		}
		if OnlySetupCluster {
			setupRapidCluster(cfactory)
			return
		}
		testutil.NewLoadGenerator(makeLoadClient(), LoadTxPerSec, LoadJobPerTick)
		if OnlyRunCluster {
			runRapidCluster(cfactory)
			return
		}
		runExperiments(cfactory)
	}
}

func setupTransport() {
	// to make load test http client efficient
	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 100
}

func runBenchmark() {
	bm := &Benchmark{
		workDir:    path.Join(WorkDir, "benchmarks"),
		duration:   BenchDuration,
		interval:   5 * time.Second,
		cfactory:   makeRemoteClusterFactory(),
		loadClient: makeLoadClient(),
	}
	bm.Run()
}

func setupRapidDocker(cfactory cluster.ClusterFactory) {
	if cls, err := cfactory.SetupCluster("docker_template"); err == nil {
		dockerCompose := cluster.NewDockerCompose(cls)
		if err = cluster.WriteYamlFile(cfactory.TemplateDir(), dockerCompose); err == nil {
			fmt.Printf("docker-compose -f %s up -d\n", path.Join(cfactory.TemplateDir(), "docker-compose.yaml"))
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func setupRapidCluster(cfactory cluster.ClusterFactory) {
	if cls, err := cfactory.SetupCluster("cluster_template"); err == nil {
		fmt.Println("The cluster startup command is as follows.")
		for i := 0; i < cls.NodeCount(); i++ {
			fmt.Println(cls.GetNode(i).PrintCmd())
		}
	} else {
		fmt.Println(err)
	}
}

type KeepAliveRunning struct{}

func (expm *KeepAliveRunning) Name() string {
	return "keep_alive_running"
}

func (expm *KeepAliveRunning) Run(*cluster.Cluster) error {
	select {}
}

func runRapidCluster(cfactory cluster.ClusterFactory) {
	r := &ExperimentRunner{cfactory: cfactory}
	if err := r.runSingleExperiment(&KeepAliveRunning{}); err != nil {
		fmt.Println(err)
	}
}

func runExperiments(cfactory cluster.ClusterFactory) {
	r := &ExperimentRunner{
		experiments: setupExperiments(),
		cfactory:    cfactory,
	}
	pass, fail := r.Run()
	fmt.Printf("\nTotal: %d  |  Pass: %d  |  Fail: %d\n", len(r.experiments), pass, fail)
}

func printAndCheckVars() {
	fmt.Println("NodeCount =", NodeCount)
	fmt.Println("LoadJobPerTick =", LoadJobPerTick)
	fmt.Println("LoadSubmitNodes =", LoadSubmitNodes)
	fmt.Println("LoadBatchSubmit =", LoadBatchSubmit)
	fmt.Println("EmptyChainCode =", EmptyChainCode)
	fmt.Println("CheckRotation =", CheckRotation)
	fmt.Println("BroadcastTx =", BroadcastTx)
	fmt.Println("RemoteLinuxCluster =", RemoteLinuxCluster)
	fmt.Println("RunBenchmark =", RunBenchmark)
	fmt.Println("BenchLoads =", BenchLoads)
	fmt.Println("consensus.ExecuteTxFlag =", consensus.ExecuteTxFlag)
	fmt.Println("consensus.PreserveTxFlag =", consensus.PreserveTxFlag)
	fmt.Println("consensus.GenerateTxFlag =", consensus.GenerateTxFlag)
	fmt.Println("consensus.VoteBatchFlag =", consensus.VoteBatchFlag)
	fmt.Println()
	pass := true
	if !RunBenchmark && len(LoadSubmitNodes) != 0 {
		fmt.Println("!RunBenchmark =?=> len(LoadSubmitNodes)=0")
	}
	if RunBenchmark && !LoadBatchSubmit {
		fmt.Println("RunBenchmark =?=> LoadBatchSubmit")
	}
	if !consensus.ExecuteTxFlag && !EmptyChainCode {
		fmt.Println("!consensus.ExecuteTxFlag ===> EmptyChainCode")
		pass = false
	}
	if !RunBenchmark && !CheckRotation {
		fmt.Println("!RunBenchmark ===> CheckRotation")
	}
	if RunBenchmark && !RemoteLinuxCluster {
		fmt.Println("RunBenchmark ===> RemoteLinuxCluster")
		pass = false
	}
	if OnlySetupDocker && RemoteLinuxCluster {
		fmt.Println("OnlySetupDocker ===> !RemoteLinuxCluster")
		pass = false
	}
	if OnlySetupDocker && RunBenchmark {
		fmt.Println("OnlySetupDocker ===> !RunBenchmark")
		pass = false
	}
	if OnlySetupCluster && RunBenchmark {
		fmt.Println("OnlySetupCluster ===> !RunBenchmark")
		pass = false
	}
	if OnlyRunCluster && RunBenchmark {
		fmt.Println("OnlyRunCluster ===> !RunBenchmark")
		pass = false
	}
	if !consensus.ExecuteTxFlag && !RunBenchmark {
		fmt.Println("!consensus.ExecuteTxFlag ===> RunBenchmark")
	}
	if consensus.PreserveTxFlag && !RunBenchmark {
		fmt.Println("consensus.PreserveTxFlag ===> RunBenchmark")
	}
	if consensus.ExecuteTxFlag && consensus.PreserveTxFlag {
		fmt.Println("consensus.ExecuteTxFlag ===> !consensus.PreserveTxFlag")
		pass = false
	}
	if consensus.ExecuteTxFlag && consensus.GenerateTxFlag {
		fmt.Println("consensus.ExecuteTxFlag ===> !consensus.GenerateTxFlag")
		pass = false
	}
	if pass {
		fmt.Println("passed testing parameters validation")
	} else {
		os.Exit(1)
	}
	fmt.Println()
}

func buildPPoV() {
	cmd := exec.Command("go", "build", "../cmd/chain")
	if RemoteLinuxCluster {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=linux")
		fmt.Printf(" $ export %s\n", "GOOS=linux")
	}
	fmt.Printf(" $ %s\n\n", strings.Join(cmd.Args, " "))
	check(cmd.Run())
}

func makeLoadClient() testutil.LoadClient {
	fmt.Println("Preparing load client")
	if EmptyChainCode {
		return testutil.NewEmptyClient(LoadSubmitNodes)
	}
	var binccPath string
	if PPoVCoinBinCC {
		buildPPoVCoinBinCC()
		binccPath = "./ppovcoin"
	}
	mintAccounts := 100
	destAccounts := 10000 // increase dest accounts for benchmark
	return testutil.NewPPoVCoinClient(LoadSubmitNodes, mintAccounts, destAccounts, binccPath)
}

func buildPPoVCoinBinCC() {
	cmd := exec.Command("go", "build")
	cmd.Args = append(cmd.Args, "-ldflags", "-s -w")
	cmd.Args = append(cmd.Args, "../execution/bincc/ppovcoin")
	if RemoteLinuxCluster {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=linux")
		fmt.Printf(" $ export %s\n", "GOOS=linux")
	}
	fmt.Printf(" $ %s\n\n", strings.Join(cmd.Args, " "))
	check(cmd.Run())
}

func makeLocalClusterFactory() *cluster.LocalFactory {
	ftry, err := cluster.NewLocalFactory(cluster.LocalFactoryParams{
		BinPath:     "./chain",
		WorkDir:     path.Join(WorkDir, "local-clusters"),
		NodeCount:   NodeCount,
		WorkerCount: WorkerCount,
		VoterCount:  VoterCount,
		SetupDocker: OnlySetupDocker,
		NodeConfig:  getNodeConfig(),
	})
	check(err)
	return ftry
}

func makeRemoteClusterFactory() *cluster.RemoteFactory {
	ftry, err := cluster.NewRemoteFactory(cluster.RemoteFactoryParams{
		BinPath:         "./chain",
		WorkDir:         path.Join(WorkDir, "remote-clusters"),
		NodeCount:       NodeCount,
		WorkerCount:     WorkerCount,
		VoterCount:      VoterCount,
		NodeConfig:      getNodeConfig(),
		KeySSH:          RemoteKeySSH,
		HostsPath:       RemoteHostsPath,
		SetupRequired:   RemoteSetupRequired,
		InstallRequired: RemoteInstallRequired,
	})
	check(err)
	return ftry
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
