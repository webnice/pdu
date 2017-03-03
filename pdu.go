// Decode and encode SMS in PDU format
// Articles describing format:
//   https://hiteshagja.wordpress.com/2010/04/04/send-long-sms/
//   http://hardisoft.ru/soft/samodelkin-soft/poluchenie-i-dekodirovanie-sms-soobshhenij-v-formate-pdu/
//   http://hardisoft.ru/soft/samodelkin-soft/otpravka-dlinnyx-sms-soobshhenij-v-formate-pdu/
//   http://www.smartposition.nl/resources/sms_pdu.html
//   https://geektimes.ru/post/257884/
//   http://www.lib.ru/unixhelp/modemmin.txt
//   http://www.netlab.linkpc.net/wiki/ru:hardware:huawei:e3272
//   http://alex-exe.ru/radio/wireless/gsm-sim900-at-command/

package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"container/list"
	"log"
	"runtime"
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

// Done Waiting for processing all incoming messages
func (pdu *impl) Done() { pdu.doCount.Wait() }

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
					log.Printf("Error decode: %s", err.Error())
				}
			}(buf)
			//		case <-tmr.C:
			//			pdu.CheckIncomleteMessages()
		}

	}
}
