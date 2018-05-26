package crypto

import (
	"fmt"
	"sort"
	"testing"
)

func TestSortPublicKeys(t *testing.T) {
	count := 10
	publicKeys := make([]*PublicKey, 0, count)
	dupPubKeys := make([]*PublicKey, 0, count)
	for i := 0; i < count; i++ {
		_, public, err := GenerateKeyPair()
		if err != nil {
			t.Errorf("Generate key pair failed, error %s", err.Error())
		}
		publicKeys = append(publicKeys, public)
		dupPubKeys = append(dupPubKeys, public)
	}

	SortPublicKeys(publicKeys)
	sort.Sort(pubKeySlice(dupPubKeys))

	for i, pubKey := range publicKeys {
		if !Equal(pubKey, dupPubKeys[i]) {
			t.Errorf("Sorted public keys not the same")
		}
	}

	fmt.Println(publicKeys)
	fmt.Println(dupPubKeys)
}

type pubKeySlice []*PublicKey

func (p pubKeySlice) Len() int { return len(p) }
func (p pubKeySlice) Less(i, j int) bool {
	r := p[i].X.Cmp(p[j].X)
	if r <= 0 {
		return true
	}
	return false
}
func (p pubKeySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
