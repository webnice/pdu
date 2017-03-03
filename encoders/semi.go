package encoders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// NewSemiOctet New semi-octet encoder object and return interface
func NewSemiOctet() EncodeSemiOctet {
	var semi = new(implsemi)
	return semi
}
