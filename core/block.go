// Copyright (C) 2021 Aung Maw
// Licensed under the GNU General Public License v3.0

package core

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/wooyang2018/ppov-blockchain/pb"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// errors
var (
	ErrInvalidBlockHash = errors.New("invalid block hash")
	ErrNilBlock         = errors.New("nil block")
)

// Block type
type Block struct {
	data         *pb.Block
	proposer     *PublicKey
	quorumCert   *QuorumCert
	batchs       []*Batch
	txList       *TxList
	transactions [][]byte
}

var _ json.Marshaler = (*Block)(nil)
var _ json.Unmarshaler = (*Block)(nil)

func NewBlock() *Block {
	return &Block{
		data: new(pb.Block),
	}
}

// Sum returns sha3 sum of block
func (blk *Block) Sum() []byte {
	h := sha3.New256()
	binary.Write(h, binary.BigEndian, blk.data.Height)
	h.Write(blk.data.ParentHash)
	h.Write(blk.data.Proposer)
	if blk.data.QuorumCert != nil {
		h.Write(blk.data.QuorumCert.BlockHash) // qc reference block hash
	}
	binary.Write(h, binary.BigEndian, blk.data.ExecHeight)
	h.Write(blk.data.MerkleRoot)
	binary.Write(h, binary.BigEndian, blk.data.Timestamp)
	for _, batch := range blk.data.Batchs {
		h.Write(batch.Hash)
	}
	return h.Sum(nil)
}

// Validate block
func (blk *Block) Validate(vs ValidatorStore) error {
	if blk.data == nil {
		return ErrNilBlock
	}
	if !blk.IsGenesis() { // skip quorum cert validation for genesis block
		if err := blk.quorumCert.Validate(vs); err != nil {
			return err
		}
		for _, batch := range blk.Batchs() {
			if err := batch.Validate(vs); err != nil {
				return err
			}
		}
	}
	if !bytes.Equal(blk.Sum(), blk.Hash()) {
		return ErrInvalidBlockHash
	}
	sig, err := newSignature(&pb.Signature{
		PubKey: blk.data.Proposer,
		Value:  blk.data.Signature,
	})
	if !vs.IsWorker(sig.PublicKey()) {
		return ErrInvalidValidator
	}
	if err != nil {
		return err
	}
	if !sig.Verify(blk.data.Hash) {
		return ErrInvalidSig
	}
	return nil
}

// Vote creates a vote for block
func (blk *Block) Vote(signer Signer) *Vote {
	vote := NewVote()
	vote.setData(&pb.Vote{
		BlockHash: blk.data.Hash,
		Signature: signer.Sign(blk.data.Hash).data,
	})
	return vote
}

func (blk *Block) ProposerVote() *Vote {
	vote := NewVote()
	vote.setData(&pb.Vote{
		BlockHash: blk.data.Hash,
		Signature: &pb.Signature{
			PubKey: blk.data.Proposer,
			Value:  blk.data.Signature,
		},
	})
	return vote
}

func (blk *Block) setData(data *pb.Block) error {
	blk.data = data
	if !blk.IsGenesis() { // every block contains qc except for genesis
		blk.quorumCert = NewQuorumCert()
		if err := blk.quorumCert.setData(data.QuorumCert); err != nil {
			return err
		}
		blk.batchs = make([]*Batch, len(data.Batchs))
		for index := range blk.batchs {
			blk.batchs[index] = NewBatch()
			if err := blk.batchs[index].setData(data.Batchs[index]); err != nil {
				return err
			}
		}
	}
	proposer, err := NewPublicKey(blk.data.Proposer)
	if err != nil {
		return err
	}
	blk.proposer = proposer
	return blk.setInternalData()
}

func (blk *Block) setInternalData() error {
	//使用集合对Batch中的交易进行去重
	txSet := make(map[string]struct{})
	txList := make([]*Transaction, 0)
	transactions := make([][]byte, 0)
	for _, batch := range blk.batchs {
		for _, tx := range *batch.TxList() {
			if _, ok := txSet[string(tx.Hash())]; !ok {
				txSet[string(tx.Hash())] = struct{}{}
				txList = append(txList, tx)
				transactions = append(transactions, tx.Hash())
			}
		}
	}
	res := TxList(txList)
	blk.txList = &res
	blk.transactions = transactions
	return nil
}

func (blk *Block) SetHeight(val uint64) *Block {
	blk.data.Height = val
	return blk
}

func (blk *Block) SetParentHash(val []byte) *Block {
	blk.data.ParentHash = val
	return blk
}

func (blk *Block) SetQuorumCert(val *QuorumCert) *Block {
	blk.quorumCert = val
	blk.data.QuorumCert = val.data
	return blk
}

func (blk *Block) SetExecHeight(val uint64) *Block {
	blk.data.ExecHeight = val
	return blk
}

func (blk *Block) SetMerkleRoot(val []byte) *Block {
	blk.data.MerkleRoot = val
	return blk
}

func (blk *Block) SetTimestamp(val int64) *Block {
	blk.data.Timestamp = val
	return blk
}

func (blk *Block) SetBatchs(val []*Batch) *Block {
	blk.data.Batchs = make([]*pb.Batch, len(val))
	blk.batchs = make([]*Batch, len(val))
	for index := range val {
		blk.batchs[index] = val[index]
		blk.data.Batchs[index] = val[index].data
	}
	blk.setInternalData()
	return blk
}

func (blk *Block) Sign(signer Signer) *Block {
	blk.proposer = signer.PublicKey()
	blk.data.Proposer = signer.PublicKey().key
	blk.data.Hash = blk.Sum()
	blk.data.Signature = signer.Sign(blk.data.Hash).data.Value
	return blk
}

func (blk *Block) Hash() []byte            { return blk.data.Hash }
func (blk *Block) Height() uint64          { return blk.data.Height }
func (blk *Block) ParentHash() []byte      { return blk.data.ParentHash }
func (blk *Block) Proposer() *PublicKey    { return blk.proposer }
func (blk *Block) QuorumCert() *QuorumCert { return blk.quorumCert }
func (blk *Block) ExecHeight() uint64      { return blk.data.ExecHeight }
func (blk *Block) MerkleRoot() []byte      { return blk.data.MerkleRoot }
func (blk *Block) Timestamp() int64        { return blk.data.Timestamp }
func (blk *Block) Batchs() []*Batch        { return blk.batchs }
func (blk *Block) IsGenesis() bool         { return blk.Height() == 0 }
func (blk *Block) Transactions() [][]byte  { return blk.transactions }
func (blk *Block) TxList() *TxList         { return blk.txList }

// Marshal encodes blk as bytes
func (blk *Block) Marshal() ([]byte, error) {
	return proto.Marshal(blk.data)
}

// UnmarshalBlock decodes block from bytes
func (blk *Block) Unmarshal(b []byte) error {
	data := new(pb.Block)
	if err := proto.Unmarshal(b, data); err != nil {
		return err
	}
	return blk.setData(data)
}

func (blk *Block) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(blk.data)
}

func (blk *Block) UnmarshalJSON(b []byte) error {
	data := new(pb.Block)
	if err := protojson.Unmarshal(b, data); err != nil {
		return err
	}
	return blk.setData(data)
}