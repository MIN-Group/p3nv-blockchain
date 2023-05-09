// Copyright (C) 2021 Aung Maw
// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package merkle

import (
	"crypto"
	"math/big"
)

// Position of a node in the tree
type Position struct {
	level uint8    //节点在树中的层次，父节点的层次大1
	index *big.Int //节点在层中的索引
	bytes []byte
}

// UnmarshalPosition unmarshals position from raw bytes
func UnmarshalPosition(b []byte) *Position {
	p := new(Position)
	p.bytes = b
	p.level = b[0]
	p.index = big.NewInt(0).SetBytes(b[1:])
	return p
}

// NewPosition create a new position
func NewPosition(level uint8, index *big.Int) *Position {
	p := new(Position)
	p.level = level
	p.index = index
	p.setBytes()
	return p
}

func (p *Position) setBytes() {
	ib := p.index.Bytes()
	if len(ib) == 0 {
		ib = []byte{0}
	}
	p.bytes = make([]byte, 0, 1+len(ib))
	p.bytes = append(p.bytes, p.level)
	p.bytes = append(p.bytes, ib...)
}

// Level gives the level of position
func (p *Position) Level() uint8 {
	return p.level
}

// Index gives the index of position
// NOTE: the value of index must not be changed
func (p *Position) Index() *big.Int {
	return p.index
}

// Bytes returns the serialized bytes of position
func (p *Position) Bytes() []byte {
	return p.bytes
}

func (p *Position) String() string {
	return string(p.bytes)
}

// Node type
type Node struct {
	Position *Position
	Data     []byte
}

// Group is a Group of child nodes under the same parent node
type Group struct {
	hashFunc       crypto.Hash
	tc             *TreeCalc
	store          Store
	parentPosition *Position
	nodes          []*Node
}

// NewGroup creates a new Group
func NewGroup(h crypto.Hash, tc *TreeCalc, store Store, pPos *Position) *Group {
	return &Group{
		hashFunc:       h,
		tc:             tc,
		store:          store,
		parentPosition: pPos,
		nodes:          make([]*Node, int(tc.BranchFactor())),
	}
}

// SetNode sets the node at the corresponding index in the block
func (b *Group) SetNode(n *Node) *Group {
	i := b.tc.NodeIndexInGroup(n.Position.Index()) //获取节点在group中的索引
	if i < len(b.nodes) {
		b.nodes[i] = n
	}
	return b
}

// Load loads the child nodes from the store
func (b *Group) Load(rowSize *big.Int) *Group {
	offset := b.tc.FirstNodeOfGroup(b.parentPosition.Index()) //根据父节点索引获取group的第一个节点的索引
	for i, n := range b.nodes {
		if n != nil {
			continue
		}
		index := big.NewInt(0).Add(offset, big.NewInt(int64(i)))
		if rowSize.Cmp(index) != 1 {
			break
		}
		p := NewPosition(b.parentPosition.level-1, index)

		if data := b.store.GetNode(p); data != nil {
			b.nodes[i] = &Node{p, data}
		}
	}
	return b
}

// MakeParent compute the sum of the child nodes and returns the parent node
func (b *Group) MakeParent() *Node {
	return &Node{
		Position: b.parentPosition,
		Data:     b.Sum(),
	}
}

// Sum sums the child nodes
func (b *Group) Sum() []byte {
	if b.IsEmpty() {
		return nil
	}
	h := b.hashFunc.New()
	for _, n := range b.nodes {
		if n != nil {
			h.Write(n.Data)
		}
	}
	return h.Sum(nil)
}

// IsEmpty checks whether all the child nodes are nil
func (b *Group) IsEmpty() bool {
	for _, n := range b.nodes {
		if n != nil {
			return false
		}
	}
	return true
}

// UpdateResult type
type UpdateResult struct {
	LeafCount *big.Int
	Height    uint8
	Leaves    []*Node
	Branches  []*Node
	Root      *Node
}
