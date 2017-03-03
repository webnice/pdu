package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"
	uc "unicode"

	"gopkg.in/webnice/pdu.v1/encoders"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Set Service Centre Address number
func (pdu *impl) setTpSca(m *message, sca string) (err error) {
	if len(sca) == 0 {
		return
	}
	if sca[0] == '+' {
		sca = sca[1:]
	}
	m.TpScaType = NumberTypeInternational
	m.TpScaNumericPlan = NumericPlanInternational
	m.TpScaNumber = sca
	if !rexNumeric.MatchString(m.TpScaNumber) && len(m.TpScaNumber) > 0 {
		m.TpScaType = NumberTypeAlphanumeric
		m.TpScaNumericPlan = NumericPlanAlphanumeric
	}
	m.TpScaTypeSource = encodeNumberType(m.TpScaType, m.TpScaNumericPlan)
	m.TpScaLen = uint8(len(m.TpScaNumber))
	return
}

// Set Originating address number
func (pdu *impl) setTpOa(m *message, addr string) (err error) {
	if len(addr) == 0 {
		err = ErrNoValudRecipientNumber
		return
	}
	if addr[0] == '+' {
		addr = addr[1:]
	}
	m.TpOaNumber = addr
	switch rexNumeric.MatchString(m.TpOaNumber) {
	case true:
		m.TpOaType = NumberTypeInternational
		m.TpOaNumericPlan = NumericPlanInternational
	case false:
		m.TpOaType = NumberTypeAlphanumeric
		m.TpOaNumericPlan = NumericPlanAlphanumeric
	}
	m.TpOaTypeSource = encodeNumberType(m.TpOaType, m.TpOaNumericPlan)
	m.TpOaLen = uint8(len(m.TpOaNumber))
	return
}

// Check all symbols in message
func (pdu *impl) checkUSC2Symbols(message string) (ret bool) {
	var sym rune
	for _, sym = range message {
		if sym > uc.MaxASCII || !uc.IsPrint(sym) {
			ret = true
		}
	}
	return
}

// Encoder SMS encoder
func (pdu *impl) Encoder(inp Encode) (ret []string, err error) {
	var m *message

	m = new(message)
	m.Dir = DirectionOutgoing
	m.CreateTime = time.Now().In(time.Local)

	// Set Service Centre Address number
	if err = pdu.setTpSca(m, inp.Sca); err != nil {
		return
	}

	// Set flags
	m.DscUSC2 = inp.Ucs2
	m.DscFlash = inp.Flash
	m.MtiReplyPath = inp.StatusReportRequest
	m.TpRdRejectDuplicates = inp.RejectDuplicates

	// Set UTF-16 encode if found any non ascii symbol
	if !m.DscUSC2 {
		m.DscUSC2 = pdu.checkUSC2Symbols(inp.Message)
	}

	// Set originating address number
	if err = pdu.setTpOa(m, inp.Address); err != nil {
		return
	}

	// Encoding body of sms
	if m.DscUSC2 {
		var enc = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		var tnsf = enc.NewEncoder()
		m.SmsDataSource, m.SmsDataLength, err = transform.Bytes(tnsf, []byte(inp.Message))
		if err != nil {
			return
		}
		//m.SmsData = hex.EncodeToString(m.SmsDataSource)
		m.SmsDataLength = len(m.SmsDataSource)
	} else {
		m.SmsDataSource = encoders.New7Bit().Encode(inp.Message)
		m.SmsDataLength = len(inp.Message)
		//m.SmsData = hex.EncodeToString(m.SmsDataSource)
	}

	// Set multipart message
	if m.SmsDataLength > _MaxBytes {
		m.MtiUdhiFound = true
	}

	// Message encode
	ret = m.Encode()
	err = m.Err

	return
}
