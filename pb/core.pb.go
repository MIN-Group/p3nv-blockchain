// Copyright (C) 2021 Aung Maw
// Licensed under the GNU General Public License v3.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: core.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash       []byte      `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Height     uint64      `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	ParentHash []byte      `protobuf:"bytes,3,opt,name=parentHash,proto3" json:"parentHash,omitempty"`
	Proposer   []byte      `protobuf:"bytes,4,opt,name=proposer,proto3" json:"proposer,omitempty"`
	QuorumCert *QuorumCert `protobuf:"bytes,5,opt,name=quorumCert,proto3" json:"quorumCert,omitempty"`
	ExecHeight uint64      `protobuf:"varint,6,opt,name=execHeight,proto3" json:"execHeight,omitempty"`
	MerkleRoot []byte      `protobuf:"bytes,7,opt,name=merkleRoot,proto3" json:"merkleRoot,omitempty"`
	Timestamp  int64       `protobuf:"varint,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Signature  []byte      `protobuf:"bytes,9,opt,name=signature,proto3" json:"signature,omitempty"` // signature of proposer
	Batchs     []*Batch    `protobuf:"bytes,10,rep,name=batchs,proto3" json:"batchs,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Block) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *Block) GetProposer() []byte {
	if x != nil {
		return x.Proposer
	}
	return nil
}

func (x *Block) GetQuorumCert() *QuorumCert {
	if x != nil {
		return x.QuorumCert
	}
	return nil
}

func (x *Block) GetExecHeight() uint64 {
	if x != nil {
		return x.ExecHeight
	}
	return 0
}

func (x *Block) GetMerkleRoot() []byte {
	if x != nil {
		return x.MerkleRoot
	}
	return nil
}

func (x *Block) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Block) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Block) GetBatchs() []*Batch {
	if x != nil {
		return x.Batchs
	}
	return nil
}

type Batch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash            []byte           `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Proposer        []byte           `protobuf:"bytes,2,opt,name=proposer,proto3" json:"proposer,omitempty"`
	BatchQuorumCert *BatchQuorumCert `protobuf:"bytes,3,opt,name=batchQuorumCert,proto3" json:"batchQuorumCert,omitempty"`
	Timestamp       int64            `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Transactions    [][]byte         `protobuf:"bytes,5,rep,name=transactions,proto3" json:"transactions,omitempty"` // transaction hashes
	TxList          []*Transaction   `protobuf:"bytes,6,rep,name=txList,proto3" json:"txList,omitempty"`             //transaction list
	Signature       []byte           `protobuf:"bytes,7,opt,name=signature,proto3" json:"signature,omitempty"`       // signature of proposer
}

func (x *Batch) Reset() {
	*x = Batch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Batch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Batch) ProtoMessage() {}

func (x *Batch) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Batch.ProtoReflect.Descriptor instead.
func (*Batch) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{1}
}

func (x *Batch) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Batch) GetProposer() []byte {
	if x != nil {
		return x.Proposer
	}
	return nil
}

func (x *Batch) GetBatchQuorumCert() *BatchQuorumCert {
	if x != nil {
		return x.BatchQuorumCert
	}
	return nil
}

func (x *Batch) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Batch) GetTransactions() [][]byte {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Batch) GetTxList() []*Transaction {
	if x != nil {
		return x.TxList
	}
	return nil
}

