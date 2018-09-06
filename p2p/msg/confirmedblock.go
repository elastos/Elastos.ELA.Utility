package msg

import (
	"io"

	"github.com/elastos/Elastos.ELA.Utility/p2p"
)

type ConfirmedBlock struct {
	Proposal *DposProposalVoteSlot
}

func NewConfirmedBlock(proposal *DposProposalVoteSlot) *ConfirmedBlock {
	return &ConfirmedBlock{Proposal: proposal}
}

func (msg *ConfirmedBlock) CMD() string {
	return p2p.CmdConfirmedBlock
}

func (msg *ConfirmedBlock) Serialize(w io.Writer) error {
	return msg.Proposal.Serialize(w)
}

func (msg *ConfirmedBlock) Deserialize(r io.Reader) error {
	return msg.Proposal.Deserialize(r)
}
