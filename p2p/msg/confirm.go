package msg

import (
	"github.com/elastos/Elastos.ELA.Utility/common"
	"github.com/elastos/Elastos.ELA.Utility/p2p"
)

const MaxConfirmSize = 8276 //32+36*(33+32+65+33+65)

// Ensure Confirm implement p2p.Message interface.
var _ p2p.Message = (*Confirm)(nil)

type Confirm struct {
	common.Serializable
}

func NewConfirm(confirm common.Serializable) *Confirm {
	return &Confirm{Serializable: confirm}
}

func (msg *Confirm) CMD() string {
	return p2p.CmdConfirm
}

func (msg *Confirm) MaxLength() uint32 {
	return MaxConfirmSize
}
