package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

var (
	// ErrIncorrectPDUdata Incorrect PDU data
	ErrIncorrectPDUdata = fmt.Errorf("Incorrect PDU data")
	// ErrNoValudRecipientNumber You must specify the valid recipient address
	ErrNoValudRecipientNumber = fmt.Errorf("You must specify the valid recipient address")
	// ErrEncodingNotImplementedForRecipientNumber Encoding not implemented for type of recipient address
	ErrEncodingNotImplementedForRecipientNumber = fmt.Errorf("Encoding not implemented for type of recipient address")
)
