package encoders // import "github.com/webdeskltd/pdu/encoders"

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import (
	"fmt"
)

// Encode encodes the given UTF-8 text into GSM 7-bit (3GPP TS 23.038) encoding with packing
func (e7bit *impl7bit) Encode(str string) []byte {
	raw7 := make([]byte, 0, len(str))
	for _, r := range str {
		i := gsmTable.Index(r)
		if i < 0 {
			b := gsmEscapes.to7Bit(r)
			if b != byte(_Unknown) {
				raw7 = append(raw7, _Esc, b)
			} else {
				raw7 = append(raw7, b)
			}
			continue
		}
		raw7 = append(raw7, byte(i))
	}
	return e7bit.Pack7Bit(raw7)
}

// Pack7Bit Pack 7bit slice
func (e7bit *impl7bit) Pack7Bit(raw7 []byte) []byte {
	var oct int   // current octet in pack7
	var bit uint8 // current bit in octet
	var b byte    // current byte in raw7
	var pack7 = make([]byte, blocks(len(raw7)*7, 8))
	var pack = func(out []byte, b byte, oct int, bit uint8) (int, uint8) {
		for i := uint8(0); i < 7; i++ {
			out[oct] |= b >> i & 1 << bit
			bit++
			if bit == 8 {
				oct++
				bit = 0
			}
		}
		return oct, bit
	}
	for i := range raw7 {
		b = raw7[i]
		oct, bit = pack(pack7, b, oct, bit)
	}
	// N.B. in order to not confuse 7 zero-bits with @
	// <CR> code is added to the packed bits.
	if 8-bit == 7 {
		oct, bit = pack(pack7, _CR, oct, bit)
	} else if bit == 0 && b == _CR {
		// and if data ends with <CR> on the octet boundary,
		// then we add an additional octet with <CR>. See (3GPP TS 23.038).
		pack7 = append(pack7, 0x00)
		oct, bit = pack(pack7, _CR, oct, bit)
	}
	return pack7
}

// DisplayPack Display pack as string
func (e7bit *impl7bit) DisplayPack(buf []byte) (out string) {
	for i := 0; i < len(buf)*8; i++ {
		b := buf[i/8]
		if i%8 == 0 {
			out += fmt.Sprintf("\n%02X:", b)
		}
		off := 7 - uint8(i%8)
		out += fmt.Sprintf("%4d", b>>off&1)
	}
	return out
}
