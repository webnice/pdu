package encoders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// New7Bit New 7bit encoder object and return interface
func New7Bit() Encode7bit {
	var e7bit = new(impl7bit)
	return e7bit
}

func blocks(n, block int) int {
	if n%block == 0 {
		return n / block
	}
	return n/block + 1
}

func pad(n, block int) int {
	if n%block == 0 {
		return n
	}
	return (n/block + 1) * block
}

func (et *escapeTable) to7Bit(r rune) byte {
	for _, esc := range et {
		if esc.to == r {
			return esc.from
		}
	}
	return byte(_Unknown)
}

func (et *escapeTable) from7Bit(b byte) rune {
	for _, esc := range et {
		if esc.from == b {
			return esc.to
		}
	}
	return _Unknown
}

func (rt *runeTable) Index(r rune) int {
	for i := range rt {
		if rt[i] == r {
			return i
		}
	}
	return -1
}

func (rt *runeTable) Rune(idx int) rune {
	if idx >= 0 && idx < len(rt) {
		return rt[idx]
	}
	return _Unknown
}
