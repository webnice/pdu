package encoders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Packs the given numerical chunks in a semi-octet representation as described in 3GPP TS 23.040.
func (semi *implsemi) Encode(chunks ...uint64) []byte {
	var digits = make([]uint8, 0, len(chunks))
	for _, c := range chunks {
		var bucket []uint8
		if c < 10 {
			digits = append(digits, 0)
		}
		for c > 0 {
			d := c % 10
			bucket = append(bucket, uint8(d))
			c = (c - d) / 10
		}
		for i := range bucket {
			digits = append(digits, bucket[len(bucket)-1-i])
		}
	}
	octets := make([]byte, 0, len(digits)/2+1)
	for i := 0; i < len(digits); i += 2 {
		if len(digits)-i < 2 {
			octets = append(octets, 0xF0|digits[i])
			return octets
		}
		octets = append(octets, digits[i+1]<<4|digits[i])
	}
	return octets
}
