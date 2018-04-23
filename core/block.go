package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"time"

	. "github.com/elastos/Elastos.ELA.Utility/common"
	"github.com/elastos/Elastos.ELA.Utility/crypto"
)

const (
	FoundationAddress        = "8VYXVxKKSAxkmRrfmGpQR2Kc66XhG6m3ta"
	BlockVersion      uint32 = 0
	GenesisNonce      uint32 = 2083236893
	InvalidBlockSize         = -1
)

type Block struct {
	Header       *Header
	Transactions []*Transaction
}

func (b *Block) Serialize(w io.Writer) error {
	b.Header.Serialize(w)
	err := WriteUint32(w, uint32(len(b.Transactions)))
	if err != nil {
		return errors.New("Block item Transactions length serialization failed.")
	}

	for _, transaction := range b.Transactions {
		transaction.Serialize(w)
	}
	return nil
}

func (b *Block) Deserialize(r io.Reader) error {
	if b.Header == nil {
		b.Header = new(Header)
	}
	err := b.Header.Deserialize(r)
	if err != nil {
		return err
	}

	//Transactions
	var i uint32
	len, err := ReadUint32(r)
	if err != nil {
		return err
	}
	var tharray []Uint256
	for i = 0; i < len; i++ {
		transaction := new(Transaction)
		transaction.Deserialize(r)
		b.Transactions = append(b.Transactions, transaction)
		tharray = append(tharray, transaction.Hash())
	}

	return nil
}

func (b *Block) Trim(w io.Writer) error {
	b.Header.Serialize(w)
	err := WriteUint32(w, uint32(len(b.Transactions)))
	if err != nil {
		return errors.New("Block item Transactions length serialization failed.")
	}
	for _, transaction := range b.Transactions {
		temp := *transaction
		hash := temp.Hash()
		hash.Serialize(w)
	}
	return nil
}

func (b *Block) FromTrimmedData(r io.Reader) error {
	if b.Header == nil {
		b.Header = new(Header)
	}
	err := b.Header.Deserialize(r)
	if err != nil {
		return err
	}

	//Transactions
	var i uint32
	Len, err := ReadUint32(r)
	if err != nil {
		return err
	}
	var txhash Uint256
	var tharray []Uint256
	for i = 0; i < Len; i++ {
		txhash.Deserialize(r)
		b.Transactions = append(b.Transactions, NewTrimmedTx(txhash))
		tharray = append(tharray, txhash)
	}

	return nil
}

func (b *Block) GetSize() int {
	var buffer bytes.Buffer
	if err := b.Serialize(&buffer); err != nil {
		return InvalidBlockSize
	}

	return buffer.Len()
}

func (b *Block) Hash() Uint256 {
	return b.Header.Hash()
}

func (b *Block) GetArbitrators(arbiters []string) ([][]byte, error) {
	//todo finish this when arbitrator election scenario is done
	var arbitersByte [][]byte
	for _, arbiter := range arbiters {
		arbiterByte, err := HexStringToBytes(arbiter)
		if err != nil {
			return nil, err
		}
		arbitersByte = append(arbitersByte, arbiterByte)
	}

	return arbitersByte, nil
}

func GetGenesisBlock() (*Block, error) {
	// header
	header := &Header{
		Version:    BlockVersion,
		Previous:   EmptyHash,
		MerkleRoot: EmptyHash,
		Timestamp:  uint32(time.Unix(time.Date(2017, time.December, 22, 10, 0, 0, 0, time.UTC).Unix(), 0).Unix()),
		Bits:       0x1d03ffff,
		Nonce:      GenesisNonce,
		Height:     uint32(0),
	}

	// register ELA coin
	registerELACoin := &Transaction{
		TxType:         RegisterAsset,
		PayloadVersion: 0,
		Payload: &PayloadRegisterAsset{
			Asset: Asset{
				Name:      "ELA",
				Precision: 0x08,
				AssetType: 0x00,
			},
			Amount:     0 * 100000000,
			Controller: Uint168{},
		},
		Attributes: []*Attribute{},
		Inputs:     []*Input{},
		Outputs:    []*Output{},
		Programs:   []*Program{},
	}

	foundation, err := Uint168FromAddress(FoundationAddress)
	if err != nil {
		return nil, err
	}

	coinBase := NewCoinBaseTransaction(&PayloadCoinBase{}, 0)
	coinBase.Outputs = []*Output{
		{
			AssetID:     registerELACoin.Hash(),
			Value:       3300 * 10000 * 100000000,
			ProgramHash: *foundation,
		},
	}

	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, rand.Uint64())
	txAttr := NewAttribute(Nonce, nonce)
	coinBase.Attributes = append(coinBase.Attributes, &txAttr)
	//block
	block := &Block{
		Header:       header,
		Transactions: []*Transaction{coinBase, registerELACoin},
	}

	hashes := make([]Uint256, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		hashes = append(hashes, tx.Hash())
	}
	block.Header.MerkleRoot, err = crypto.ComputeRoot(hashes)
	if err != nil {
		return nil, errors.New("[GenesisBlock] , BuildMerkleRoot failed.")
	}

	return block, nil
}
