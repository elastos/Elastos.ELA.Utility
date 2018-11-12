package msg

import (
	"io"

	"github.com/elastos/Elastos.ELA.Utility/p2p"
)

const MaxConfirmSize = 8276 //32+36*(33+32+65+33+65)

type Confirm struct {
	Command  string
	Proposal DPosProposalVoteSlot
}

func (msg *Confirm) CMD() string {
	return p2p.CmdConfirm
}

func (msg *Confirm) MaxLength() uint32 {
	return MaxConfirmSize
}

func (msg *Confirm) Serialize(w io.Writer) error {
	return msg.Proposal.Serialize(w)
}

func (msg *Confirm) Deserialize(r io.Reader) error {
	return msg.Proposal.Deserialize(r)
}
