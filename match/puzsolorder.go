package match

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/mit-dci/opencx/crypto/timelockencoders"
)

// SolutionOrder is an order and modulus that are together.
// This includes an order, and the puzzle modulus factors.
type SolutionOrder struct {
	P *big.Int `json:"p"`
	Q *big.Int `json:"q"`
}

// NewSolutionOrder creates a new SolutionOrder from an already
// existing AuctionOrder, with a specified number of bits for an rsa
// key.
func NewSolutionOrder(rsaKeyBits uint64) (solOrder SolutionOrder, err error) {
	rsaKeyBitsInt := int(rsaKeyBits)

	// generate primes p and q
	var rsaPrivKey *rsa.PrivateKey
	if rsaPrivKey, err = rsa.GenerateMultiPrimeKey(rand.Reader, 2, rsaKeyBitsInt); err != nil {
		err = fmt.Errorf("Could not generate primes for RSA: %s", err)
		return
	}
	if len(rsaPrivKey.Primes) != 2 {
		err = fmt.Errorf("For some reason the RSA Privkey has != 2 primes, this should not be the case for RSW, we only need p and q")
		return
	}

	// finally set p, q, and the auction order.
	solOrder.P = new(big.Int).SetBytes(rsaPrivKey.Primes[0].Bytes())
	solOrder.Q = new(big.Int).SetBytes(rsaPrivKey.Primes[1].Bytes())
	return
}

// EncryptSolutionOrder encrypts a solution order and creates a puzzle
// along with the encrypted order
func (so *SolutionOrder) EncryptSolutionOrder(auctionOrder AuctionOrder, t uint64) (encSolOrder EncryptedSolutionOrder, err error) {
	// Try serializing the solution order
	var rawSolOrder []byte = auctionOrder.Serialize()
	if encSolOrder.OrderCiphertext, encSolOrder.OrderPuzzle, err = timelockencoders.CreateRC5RSWPuzzleWithPrimes(uint64(2), t, rawSolOrder, so.P, so.Q); err != nil {
		err = fmt.Errorf("Error creating puzzle from auction order: %s", err)
		return
	}

	// make sure they match
	encSolOrder.IntendedAuction = auctionOrder.AuctionID
	encSolOrder.IntendedPair = auctionOrder.TradingPair
	return
}

// Serialize uses gob encoding to turn the solution order into bytes.
func (so *SolutionOrder) Serialize() (raw []byte, err error) {
	var b bytes.Buffer

	// register SolutionOrder interface
	gob.Register(SolutionOrder{})

	// create a new encoder writing to the buffer
	enc := gob.NewEncoder(&b)

	// encode the puzzle in the buffer
	if err = enc.Encode(so); err != nil {
		err = fmt.Errorf("Error encoding solutionorder: %s", err)
		return
	}

	// Get the bytes from the buffer
	raw = b.Bytes()
	return
}

// Deserialize turns the solution order from bytes into a usable
// struct.
func (so *SolutionOrder) Deserialize(raw []byte) (err error) {
	var b *bytes.Buffer
	b = bytes.NewBuffer(raw)

	// register SolutionOrder
	gob.Register(SolutionOrder{})

	// create a new decoder writing to the buffer
	dec := gob.NewDecoder(b)

	// decode the solutionorder in the buffer
	if err = dec.Decode(so); err != nil {
		err = fmt.Errorf("Error decoding solutionorder: %s", err)
		return
	}

	return
}
