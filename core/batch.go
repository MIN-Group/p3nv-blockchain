package core

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/wooyang2018/ppov-blockchain/logger"

	"github.com/wooyang2018/ppov-blockchain/pb"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// errors
var (
	ErrInvalidBatchHash = errors.New("invalid batch hash")
	ErrNilBatch         = errors.New("nil batch")
)

type Batch struct {
	data            *pb.Batch
	txList          *TxList
	proposer        *PublicKey
	batchQuorumCert *BatchQuorumCert
}

var _ json.Marshaler = (*Batch)(nil)
var _ json.Unmarshaler = (*Batch)(nil)

func NewBatch() *Batch {
	return &Batch{
		data: new(pb.Batch),
	}
}

// Sum returns sha3 sum of batch
func (b *Batch) Sum() []byte {
	h := sha3.New256()
	h.Write(b.data.Proposer)
	binary.Write(h, binary.BigEndian, b.data.Timestamp)
	for _, txHash := range b.data.Transactions {
		h.Write(txHash)
	}
	return h.Sum(nil)
}

// Validate batch
func (b *Batch) Validate(vs ValidatorStore) error {
	if b.data == nil {
		return ErrNilBatch
	}
	if b.batchQuorumCert != nil {
		if err := b.batchQuorumCert.Validate(vs); err != nil {
			return err
		}
	}
	if !bytes.Equal(b.Sum(), b.Hash()) {
		return ErrInvalidBatchHash
	}
	sig, err := newSignature(&pb.Signature{
		PubKey: b.data.Proposer,
		Value:  b.data.Signature,
	})
	if err != nil {
		return err
	}
	if !vs.IsWorker(sig.PublicKey()) {
		return ErrInvalidValidator
	}
	if !sig.Verify(b.data.Hash) {
		return ErrInvalidSig
	}
	return nil
}

func (b *Batch) setData(data *pb.Batch) error {
	b.data = data
	if data.BatchQuorumCert != nil && data.BatchQuorumCert.Signatures != nil {
		b.batchQuorumCert = NewBatchQuorumCert()
		if err := b.batchQuorumCert.setData(data.BatchQuorumCert); err != nil {
			return err
		}
	}

	proposer, err := NewPublicKey(b.data.Proposer)
	if err != nil {
		return err
	}
	b.proposer = proposer

	txs := make([]*Transaction, len(data.TxList), len(data.TxList))
	for i, v := range data.TxList {
		tx := NewTransaction()
		if err := tx.setData(v); err != nil {
			return err
		}
		txs[i] = tx
	}
	b.txList = (*TxList)(&txs)

	return nil
}

func (b *Batch) SetBatchQuorumCert(val *BatchQuorumCert) *Batch {
	if string(b.Hash()) == string(val.BatchHash()) {
		b.batchQuorumCert = val
		b.data.BatchQuorumCert = val.data
	}
	return b
}

func (b *Batch) SetTimestamp(val int64) *Batch {
	b.data.Timestamp = val
	return b
}

func (b *Batch) SetTransactions(val []*Transaction) *Batch {
	hashes := make([][]byte, len(val), len(val))
	data := make([]*pb.Transaction, len(val), len(val))
	txs := make([]*Transaction, len(val), len(val))
	for i, v := range val {
		hashes[i] = v.Hash()
		data[i] = v.data

		tx := NewTransaction()
		if err := tx.setData(v.data); err != nil {
			logger.I().Errorw("set transaction failed", "error", err)
		}
		txs[i] = tx
	}
	b.data.Transactions = hashes
	b.data.TxList = data
	b.txList = (*TxList)(&txs)
	return b
}

func (b *Batch) Sign(signer Signer) *Batch {
	b.proposer = signer.PublicKey()
	b.data.Proposer = signer.PublicKey().key
	b.data.Hash = b.Sum()
	b.data.Signature = signer.Sign(b.data.Hash).data.Value
	return b
}

func (b *Batch) Hash() []byte                      { return b.data.Hash }
func (b *Batch) Proposer() *PublicKey              { return b.proposer }
func (b *Batch) BatchQuorumCert() *BatchQuorumCert { return b.batchQuorumCert }
func (b *Batch) Timestamp() int64                  { return b.data.Timestamp }
func (b *Batch) Transactions() [][]byte            { return b.data.Transactions }
func (b *Batch) TxList() *TxList                   { return b.txList }

func (b *Batch) Marshal() ([]byte, error) {
	return proto.Marshal(b.data)
}

func (b *Batch) Unmarshal(bytes []byte) error {
	data := new(pb.Batch)
	if err := proto.Unmarshal(bytes, data); err != nil {
		return err
	}
	return b.setData(data)
}

func (b *Batch) UnmarshalJSON(bytes []byte) error {
	data := new(pb.Batch)
	if err := protojson.Unmarshal(bytes, data); err != nil {
		return err
	}
	return b.setData(data)
}

func (b *Batch) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(b.data)
}
