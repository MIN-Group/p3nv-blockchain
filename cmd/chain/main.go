// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/wooyang2018/ppov-blockchain/node"
)

const (
	FlagDebug   = "debug"
	FlagDataDir = "dataDir"

	FlagPort    = "port"
	FlagAPIPort = "apiPort"

	FlagBroadcastTx = "broadcast-tx"

	// storage
	FlagMerkleBranchFactor = "storage-merkleBranchFactor"

	// execution
	FlagTxExecTimeout       = "execution-txExecTimeout"
	FlagExecConcurrentLimit = "execution-concurrentLimit"

	// consensus
	FlagChainID         = "chainID"
	FlagBatchTxLimit    = "consensus-batchTxLimit"
	FlagBlockBatchLimit = "consensus-blockBatchLimit"
	FlagVoteBatchLimit  = "consensus-voteBatchLimit"
	FlagTxWaitTime      = "consensus-txWaitTime"
	FlagBatchWaitTime   = "consensus-batchWaitTime"
	FlagProposeTimeout  = "consensus-proposeTimeout"
	FlagBatchTimeout    = "consensus-batchTimeout"
	FlagBlockDelay      = "consensus-blockDelay"
	FlagViewWidth       = "consensus-viewWidth"
	FlagLeaderTimeout   = "consensus-leaderTimeout"
	FlagBenchmarkPath   = "consensus-benchmarkPath"
)

var nodeConfig = node.DefaultConfig

var rootCmd = &cobra.Command{
	Use:   "chain",
	Short: "ppov blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		node.Run(nodeConfig)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&nodeConfig.Debug,
		FlagDebug, false, "debug mode")

	rootCmd.PersistentFlags().StringVarP(&nodeConfig.Datadir,
		FlagDataDir, "d", "", "blockchain data directory")
	rootCmd.MarkPersistentFlagRequired(FlagDataDir)

	rootCmd.Flags().IntVarP(&nodeConfig.Port,
		FlagPort, "p", nodeConfig.Port, "p2p port")

	rootCmd.Flags().IntVarP(&nodeConfig.APIPort,
		FlagAPIPort, "P", nodeConfig.APIPort, "node api port")

	rootCmd.Flags().BoolVar(&nodeConfig.BroadcastTx,
		FlagBroadcastTx, false, "whether to broadcast transaction")

	rootCmd.Flags().Uint8Var(&nodeConfig.StorageConfig.MerkleBranchFactor,
		FlagMerkleBranchFactor, nodeConfig.StorageConfig.MerkleBranchFactor,
		"merkle tree branching factor")

	rootCmd.Flags().DurationVar(&nodeConfig.ExecutionConfig.TxExecTimeout,
		FlagTxExecTimeout, nodeConfig.ExecutionConfig.TxExecTimeout,
		"tx execution timeout")

	rootCmd.Flags().IntVar(&nodeConfig.ExecutionConfig.ConcurrentLimit,
		FlagExecConcurrentLimit, nodeConfig.ExecutionConfig.ConcurrentLimit,
		"concurrent tx execution limit")

	rootCmd.Flags().Int64Var(&nodeConfig.ConsensusConfig.ChainID,
		FlagChainID, nodeConfig.ConsensusConfig.ChainID,
		"chainid is used to create genesis block")

	rootCmd.Flags().IntVar(&nodeConfig.ConsensusConfig.BatchTxLimit,
		FlagBatchTxLimit, nodeConfig.ConsensusConfig.BatchTxLimit,
		"maximum tx count in a batch")

	rootCmd.Flags().IntVar(&nodeConfig.ConsensusConfig.BlockBatchLimit,
		FlagBlockBatchLimit, nodeConfig.ConsensusConfig.BlockBatchLimit,
		"maximum batch count in a block")

	rootCmd.Flags().IntVar(&nodeConfig.ConsensusConfig.VoteBatchLimit,
		FlagVoteBatchLimit, nodeConfig.ConsensusConfig.VoteBatchLimit,
		"batch count in a batch vote")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.TxWaitTime,
		FlagTxWaitTime, nodeConfig.ConsensusConfig.TxWaitTime,
		"block creation delay if no transactions in the pool")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.BatchWaitTime,
		FlagBatchWaitTime, nodeConfig.ConsensusConfig.BatchWaitTime,
		"maximum delay the leader waits for voting on a batch")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.ProposeTimeout,
		FlagProposeTimeout, nodeConfig.ConsensusConfig.ProposeTimeout,
		"duration to wait to propose next block if leader cannot create qc")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.BatchTimeout,
		FlagBatchTimeout, nodeConfig.ConsensusConfig.BatchTimeout,
		"duration to wait to propose next batch if leader cannot create qc")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.BlockDelay,
		FlagBlockDelay, nodeConfig.ConsensusConfig.BlockDelay,
		"minimum delay between blocks")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.ViewWidth,
		FlagViewWidth, nodeConfig.ConsensusConfig.ViewWidth,
		"view duration for a leader")

	rootCmd.Flags().DurationVar(&nodeConfig.ConsensusConfig.LeaderTimeout,
		FlagLeaderTimeout, nodeConfig.ConsensusConfig.LeaderTimeout,
		"leader must create next qc in this duration")

	rootCmd.Flags().StringVar(&nodeConfig.ConsensusConfig.BenchmarkPath,
		FlagBenchmarkPath, nodeConfig.ConsensusConfig.BenchmarkPath,
		"path to save the benchmark log of the consensus algorithm")
}
