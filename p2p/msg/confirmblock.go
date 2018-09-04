package msg

import (
	"io"

	"github.com/elastos/Elastos.ELA.Utility/common"
	"github.com/elastos/Elastos.ELA.Utility/p2p"
)

type ConfirmBlock struct {
	Hash common.Uint256
	Proposal DposProposalVoteSlot
}

func NewConfirmBlock(hash common.Uint256) *ConfirmBlock {
	return &ConfirmBlock{Hash: hash}
}

func (msg *ConfirmBlock) CMD() string {
	return p2p.CmdConfirmBlock
}

func (msg *ConfirmBlock) Serialize(w io.Writer) error {
	return msg.Proposal.Serialize(w)
}

func (msg *ConfirmBlock) Deserialize(r io.Reader) error {
	return msg.Proposal.Deserialize(r)
}
