// +build ignore

package main

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import (
	"github.com/webdeskltd/pdu"
	"runtime"
	"time"
)

var (
	messages = []string{
		`07919761989901F0040B919701119905F80000211062320150610CC8329BFD065DDF72363904`,
	}
)

func main() {
	Main()
	println("exit.")
	runtime.Gosched()
	runtime.GC()
	time.Sleep(time.Second)
}

func Main() {
	var pduDecoder pdu.Interface

	pduDecoder = pdu.New().Decoder(messageReceiver)
	defer pduDecoder.Done()

	for i := range messages {
		pduDecoder.Writer().Write([]byte(messages[i]))
	}

}

// Receive new messages
func messageReceiver(msg pdu.Message) {
	println("New message found")
	if msg.Error() != nil {
		println("Message error")
		println(msg.Error().Error())
		return
	}

	print(" SMSC:")
	println(msg.ServiceCentreAddress())
	print(" From:")
	println(msg.OriginatingAddress())
	print(" SMS (")
	print(msg.DataParts())
	print("):")
	println(msg.Data())
}