func (x *Batch) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type BlockCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash          []byte         `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	ElapsedExec   float64        `protobuf:"fixed64,2,opt,name=elapsedExec,proto3" json:"elapsedExec,omitempty"`
	ElapsedMerkle float64        `protobuf:"fixed64,3,opt,name=elapsedMerkle,proto3" json:"elapsedMerkle,omitempty"`
	OldBlockTxs   [][]byte       `protobuf:"bytes,5,rep,name=oldBlockTxs,proto3" json:"oldBlockTxs,omitempty"`
	StateChanges  []*StateChange `protobuf:"bytes,6,rep,name=stateChanges,proto3" json:"stateChanges,omitempty"`
	LeafCount     []byte         `protobuf:"bytes,7,opt,name=leafCount,proto3" json:"leafCount,omitempty"`
	MerkleRoot    []byte         `protobuf:"bytes,8,opt,name=merkleRoot,proto3" json:"merkleRoot,omitempty"`
}

func (x *BlockCommit) Reset() {
	*x = BlockCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockCommit) ProtoMessage() {}

func (x *BlockCommit) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockCommit.ProtoReflect.Descriptor instead.
func (*BlockCommit) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{2}
}

func (x *BlockCommit) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *BlockCommit) GetElapsedExec() float64 {
	if x != nil {
		return x.ElapsedExec
	}
	return 0
}

func (x *BlockCommit) GetElapsedMerkle() float64 {
	if x != nil {
		return x.ElapsedMerkle
	}
	return 0
}

func (x *BlockCommit) GetOldBlockTxs() [][]byte {
	if x != nil {
		return x.OldBlockTxs
	}
	return nil
}

func (x *BlockCommit) GetStateChanges() []*StateChange {
	if x != nil {
		return x.StateChanges
	}
	return nil
}

func (x *BlockCommit) GetLeafCount() []byte {
	if x != nil {
		return x.LeafCount
	}
	return nil
}

func (x *BlockCommit) GetMerkleRoot() []byte {
	if x != nil {
		return x.MerkleRoot
	}
	return nil
}

type Signature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PubKey []byte `protobuf:"bytes,1,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	Value  []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Signature) Reset() {
	*x = Signature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Signature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Signature) ProtoMessage() {}

func (x *Signature) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Signature.ProtoReflect.Descriptor instead.
func (*Signature) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{3}
}

func (x *Signature) GetPubKey() []byte {
	if x != nil {
		return x.PubKey
	}
	return nil
}

func (x *Signature) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type QuorumCert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash  []byte       `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Signatures []*Signature `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (x *QuorumCert) Reset() {
	*x = QuorumCert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuorumCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuorumCert) ProtoMessage() {}

func (x *QuorumCert) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuorumCert.ProtoReflect.Descriptor instead.
func (*QuorumCert) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{4}
}

func (x *QuorumCert) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *QuorumCert) GetSignatures() []*Signature {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type BatchQuorumCert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchHash  []byte       `protobuf:"bytes,1,opt,name=batchHash,proto3" json:"batchHash,omitempty"`
	Signatures []*Signature `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (x *BatchQuorumCert) Reset() {
	*x = BatchQuorumCert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchQuorumCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchQuorumCert) ProtoMessage() {}

func (x *BatchQuorumCert) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchQuorumCert.ProtoReflect.Descriptor instead.
func (*BatchQuorumCert) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{5}
}

func (x *BatchQuorumCert) GetBatchHash() []byte {
	if x != nil {
		return x.BatchHash
	}
	return nil
}

func (x *BatchQuorumCert) GetSignatures() []*Signature {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type Vote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash []byte     `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Signature *Signature `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Vote) Reset() {
	*x = Vote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vote) ProtoMessage() {}

func (x *Vote) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vote.ProtoReflect.Descriptor instead.
func (*Vote) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{6}
}

func (x *Vote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *Vote) GetSignature() *Signature {
	if x != nil {
		return x.Signature
	}
	return nil
}

type BatchVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Batchs     []*Batch     `protobuf:"bytes,1,rep,name=batchs,proto3" json:"batchs,omitempty"`
	Signatures []*Signature `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (x *BatchVote) Reset() {
	*x = BatchVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchVote) ProtoMessage() {}

func (x *BatchVote) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchVote.ProtoReflect.Descriptor instead.
func (*BatchVote) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{7}
}

func (x *BatchVote) GetBatchs() []*Batch {
	if x != nil {
		return x.Batchs
	}
	return nil
}

func (x *BatchVote) GetSignatures() []*Signature {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash      []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	Nonce     int64  `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Sender    []byte `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
	CodeAddr  []byte `protobuf:"bytes,5,opt,name=codeAddr,proto3" json:"codeAddr,omitempty"`
	Input     []byte `protobuf:"bytes,6,opt,name=input,proto3" json:"input,omitempty"`
	Expiry    uint64 `protobuf:"varint,7,opt,name=expiry,proto3" json:"expiry,omitempty"` // expiry block height
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{8}
}

