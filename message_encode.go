package pdu // import "github.com/webdeskltd/pdu"

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/webdeskltd/pdu/encoders"
)

// Make SCA
func (msg *message) makeSca() (ret *bytes.Buffer) {
	var num uint64
	var buf []byte
	ret = bytes.NewBufferString(``)
	if msg.TpScaLen == 0 {
		ret.WriteString(fmt.Sprintf("%02d", 0))
		return
	}
	switch msg.TpScaType {
	case NumberTypeInternational:
		num, msg.Err = strconv.ParseUint(msg.TpScaNumber, 0, 64)
		if msg.Err != nil {
			return
		}
		buf = encoders.NewSemiOctet().Encode(num)
		ret.WriteString(fmt.Sprintf("%02d", len(buf)+1))
		ret.WriteString(hex.EncodeToString([]byte{msg.TpScaTypeSource}))
		ret.WriteString(hex.EncodeToString(buf))
	default:
		ret.WriteString(fmt.Sprintf("%02d", 0))
	}
	return
}

// Originating Address
func (msg *message) makeTpDa() (ret *bytes.Buffer) {
	var num uint64
	var buf []byte
	ret = bytes.NewBufferString(``)
	switch msg.TpOaType {
	case NumberTypeInternational:
		num, msg.Err = strconv.ParseUint(msg.TpOaNumber, 0, 64)
		if msg.Err != nil {
			return
		}
		ret.WriteString(hex.EncodeToString([]byte{msg.TpOaLen}))
		ret.WriteString(hex.EncodeToString([]byte{msg.TpOaTypeSource}))
		buf = encoders.NewSemiOctet().Encode(num)
		if len(buf) < int(msg.TpOaLen) {
			var nb []byte
			// first zero
			for i := 0; i < int(msg.TpOaLen)/2-len(buf)+(int(msg.TpOaLen)/2)%2; i++ {
				nb = append(nb, 0x0)
			}
			nb = append(nb, buf...)
			buf = nb
		}
		ret.WriteString(hex.EncodeToString(buf))
	default:
		msg.Err = ErrNoValudRecipientNumber
		return
	}
	return
}

// Header information, part one
func (msg *message) makeTpDuHead() (ret *bytes.Buffer) {
	var pduType uint8

	ret = bytes.NewBufferString(``)
	// TP-RP ignore
	// TP-UDHI data header present
	if msg.MtiUdhiFound {
		pduType = pduType ^ 1<<6
	}
	// TP-SRR
	if msg.MtiReplyPath {
		pduType = pduType ^ 1<<5
	}
	// TP-VPF off
	// TP-RD Reject duplicates
	if msg.TpRdRejectDuplicates {
		pduType = pduType ^ 1<<2
	}
	// TP-MTI direction
	pduType = pduType ^ 1<<0
	ret.WriteString(hex.EncodeToString([]byte{pduType}))

	// TP-MR set 00
	ret.WriteString(`00`)

	// TP-DA
	ret.WriteString(msg.makeTpDa().String())
	if msg.Err != nil {
		return
	}

	// TP-PID
	ret.WriteString(hex.EncodeToString([]byte{msg.Pid}))

	// TP-Data-Coding-Scheme
	if msg.DscUSC2 {
		msg.DcsSource = 0x08
	}
	if msg.DscFlash {
		msg.DcsSource = msg.DcsSource ^ 1<<4
	}
	ret.WriteString(hex.EncodeToString([]byte{msg.DcsSource}))

	return
}

// Make all UD with length
func (msg *message) makeUdAll() (ret *bytes.Buffer) {
	ret = bytes.NewBufferString(``)
	ret.WriteString(hex.EncodeToString([]byte{byte(msg.SmsDataLength)}))
	ret.WriteString(hex.EncodeToString(msg.SmsDataSource))
	return
}

