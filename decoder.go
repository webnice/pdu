package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"container/list"
	"io"
	"time"
)

// Decoder Register function is invoked when decoding a new message
func (pdu *impl) Decoder(fn FnDecoder) Interface { pdu.DecFn = fn; return pdu }

// Writer Return writer
func (pdu *impl) Writer() io.Writer { return pdu }

// Write Writer implementation
func (pdu *impl) Write(in []byte) (l int, err error) {
	pdu.doCount.Add(1)
	pdu.Dec <- bytes.NewBuffer(in)
	l = len(in)
	return
}

// CheckIncomleteMessages Check stored messages
func (pdu *impl) CheckIncomleteMessages() {
	var m map[uint16]*countParts
	var ok bool
	var i uint16
	var elm *list.Element

	// Calculate
	m = make(map[uint16]*countParts)
	for elm = pdu.IncomleteMessages.Front(); elm != nil; elm = elm.Next() {
		var item = elm.Value.(*message)
		if _, ok = m[item.UdhiIedID]; ok {
			m[item.UdhiIedID].Count += 1
		} else {
			m[item.UdhiIedID] = &countParts{
				Count:       1,
				NumberParts: item.UdhiNumberParts,
			}
		}
	}

	// Check and paste together
	for i = range m {
		if m[i].Count == m[i].NumberParts {
			pdu.MessagePasteTogether(i)
		}
	}
}

// MessagePasteTogether Merge all messages with the specified ID
func (pdu *impl) MessagePasteTogether(id uint16) {
	var m *message
	var elm *list.Element
	var elms []*message
	var count, max uint8

	for elm = pdu.IncomleteMessages.Front(); elm != nil; elm = elm.Next() {
		var item = elm.Value.(*message)
		if item.UdhiIedID == id {
			elms = append(elms, item)
		}
	}

	// Find first
	for i := range elms {
		if elms[i].UdhiSequenceID == 1 && m == nil {
			m = elms[i]
			max = m.UdhiNumberParts
			count++
		}
	}

	// By order append sms body
	for {
		for i := range elms {
			if elms[i].UdhiSequenceID == count+1 {
				count++
				m.SmsData += elms[i].SmsData
				m.UdhiSequenceID = elms[i].UdhiSequenceID
				m.SmsDataLength += elms[i].SmsDataLength
				m.SmsDataSourceLength += elms[i].SmsDataSourceLength
			}
		}
		if count == max {
			break
		}
	}
	m.End = true
	pdu.DecFn(m)
}

// Decode source data to message
func (pdu *impl) Decode(src *bytes.Buffer) (err error) {
	var m *message
	defer func() {
		if e := recover(); err != nil {
			err = e.(error)
			m.Err = err
		}
	}()
	m = new(message)
	m.Dir = DirectionIncomming
	m.CreateTime = time.Now().In(time.Local)
	m.Scan(src)
	if !m.Complete() && m.Err == nil {
		pdu.IncomleteMessages.PushBack(m)
		pdu.CheckIncomleteMessages()
		return
	}
	pdu.DecFn(m)
	return
}
