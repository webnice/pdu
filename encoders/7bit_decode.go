package encoders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
)

// Decode the given GSM 7-bit packed octet data (3GPP TS 23.038) into a UTF-8 encoded string
func (e7bit *impl7bit) Decode(octets []byte) (str string, err error) {
	var escaped bool
	var r rune
	var raw7 = e7bit.Unpack7Bit(octets)
	for _, b := range raw7 {
		if b > _Max {
			err = ErrUnexpectedByte
			return
		} else if escaped {
			r = gsmEscapes.from7Bit(b)
			escaped = false
		} else if b == _Esc {
			escaped = true
			continue
		} else {
			r = gsmTable.Rune(int(b))
		}
		str += string(r)
	}
	return
}

// Unpack7Bit Unpack 7bit slice
func (e7bit *impl7bit) Unpack7Bit(pack7 []byte) []byte {
	var sep byte  // current septet
	var bit uint8 // current bit in septet
	var raw7 = make([]byte, 0, len(pack7))
	for _, oct := range pack7 {
		for i := uint8(0); i < 8; i++ {
			sep |= oct >> i & 1 << bit
			bit++
			if bit == 7 {
				raw7 = append(raw7, sep)
				sep = 0
				bit = 0
			}
		}
	}
	if bytes.HasSuffix(raw7, crcr) || bytes.HasSuffix(raw7, cr) {
		raw7 = raw7[:len(raw7)-1]
	}
	return raw7
}
