package encoders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

// Decode numerical chunks from the given semi-octet encoded data
func (semi *implsemi) Decode(octets []byte) []int {
	var chunks = make([]int, 0, len(octets)*2)
	for _, oct := range octets {
		half := oct >> 4
		if half == 0xF {
			chunks = append(chunks, int(oct&0x0F))
			return chunks
		}
		chunks = append(chunks, int(oct&0x0F)*10+int(half))
	}
	return chunks
}

// DecodeAddress phone numbers from the given semi-octet encoded data
// This method is different from DecodeSemi because a 0x00 byte should be interpreted as
// two distinct digits. There 0x00 will be "00"
func (semi *implsemi) DecodeAddress(octets []byte) (str string) {
	for _, oct := range octets {
		half := oct >> 4
		if half == 0xF {
			str += fmt.Sprintf("%d", oct&0x0F)
			return
		}
		str += fmt.Sprintf("%d%d", oct&0x0F, half)
	}
	return
}
