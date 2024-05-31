package consensus

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/wooyang2018/ppov-blockchain/core"
	"github.com/wooyang2018/ppov-blockchain/logger"
)

type voterState struct {
	batchQ         []*core.BatchHeader //待投票的Batch队列
	voteBatchLimit int
	mtxState       sync.RWMutex

	writer  *csv.Writer
	index   int
	preTime int64
	txCount int
}

func newVoterState() *voterState {
	return &voterState{
		batchQ: make([]*core.BatchHeader, 0),
	}
}

func (v *voterState) setVoteBatchLimit(voteBatchLimit int) *voterState {
	v.mtxState.Lock()
	defer v.mtxState.Unlock()
	v.voteBatchLimit = voteBatchLimit
	return v
}

func (v *voterState) addBatch(batch *core.BatchHeader, widx int, txs int) {
	t1 := time.Now().UnixNano()
	v.saveItem(widx, t1, txs)

	v.mtxState.Lock()
	defer v.mtxState.Unlock()
	if rand.Intn(100) < 95 {
		v.batchQ = append(v.batchQ, batch)
	}
	if PreserveTxFlag && len(v.batchQ) > 3*v.voteBatchLimit {
		v.batchQ = v.batchQ[1:]
	}
}

func (v *voterState) hasEnoughBatch() bool {
	v.mtxState.RLock()
	defer v.mtxState.RUnlock()
	return len(v.batchQ) >= v.voteBatchLimit
}

func (v *voterState) getBatchNum() int {
	v.mtxState.RLock()
	defer v.mtxState.RUnlock()
	return len(v.batchQ)
}

// popBatchHeaders 从队列头部弹出num个Batch
func (v *voterState) popBatchHeaders() []*core.BatchHeader {
	v.mtxState.Lock()
	defer v.mtxState.Unlock()

	num := v.voteBatchLimit
	res := make([]*core.BatchHeader, 0, num)
	for _, batch := range v.batchQ {
		if num <= 0 {
			break
		}
		num--
		res = append(res, batch)
	}
	v.batchQ = v.batchQ[len(res):]
	return res
}

func (v *voterState) newTester(file *os.File) {
	if file == nil {
		return
	}
	v.writer = csv.NewWriter(file)
	v.writer.Write([]string{
		"Index",
		"Proposer",
		"ReceiveTime",
		"TxCount",
		"Lambda",
	})
}

func (v *voterState) saveItem(widx int, t1 int64, txs int) {
	if v.writer == nil {
		return
	}
	if v.preTime == 0 {
		v.preTime = t1
	}
	v.txCount += txs

	var lambda float64
	if v.index%20 == 0 && v.index != 0 {
		elapsed := t1 - v.preTime
		lambda = float64(v.txCount) * 1e9 / float64(elapsed)
		v.preTime = t1
		v.txCount = 0
		v.writer.Flush()
	}

	v.writer.Write([]string{
		strconv.Itoa(v.index),
		strconv.Itoa(widx),
		strconv.FormatInt(t1, 10),
		strconv.Itoa(txs),
		strconv.FormatFloat(lambda, 'f', 2, 64),
	})
	v.index++
}

type leaderState struct {
	batchMap    map[string]*core.BatchHeader //batch hash -> batch
	batchSigns  map[string][]*core.Signature //batch hash -> signature list
	batchStopCh map[string]chan struct{}

	batchReadyQ []*core.BatchHeader //就绪Batch队列

	batchWaitTime   time.Duration //Batch超时时间
	blockBatchLimit int
	batchSignLimit  int

	mtxState sync.RWMutex //TODO 锁粒度优化
}

func newLeaderState() *leaderState {
	return &leaderState{
		batchMap:    make(map[string]*core.BatchHeader),
		batchSigns:  make(map[string][]*core.Signature),
		batchStopCh: make(map[string]chan struct{}),
		batchReadyQ: make([]*core.BatchHeader, 0),
	}
}

func (l *leaderState) setBatchSignLimit(batchSignLimit int) *leaderState {
	l.mtxState.Lock()
	defer l.mtxState.Unlock()
	l.batchSignLimit = batchSignLimit
	return l
}

func (l *leaderState) setBlockBatchLimit(blockBatchLimit int) *leaderState {
	l.mtxState.Lock()
	defer l.mtxState.Unlock()
	l.blockBatchLimit = blockBatchLimit
	return l
}

func (l *leaderState) setBatchWaitTime(batchWaitTime time.Duration) *leaderState {
	l.mtxState.Lock()
	defer l.mtxState.Unlock()
	l.batchWaitTime = batchWaitTime
	return l
}

func (l *leaderState) getBatchReadyNum() int {
	l.mtxState.Lock()
	defer l.mtxState.Unlock()
	return len(l.batchReadyQ)
}

func (l *leaderState) addBatchVote(vote *core.BatchVote) {
	l.mtxState.Lock()
	defer l.mtxState.Unlock()
	for index, sig := range vote.Signatures() {
		batch := vote.BatchHeaders()[index]
		hash := string(batch.Hash())
		//如果第一次收到对该Batch的投票
		if _, ok := l.batchMap[hash]; !ok {
			l.batchMap[hash] = batch
			l.batchSigns[hash] = make([]*core.Signature, 0)
			l.batchStopCh[hash] = make(chan struct{})
			go l.waitCleanState(hash, time.NewTimer(l.batchWaitTime), l.batchStopCh[hash])
		}
		l.batchSigns[hash] = append(l.batchSigns[hash], sig)
		if len(l.batchSigns[hash]) >= l.batchSignLimit {
			batchQC := core.NewBatchQuorumCert().Build(batch.Hash(), l.batchSigns[hash])
			batch.SetBatchQuorumCert(batchQC)
			l.batchReadyQ = append(l.batchReadyQ, batch)
			logger.I().Debugw("generated ready batch", "txs", len(batch.Transactions()))
			close(l.batchStopCh[hash])
			delete(l.batchMap, hash)
			delete(l.batchSigns, hash)
			delete(l.batchStopCh, hash)
		}
	}
}

// popReadyHeaders 从就绪队列的头部弹出num个Batch
func (l *leaderState) popReadyHeaders() []*core.BatchHeader {
	l.mtxState.Lock()
	defer l.mtxState.Unlock()

	num := l.blockBatchLimit
	res := make([]*core.BatchHeader, 0, num)
	for _, batch := range l.batchReadyQ {
		if num <= 0 { //TODO 改写法
			break
		}
		num--
		res = append(res, batch)
	}
	l.batchReadyQ = l.batchReadyQ[len(res):]
	return res
}

func (l *leaderState) waitCleanState(hash string, timer *time.Timer, stopCh chan struct{}) {
	select {
	case <-timer.C:
		l.mtxState.Lock()
		delete(l.batchMap, hash)
		delete(l.batchSigns, hash)
		delete(l.batchStopCh, hash)
		l.mtxState.Unlock()
		logger.I().Debug("deleted batch which is not ready")
	case <-stopCh:
	}
	timer.Stop()
}
