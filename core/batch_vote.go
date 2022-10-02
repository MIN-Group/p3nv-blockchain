package core

import (
	"errors"

	"github.com/wooyang2018/ppov-blockchain/pb"
	"google.golang.org/protobuf/proto"
)

// errors
var (
	ErrNilBatchVote    = errors.New("nil batch vote")
	ErrDifferentSigner = errors.New("batch vote with differnt signers")
)

// BatchVote type
type BatchVote struct {
	data   *pb.BatchVote
	voter  *PublicKey
	batchs []*Batch
	sigs   sigList
}

func NewBatchVote() *BatchVote {
	return &BatchVote{
		data: new(pb.BatchVote),
	}
}

// Validate batch vote
func (vote *BatchVote) Validate(vs ValidatorStore) error {
	if vote.data == nil {
		return ErrNilBatchVote
	}
	for i, s := range vote.data.Signatures {
		sig, err := newSignature(s)
		if err != nil {
			return err
		}
		if !vs.IsVoter(sig.PublicKey()) {
			return ErrInvalidBatchVoter
		}
		if !sig.Verify(vote.data.Batchs[i].Hash) {
			return ErrInvalidSig
		}
	}

	return nil
}

func (vote *BatchVote) setData(data *pb.BatchVote) error {
	if data == nil {
		return ErrNilBatchVote
	}
	vote.data = data
	length := len(vote.data.Signatures)
	vote.sigs = make(sigList, 0, length)
	vote.batchs = make([]*Batch, 0, length)

	for i := 0; i < length; i++ {
		sig, err := newSignature(vote.data.Signatures[i])
		if err != nil {
			return err
		}
		vote.sigs = append(vote.sigs, sig)

		batch := NewBatch()
		if err := batch.setData(vote.data.Batchs[i]); err != nil {
			return err
		}
		vote.batchs = append(vote.batchs, batch)

		if vote.voter == nil {
			vote.voter = sig.pubKey
		} else {
			if vote.voter.String() != sig.PublicKey().String() {
				return ErrDifferentSigner
			}
		}

	}
	return nil
}

// Build generate a vote for a bunch of batches
func (vote *BatchVote) Build(batchs []*Batch, signer Signer) *BatchVote {
	data := new(pb.BatchVote)
	length := len(batchs)
	data.Batchs = make([]*pb.Batch, 0, length)
	data.Signatures = make([]*pb.Signature, 0, length)
	for i := 0; i < length; i++ {
		data.Batchs = append(data.Batchs, batchs[i].data)
		data.Signatures = append(data.Signatures, signer.Sign(batchs[i].data.Hash).data)
	}
	vote.setData(data)
	return vote
}

func (vote *BatchVote) Batchs() []*Batch         { return vote.batchs }
func (vote *BatchVote) Voter() *PublicKey        { return vote.voter }
func (vote *BatchVote) Signatures() []*Signature { return vote.sigs }

// Marshal encodes batch vote as bytes
func (vote *BatchVote) Marshal() ([]byte, error) {
	return proto.Marshal(vote.data)
}

// Unmarshal decodes batch vote from bytes
func (vote *BatchVote) Unmarshal(b []byte) error {
	data := new(pb.BatchVote)
	if err := proto.Unmarshal(b, data); err != nil {
		return err
	}
	return vote.setData(data)
}
