package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"
)

// Complete return true if decoding of message completed
func (msg *message) Complete() bool { return msg.End }

// Error Last error
func (msg *message) Error() error { return msg.Err }

// Direction Message direction
func (msg *message) Direction() MessageDirection { return msg.Dir }

// Create Date and time begin of decode message
func (msg *message) Create() time.Time { return msg.CreateTime }

// Command If the message contained a command, this function returns it
func (msg *message) Command() string { return msg.Cmd }

// ServiceCentreAddress Return service centre address
func (msg *message) ServiceCentreAddress() string { return msg.TpScaNumber }

// ServiceCentreType Return service centre address type
func (msg *message) ServiceCentreType() NumberType { return msg.TpScaType }

// ServiceCentreNumericPlan Return service centre numbering plan identifier
func (msg *message) ServiceCentreNumericPlan() NumberNumericPlan { return msg.TpScaNumericPlan }

// ServiceCentreTime Service centre time stamp
func (msg *message) ServiceCentreTime() time.Time { return msg.ServiceCentreTimeStamp }

// Type Return message type indicator (MTI)
func (msg *message) Type() SmsType { return msg.MtiSmsType }

// IsStatusReport Status report indication (TP-SRI)
func (msg *message) IsStatusReport() bool { return msg.MtiSmsType == TypeSmsStatusReport }

// Reply path (TP-RP) if =true-A response is requested
func (msg *message) IsReplyPath() bool { return msg.MtiReplyPath }

// OriginatingAddress Originating address
func (msg *message) OriginatingAddress() string { return msg.TpOaNumber }

// OriginatingAddressType Originating address type
func (msg *message) OriginatingAddressType() NumberType { return msg.TpOaType }

// OriginatingAddressNumericPlan Originating address numbering plan identifier
func (msg *message) OriginatingAddressNumericPlan() NumberNumericPlan { return msg.TpOaNumericPlan }

// ProtocolIdentifier Protocol identifier
func (msg *message) ProtocolIdentifier() uint8 { return msg.Pid }

// IsFlash Message is flash
func (msg *message) IsFlash() bool { return msg.DscFlash }

// IsEncode7Bit Message encoded as 7bit asci
func (msg *message) IsEncode7Bit() bool { return !msg.DscUSC2 }

// IsEncodeUSC2 Message encoded as UCS2 (UTF-16)
func (msg *message) IsEncodeUSC2() bool { return msg.DscUSC2 }

// Data Decoded message data
func (msg *message) Data() string { return msg.SmsData }

// DataParts The number of SMS (parts)
func (msg *message) DataParts() (ret int) {
	ret = 1
	if msg.MtiUdhiFound {
		ret = int(msg.UdhiNumberParts)
	}
	return
}

// DischargeTime Status report field TP-DT - Discharge Time
func (msg *message) DischargeTime() time.Time { return msg.TpDischargeTime }

// ReportStatus Status report field TP-ST
func (msg *message) ReportStatus() StatusReport { return msg.TpStType }

// MessageReference Status report field TP-MR
func (msg *message) MessageReference() uint8 { return msg.TpMr }
