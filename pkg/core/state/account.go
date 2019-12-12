package state

import (
	"github.com/CityOfZion/neo-go/pkg/crypto/keys"
	"github.com/CityOfZion/neo-go/pkg/io"
	"github.com/CityOfZion/neo-go/pkg/util"
)

// UnspentBalance contains input/output transactons that sum up into the
// account balance for the given asset.
type UnspentBalance struct {
	Tx    util.Uint256 `json:"txid"`
	Index uint16       `json:"n"`
	Value util.Fixed8  `json:"value"`
}

// UnspentBalances is a slice of UnspentBalance (mostly needed to sort them).
type UnspentBalances []UnspentBalance

// Account represents the state of a NEO account.
type Account struct {
	Version    uint8
	ScriptHash util.Uint160
	IsFrozen   bool
	Votes      []*keys.PublicKey
	Balances   map[util.Uint256][]UnspentBalance
}

// NewAccount returns a new Account object.
func NewAccount(scriptHash util.Uint160) *Account {
	return &Account{
		Version:    0,
		ScriptHash: scriptHash,
		IsFrozen:   false,
		Votes:      []*keys.PublicKey{},
		Balances:   make(map[util.Uint256][]UnspentBalance),
	}
}

// DecodeBinary decodes Account from the given BinReader.
func (s *Account) DecodeBinary(br *io.BinReader) {
	s.Version = uint8(br.ReadB())
	br.ReadBytes(s.ScriptHash[:])
	s.IsFrozen = br.ReadBool()
	br.ReadArray(&s.Votes)

	s.Balances = make(map[util.Uint256][]UnspentBalance)
	lenBalances := br.ReadVarUint()
	for i := 0; i < int(lenBalances); i++ {
		key := util.Uint256{}
		br.ReadBytes(key[:])
		ubs := make([]UnspentBalance, 0)
		br.ReadArray(&ubs)
		s.Balances[key] = ubs
	}
}

// EncodeBinary encodes Account to the given BinWriter.
func (s *Account) EncodeBinary(bw *io.BinWriter) {
	bw.WriteB(byte(s.Version))
	bw.WriteBytes(s.ScriptHash[:])
	bw.WriteBool(s.IsFrozen)
	bw.WriteArray(s.Votes)

	bw.WriteVarUint(uint64(len(s.Balances)))
	for k, v := range s.Balances {
		bw.WriteBytes(k[:])
		bw.WriteArray(v)
	}
}

// DecodeBinary implements io.Serializable interface.
func (u *UnspentBalance) DecodeBinary(r *io.BinReader) {
	u.Tx.DecodeBinary(r)
	u.Index = r.ReadU16LE()
	u.Value.DecodeBinary(r)
}

// EncodeBinary implements io.Serializable interface.
func (u *UnspentBalance) EncodeBinary(w *io.BinWriter) {
	u.Tx.EncodeBinary(w)
	w.WriteU16LE(u.Index)
	u.Value.EncodeBinary(w)
}

// GetBalanceValues sums all unspent outputs and returns a map of asset IDs to
// overall balances.
func (s *Account) GetBalanceValues() map[util.Uint256]util.Fixed8 {
	res := make(map[util.Uint256]util.Fixed8)
	for k, v := range s.Balances {
		balance := util.Fixed8(0)
		for _, b := range v {
			balance += b.Value
		}
		res[k] = balance
	}
	return res
}

// Len returns the length of UnspentBalances (used to sort things).
func (us UnspentBalances) Len() int { return len(us) }

// Less compares two elements of UnspentBalances (used to sort things).
func (us UnspentBalances) Less(i, j int) bool { return us[i].Value < us[j].Value }

// Swap swaps two elements of UnspentBalances (used to sort things).
func (us UnspentBalances) Swap(i, j int) { us[i], us[j] = us[j], us[i] }
