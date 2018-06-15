package p2p

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/elastos/Elastos.ELA.Utility/common"
)

const (
	CMDLEN      = 12
	CMDOFFSET   = 4
	CHECKSUMLEN = 4
	HEADERLEN   = 24
)

type header struct {
	Magic    uint32
	CMD      [CMDLEN]byte
	Length   uint32
	Checksum [CHECKSUMLEN]byte
}

func buildHeader(magic uint32, cmd string, body []byte) *header {
	// Calculate Checksum
	checksum := common.Sha256D(body)

	header := new(header)
	// Write Magic
	header.Magic = magic
	// Write CMD
	copy(header.CMD[:len(cmd)], cmd)
	// Write Length
	header.Length = uint32(len(body))
	// Write Checksum
	copy(header.Checksum[:], checksum[:CHECKSUMLEN])

	return header
}

func (header *header) Verify(buf []byte) error {
	// Verify Checksum
	sum := common.Sha256D(buf)
	checksum := sum[:CHECKSUMLEN]
	if !bytes.Equal(header.Checksum[:], checksum) {
		return fmt.Errorf("unmatched body checksum")
	}

	return nil
}

func (header *header) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (header *header) Deserialize(buf []byte) error {
	cmd := buf[CMDOFFSET : CMDOFFSET+CMDLEN]
	end := bytes.IndexByte(cmd, 0)
	if end < 0 || end >= CMDLEN {
		return errors.New("unexpected length of CMD")
	}

	hdr := bytes.NewReader(buf[:HEADERLEN])
	return binary.Read(hdr, binary.LittleEndian, header)
}

func (header *header) GetCMD() string {
	end := bytes.IndexByte(header.CMD[:], 0)
	return string(header.CMD[:end])
}
