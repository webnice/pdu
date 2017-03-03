package encoders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// EncodeSemiOctet interface
type EncodeSemiOctet interface {
	// Decode numerical chunks from the given semi-octet encoded data
	Decode([]byte) []int
	// DecodeAddress phone numbers from the given semi-octet encoded data
	DecodeAddress([]byte) string
	// Packs the given numerical chunks in a semi-octet representation as described in 3GPP TS 23.040
	Encode(...uint64) []byte
}

type implsemi struct {
}
