// Copyright (C) 2020 Aung Maw
// Licensed under the GNU General Public License v3.0

package txpool

import (
	"time"

	"github.com/wooyang2018/ppov-blockchain/core"
)

type broadcaster struct {
	msgSvc MsgService //通信服务

	queue     chan *core.Transaction //待广播交易的chan
	txBatch   []*core.Transaction    //待广播交易的切片
	batchSize int                    //当len(txBatch)>=batchSize时广播txBatch

	timeout time.Duration //当timeout超时广播txBatch
	timer   *time.Timer
}

func newBroadcaster(msgSvc MsgService) *broadcaster {
	b := &broadcaster{
		msgSvc:    msgSvc,
		queue:     make(chan *core.Transaction, 1000),
		batchSize: 100,
		timeout:   5 * time.Millisecond,
	}
	b.txBatch = make([]*core.Transaction, 0, b.batchSize)
	b.timer = time.NewTimer(b.timeout)
	go b.run()

	return b
}

func (b *broadcaster) run() {
	for {
		select {
		case <-b.timer.C:
			if len(b.txBatch) > 0 {
				b.broadcastBatch()
			}
			b.timer.Reset(b.timeout)

		case tx := <-b.queue:
			b.txBatch = append(b.txBatch, tx)
			if len(b.txBatch) >= b.batchSize {
				b.broadcastBatch()
			}
		}
	}
}

func (b *broadcaster) broadcastBatch() {
	b.msgSvc.BroadcastTxList((*core.TxList)(&b.txBatch))
	b.txBatch = make([]*core.Transaction, 0, b.batchSize)
}