// Make UD header
// UDH = UDHL + HL + IEI + IEDL + RefID + Parts + Part number
func (msg *message) makeUdh() {
	var size int
	msg.UdhiLength = 0x05
	// IEI
	if msg.DscUSC2 {
		msg.UdhiIei = 0x00
	} else {
		msg.UdhiIei = 0x08
	}
	// IEDL
	switch msg.UdhiIei {
	case 0x00:
		msg.UdhiIedl = 0x03
	case 0x08:
		msg.UdhiIedl = 0x04
		msg.UdhiLength++
	}
	// Size and parts
	size = _MaxBytes - int(msg.UdhiLength)
	msg.UdhiNumberParts = uint8(msg.SmsDataLength / size)
	if int(msg.UdhiNumberParts)*size < msg.SmsDataLength {
		msg.UdhiNumberParts++
	}
	// RefID
	rand.Seed(time.Now().UnixNano())
	msg.UdhiIedID = uint16(rand.Uint32())
}

// Return part of message
// UDH = UDHL + HL + IEI + IEDL + RefID + Parts + Part number
func (msg *message) getUdh(i uint8) (ret *bytes.Buffer) {
	var size int
	ret = bytes.NewBufferString(``)
	size = _MaxBytes - int(msg.UdhiLength) - 1
	msg.Lp = size * int(i)
	// UDHL
	if len(msg.SmsDataSource[msg.Lp:]) >= size {
		ret.WriteString(hex.EncodeToString([]byte{byte(size + int(msg.UdhiLength) + 1)}))
	} else {
		ret.WriteString(hex.EncodeToString([]byte{byte(len(msg.SmsDataSource[msg.Lp:]) + int(msg.UdhiLength) + 1)}))
	}
	// HL
	ret.WriteString(hex.EncodeToString([]byte{byte(msg.UdhiLength)}))
	// IEI
	ret.WriteString(hex.EncodeToString([]byte{byte(msg.UdhiIei)}))
	// IEDL
	ret.WriteString(hex.EncodeToString([]byte{byte(msg.UdhiIedl)}))
	// RefID
	if msg.UdhiIedl == 0x3 {
		ret.WriteString(hex.EncodeToString([]byte{
			byte(msg.UdhiIedID & 0xFF),
		}))
	} else {
		ret.WriteString(hex.EncodeToString([]byte{
			byte(msg.UdhiIedID >> 8),
			byte(msg.UdhiIedID & 0xFF),
		}))
	}
	// Parts Number of parts
	ret.WriteString(hex.EncodeToString([]byte{byte(msg.UdhiNumberParts)}))
	// Part number
	ret.WriteString(hex.EncodeToString([]byte{byte(i + 1)}))
	// Data
	if len(msg.SmsDataSource[msg.Lp:]) >= size {
		ret.WriteString(hex.EncodeToString(msg.SmsDataSource[msg.Lp : msg.Lp+size]))
	} else {
		ret.WriteString(hex.EncodeToString(msg.SmsDataSource[msg.Lp:]))
	}
	return
}

// Calculate CSA length in octets
func (msg *message) getScaLength() (ret int) {
	if msg.TpScaLen == 0 {
		ret = 1
	} else {
		//		ret = int(msg.TpScaLen) * 7 / 8
		//		if ret*8/7 < int(msg.TpScaLen) {
		//			ret++
		//		}
		ret = 8
	}
	return
}

func (msg *message) Encode() (ret []string) {
	var sca, duHead, out string
	var i uint8

	sca = strings.ToUpper(msg.makeSca().String())
	if msg.Err != nil {
		return
	}

	duHead = strings.ToUpper(msg.makeTpDuHead().String())
	if msg.Err != nil {
		return
	}

	// Single SMS
	if msg.SmsDataLength <= _MaxBytes {
		out = sca + duHead + strings.ToUpper(msg.makeUdAll().String())
		out = fmt.Sprintf("AT+CMGS=%d\r\n%s",
			len(out)/2-msg.getScaLength(),
			out,
		)
		ret = append(ret, out)
		return
	}
	msg.makeUdh()

	// Get all parts
	for i = 0; i < msg.UdhiNumberParts; i++ {
		out = sca + duHead + strings.ToUpper(msg.getUdh(i).String())
		out = fmt.Sprintf("AT+CMGS=%d\r\n%s",
			len(out)/2-msg.getScaLength(),
			out,
		)
		ret = append(ret, out)
	}

	return
}