func (x *Transaction) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Transaction) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Transaction) GetNonce() int64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Transaction) GetSender() []byte {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *Transaction) GetCodeAddr() []byte {
	if x != nil {
		return x.CodeAddr
	}
	return nil
}

func (x *Transaction) GetInput() []byte {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *Transaction) GetExpiry() uint64 {
	if x != nil {
		return x.Expiry
	}
	return 0
}

type TxCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash        []byte  `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	BlockHash   []byte  `protobuf:"bytes,2,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	BlockHeight uint64  `protobuf:"varint,3,opt,name=blockHeight,proto3" json:"blockHeight,omitempty"`
	Error       string  `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	Elapsed     float64 `protobuf:"fixed64,5,opt,name=elapsed,proto3" json:"elapsed,omitempty"`
}

func (x *TxCommit) Reset() {
	*x = TxCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxCommit) ProtoMessage() {}

func (x *TxCommit) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxCommit.ProtoReflect.Descriptor instead.
func (*TxCommit) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{9}
}

func (x *TxCommit) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *TxCommit) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *TxCommit) GetBlockHeight() uint64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *TxCommit) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *TxCommit) GetElapsed() float64 {
	if x != nil {
		return x.Elapsed
	}
	return 0
}

type TxList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*Transaction `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *TxList) Reset() {
	*x = TxList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxList) ProtoMessage() {}

func (x *TxList) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxList.ProtoReflect.Descriptor instead.
func (*TxList) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{10}
}

func (x *TxList) GetList() []*Transaction {
	if x != nil {
		return x.List
	}
	return nil
}

type StateChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key           []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	PrevValue     []byte `protobuf:"bytes,3,opt,name=prevValue,proto3" json:"prevValue,omitempty"`
	TreeIndex     []byte `protobuf:"bytes,4,opt,name=treeIndex,proto3" json:"treeIndex,omitempty"`
	PrevTreeIndex []byte `protobuf:"bytes,5,opt,name=prevTreeIndex,proto3" json:"prevTreeIndex,omitempty"`
}

func (x *StateChange) Reset() {
	*x = StateChange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateChange) ProtoMessage() {}

func (x *StateChange) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateChange.ProtoReflect.Descriptor instead.
func (*StateChange) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{11}
}

func (x *StateChange) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *StateChange) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *StateChange) GetPrevValue() []byte {
	if x != nil {
		return x.PrevValue
	}
	return nil
}

func (x *StateChange) GetTreeIndex() []byte {
	if x != nil {
		return x.TreeIndex
	}
	return nil
}

func (x *StateChange) GetPrevTreeIndex() []byte {
	if x != nil {
		return x.PrevTreeIndex
	}
	return nil
}

var File_core_proto protoreflect.FileDescriptor

var file_core_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x70, 0x62, 0x22, 0xc8, 0x02, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x70,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x0a, 0x71, 0x75, 0x6f, 0x72, 0x75,
	0x6d, 0x43, 0x65, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74,
	0x52, 0x0a, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x65, 0x78, 0x65, 0x63, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x65, 0x78, 0x65, 0x63, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0a, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x70, 0x62, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x06, 0x62, 0x61, 0x74, 0x63, 0x68, 0x73,
	0x22, 0x89, 0x02, 0x0a, 0x05, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x12, 0x42, 0x0a, 0x0f, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52, 0x0f, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x22, 0x0a, 0x0c,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x0c, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x2c, 0x0a, 0x06, 0x74, 0x78, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x74, 0x78, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x83, 0x02, 0x0a,
	0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x45, 0x78, 0x65, 0x63, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x45, 0x78,
	0x65, 0x63, 0x12, 0x24, 0x0a, 0x0d, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x4d, 0x65, 0x72,
	0x6b, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x65, 0x6c, 0x61, 0x70, 0x73,
	0x65, 0x64, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x6c, 0x64, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x78, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0b, 0x6f,
	0x6c, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x78, 0x73, 0x12, 0x38, 0x0a, 0x0c, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x65, 0x61, 0x66, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x6c, 0x65, 0x61, 0x66, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x52, 0x6f, 0x6f, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x52, 0x6f,
	0x6f, 0x74, 0x22, 0x39, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x70, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x70, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x5e, 0x0a,
	0x0a, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x32, 0x0a, 0x0a, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x63, 0x0a,
	0x0f, 0x42, 0x61, 0x74, 0x63, 0x68, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x62, 0x61, 0x74, 0x63, 0x68, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x61, 0x74, 0x63, 0x68, 0x48, 0x61, 0x73, 0x68, 0x12, 0x32,
	0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x22, 0x56, 0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x30, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x67, 0x0a, 0x09, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x62, 0x61, 0x74, 0x63, 0x68,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70,
	0x62, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x06, 0x62, 0x61, 0x74, 0x63, 0x68, 0x73, 0x12,
	0x32, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x73, 0x22, 0xb7, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x64, 0x65, 0x41, 0x64, 0x64, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6f, 0x64, 0x65, 0x41, 0x64, 0x64, 0x72, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x22, 0x8e, 0x01,
	0x0a, 0x08, 0x54, 0x78, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1c,
	0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x22, 0x32,
	0x0a, 0x06, 0x54, 0x78, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x6c, 0x69,
	0x73, 0x74, 0x22, 0x97, 0x01, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72,
	0x65, 0x76, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70,
	0x72, 0x65, 0x76, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x72, 0x65, 0x65,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x72, 0x65,
	0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x54, 0x72,
	0x65, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x70,
	0x72, 0x65, 0x76, 0x54, 0x72, 0x65, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_proto_rawDescOnce sync.Once
	file_core_proto_rawDescData = file_core_proto_rawDesc
)

