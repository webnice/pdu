package pdu // import "github.com/webdeskltd/pdu"

//import "github.com/webdeskltd/debug"
import "github.com/webdeskltd/log"
import (
	"bytes"
	"container/list"
	"io"
	"runtime"
	"time"
)

// New Create new object and return interface
func New() Interface {
	var pdu = new(impl)
	pdu.doCloseUp = make(chan bool)
	pdu.Dec = make(chan *bytes.Buffer, 1000)
	pdu.IncomleteMessages = list.New()
	pdu.doCloseDone.Add(1)
	go func() {
		defer pdu.doCloseDone.Done()
		pdu.Worker()
	}()
	runtime.SetFinalizer(pdu, destructor)
	return pdu
}

// Object destructor
func destructor(pdu *impl) {
	pdu.doCloseUp <- true
	pdu.doCloseDone.Wait()
	close(pdu.doCloseUp)
	close(pdu.Dec)
}

// Writer Return writer
func (pdu *impl) Writer() io.Writer { return pdu }

// Write Writer implementation
func (pdu *impl) Write(in []byte) (l int, err error) {
	pdu.doCount.Add(1)
	pdu.Dec <- bytes.NewBuffer(in)
	l = len(in)
	return
}

// Done Waiting for processing all incoming messages
func (pdu *impl) Done() { pdu.doCount.Wait() }

// Decoder Register function is invoked when decoding a new message
func (pdu *impl) Decoder(fn FnDecoder) Interface { pdu.DecFn = fn; return pdu }

// Worker Goroudine read and decode messages
func (pdu *impl) Worker() {
	var err error
	var exit bool
	var buf *bytes.Buffer
	//	var tmr = time.NewTicker(time.Second)
	//	defer tmr.Stop()
	for !exit {
		select {
		case <-pdu.doCloseUp:
			exit = true
		case buf = <-pdu.Dec:
			// Parallel processing
			go func(b *bytes.Buffer) {
				defer pdu.doCount.Done()
				if err = pdu.Decode(b); err != nil {
					log.Errorf("Error decode: %s", err.Error())
				}
			}(buf)
			//		case <-tmr.C:
			//			pdu.CheckIncomleteMessages()
		}

	}
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
