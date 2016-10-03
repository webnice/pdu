package pdu // import "github.com/webdeskltd/pdu"

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import (
	"fmt"
)

var (
	// ErrIncorrectPDUdata Incorrect PDU data
	ErrIncorrectPDUdata = fmt.Errorf("Incorrect PDU data")
	// ErrNoValudRecipientNumber You must specify the valid recipient address
	ErrNoValudRecipientNumber = fmt.Errorf("You must specify the valid recipient address")
)
