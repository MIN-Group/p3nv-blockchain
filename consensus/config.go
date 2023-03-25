// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package consensus

import "time"

const ExecuteTxFlag = false //set to false when benchmark test

type Config struct {
	ChainID int64

	// maximum tx count in a batch
	BatchTxLimit int

	// maximum batch count in a block
	BlockBatchLimit int

	//batch count in a batch vote
	VoteBatchLimit int

	// block creation delay if no transactions in the pool
	TxWaitTime time.Duration

	// maximum delay the leader waits for voting on a batch
	BatchWaitTime time.Duration

	// duration to wait to propose next block if leader cannot create qc
	ProposeTimeout time.Duration

	// duration to wait to propose next batch if leader cannot create qc
	BatchTimeout time.Duration

	// minimum delay between each block (i.e, it can define maximum block rate)
	BlockDelay time.Duration

	// view duration for a leader
	ViewWidth time.Duration

	// leader must create next qc within this duration
	LeaderTimeout time.Duration

	// path to save the benchmark log of the consensus algorithm (it will not be saved if blank)
	BenchmarkPath string
}

var DefaultConfig = Config{
	BatchTxLimit:    200,
	BlockBatchLimit: 4,
	VoteBatchLimit:  4,
	TxWaitTime:      1 * time.Second,
	BatchWaitTime:   3 * time.Second,
	ProposeTimeout:  1500 * time.Millisecond,
	BatchTimeout:    1000 * time.Millisecond,
	BlockDelay:      100 * time.Millisecond, // maximum block rate = 10 blk per sec
	ViewWidth:       60 * time.Second,
	LeaderTimeout:   20 * time.Second,
	BenchmarkPath:   "",
}
