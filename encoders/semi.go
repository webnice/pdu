package encoders // import "github.com/webdeskltd/pdu/encoders"

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import ()

// NewSemiOctet New semi-octet encoder object and return interface
func NewSemiOctet() EncodeSemiOctet {
	var semi = new(implsemi)
	return semi
}