func file_core_proto_rawDescGZIP() []byte {
	file_core_proto_rawDescOnce.Do(func() {
		file_core_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_proto_rawDescData)
	})
	return file_core_proto_rawDescData
}

var file_core_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_core_proto_goTypes = []interface{}{
	(*Block)(nil),           // 0: core.pb.Block
	(*Batch)(nil),           // 1: core.pb.Batch
	(*BlockCommit)(nil),     // 2: core.pb.BlockCommit
	(*Signature)(nil),       // 3: core.pb.Signature
	(*QuorumCert)(nil),      // 4: core.pb.QuorumCert
	(*BatchQuorumCert)(nil), // 5: core.pb.BatchQuorumCert
	(*Vote)(nil),            // 6: core.pb.Vote
	(*BatchVote)(nil),       // 7: core.pb.BatchVote
	(*Transaction)(nil),     // 8: core.pb.Transaction
	(*TxCommit)(nil),        // 9: core.pb.TxCommit
	(*TxList)(nil),          // 10: core.pb.TxList
	(*StateChange)(nil),     // 11: core.pb.StateChange
}
var file_core_proto_depIdxs = []int32{
	4,  // 0: core.pb.Block.quorumCert:type_name -> core.pb.QuorumCert
	1,  // 1: core.pb.Block.batchs:type_name -> core.pb.Batch
	5,  // 2: core.pb.Batch.batchQuorumCert:type_name -> core.pb.BatchQuorumCert
	8,  // 3: core.pb.Batch.txList:type_name -> core.pb.Transaction
	11, // 4: core.pb.BlockCommit.stateChanges:type_name -> core.pb.StateChange
	3,  // 5: core.pb.QuorumCert.signatures:type_name -> core.pb.Signature
	3,  // 6: core.pb.BatchQuorumCert.signatures:type_name -> core.pb.Signature
	3,  // 7: core.pb.Vote.signature:type_name -> core.pb.Signature
	1,  // 8: core.pb.BatchVote.batchs:type_name -> core.pb.Batch
	3,  // 9: core.pb.BatchVote.signatures:type_name -> core.pb.Signature
	8,  // 10: core.pb.TxList.list:type_name -> core.pb.Transaction
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_core_proto_init() }
func file_core_proto_init() {
	if File_core_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Batch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockCommit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Signature); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuorumCert); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchQuorumCert); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vote); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchVote); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxCommit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_core_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateChange); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_proto_goTypes,
		DependencyIndexes: file_core_proto_depIdxs,
		MessageInfos:      file_core_proto_msgTypes,
	}.Build()
	File_core_proto = out.File
	file_core_proto_rawDesc = nil
	file_core_proto_goTypes = nil
	file_core_proto_depIdxs = nil
}