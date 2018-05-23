package msg

import (
	"io"

	. "github.com/elastos/Elastos.ELA.Utility/common"
	"github.com/elastos/Elastos.ELA.Utility/p2p"
)

type GetBlocks struct {
	Locator  []*Uint256
	HashStop Uint256
}

func NewGetBlocks(locator []*Uint256, hashStop Uint256) *GetBlocks {
	msg := new(GetBlocks)
	msg.Locator = locator
	msg.HashStop = hashStop
	return msg
}

func (msg *GetBlocks) CMD() string {
	return p2p.CmdGetBlocks
}

func (msg *GetBlocks) Serialize(writer io.Writer) error {
	return WriteElements(writer, uint32(len(msg.Locator)), msg.Locator, msg.HashStop)
}

func (msg *GetBlocks) Deserialize(reader io.Reader) error {
	count, err := ReadUint32(reader)
	if err != nil {
		return err
	}

	msg.Locator = make([]*Uint256, count)
	return ReadElements(reader, &msg.Locator, &msg.HashStop)
}
