package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"container/list"
	"io"
	"regexp"
	"sync"
	"time"
)

const _MaxBytes = int(140)

const (
	// TypeSmsDeliver Incomming SMS: SMS-DELIVER
	TypeSmsDeliver = SmsType(0x0)
	// TypeSmsStatusReport Incomming SMS: SMS-STATUS REPORT
	TypeSmsStatusReport = SmsType(0x2)
	// TypeSmsSubmitReport Incomming SMS: SMS-SUBMIT REPORT
	TypeSmsSubmitReport = SmsType(0x1)
	// TypeSmsReserved Incomming and outgoing SMS: RESERVED
	TypeSmsReserved = SmsType(0x3)
	// TypeSmsDeliverReport Outgoing SMS: SMS-DELIVER REPORT
	TypeSmsDeliverReport = SmsType(0x0f)
	// TypeSmsCommand Outgoing SMS: SMS-COMMAND
	TypeSmsCommand = SmsType(0x2f)
	// TypeSmsSubmit Outgoing SMS: SMS-SUBMIT
	TypeSmsSubmit = SmsType(0x1f)

	// DirectionIncomming Direction incomming
	DirectionIncomming = MessageDirection(`Incomming`)
	// DirectionOutgoing Direction outgoing
	DirectionOutgoing = MessageDirection(`Outgoing`)

	// NumberTypeUnknown Unknown. The cellular network does not know what the format of a number
	NumberTypeUnknown = NumberType(`Unknown. The cellular network does not know what the format of a number`)
	// NumberTypeInternational International number format
	NumberTypeInternational = NumberType(`International number format`)
	// NumberTypeInternal Internal number of the country. The prefixes of the country have no numbers
	NumberTypeInternal = NumberType(`Internal number of the country`)
	// NumberTypeService The Service network number. Used by the operator.
	NumberTypeService = NumberType(`The Service network number`)
	// NumberTypeSubscriber The subscriber's number. Used when a certain idea of short number stored in one or more of the SC as part of a high-level application
	NumberTypeSubscriber = NumberType(`The subscriber's number`)
	// NumberTypeAlphanumeric Alphanumeric encoded in 7-bit encoding
	NumberTypeAlphanumeric = NumberType(`Alphanumeric encoded in 7-bit encoding`)
	// NumberTypeReduced Reduced number
	NumberTypeReduced = NumberType(`Reduced number`)
	// NumberTypeReserved Reserved
	NumberTypeReserved = NumberType(`Reserved`)

	// NumericPlanAlphanumeric Alphanumeric encoded
	NumericPlanAlphanumeric = NumberNumericPlan(`Alphanumeric encoded`)
	// NumericPlanInternational International
	NumericPlanInternational = NumberNumericPlan(`International`)
	// NumericPlanUnknown Unknown
	NumericPlanUnknown = NumberNumericPlan(`Unknown`)

	// StatusDelivered Short message delivered successfully
	StatusDelivered = StatusReport(`Short message delivered successfully`)
	// StatusForwarded Forwarded, but status unknown
	StatusForwarded = StatusReport(`Forwarded, but status unknown`)
	// StatusReplaced Replaced
	StatusReplaced = StatusReport(`Replaced`)
	// StatusCongestion Congestion, still trying
	StatusCongestion = StatusReport(`Congestion, still trying`)
	// StatusRecipientBusy Recipient busy, still trying
	StatusRecipientBusy = StatusReport(`Recipient busy, still trying`)
	// StatusRecipientNoResponse No response recipient, still trying
	StatusRecipientNoResponse = StatusReport(`No response recipient, still trying`)
	// StatusServiceRejected Service rejected, still trying
	StatusServiceRejected = StatusReport(`Service rejected, still trying`)
	// StatusQosNotAvailable QOS not available, still trying
	StatusQosNotAvailableTrying = StatusReport(`QOS not available, still trying`)
	// StatusRecipientError Recipient error, still trying
	StatusRecipientError = StatusReport(`Recipient error, still trying`)
	// StatusRpcError RPC Error
	StatusRpcError = StatusReport(`RPC Error`)
	// StatusIncompatible Incompatible destination
	StatusIncompatible = StatusReport(`Incompatible destination`)
	// StatusConnectionRejected Connection rejected
	StatusConnectionRejected = StatusReport(`Connection rejected`)
	// StatusNotObtainable Not obtainable
	StatusNotObtainable = StatusReport(`Not obtainable`)
	// StatusQosNotAvailable QOS not available
	StatusQosNotAvailable = StatusReport(`QOS not available`)
	// StatusNoINAvailable No internetworking available
	StatusNoINAvailable = StatusReport(`No internetworking available`)
	// StatusMessageExpired Message expired
	StatusMessageExpired = StatusReport(`Message expired`)
	// StatusMessageDeletedBySender Message deleted by sender
	StatusMessageDeletedBySender = StatusReport(`Message deleted by sender`)
	// StatusMessageDeletedBySmsc Message deleted by SMSC
	StatusMessageDeletedBySmsc = StatusReport(`Message deleted by SMSC`)
	// StatusDoesNotExist Does not exist
	StatusDoesNotExist = StatusReport(`Does not exist`)
)

var (
	rexDataWithCommand    = regexp.MustCompile(`^\+([0-9A-Za-z]+)\: (\d+),([^,]*),(\d+)[\t\n\f\r ]+`)
	rexDataWithoutCommand = regexp.MustCompile(`([0-9A-Fa-f]+)$`)
	rexNumeric            = regexp.MustCompile(`^([0-9]+)$`)
)

// Interface is an interface
type Interface interface {
	// Done Waiting for processing all incoming messages
	Done()
	// Decoder Register function is invoked when decoding a new message
	Decoder(fn FnDecoder) Interface
	// Writer Return writer
	Writer() io.Writer
	// Encoder SMS encoder
	Encoder(Encode) ([]string, error)
}

// is an implementation
type impl struct {
	doCloseUp         chan bool          // Begin shutdown decoder goroutine
	doCloseDone       sync.WaitGroup     // Sync/wait when goroutine is running
	doCount           sync.WaitGroup     // Consideration received and processed messages
	Dec               chan *bytes.Buffer // Channel for decoder
	DecFn             FnDecoder          // Function call after new message decoded
	IncomleteMessages *list.List         // Temporary storage of partially received SMS messages
}

// Message SMS message
type Message interface {
	// Complete return true if decoding of message completed
	Complete() bool
	// Error Last error
	Error() error
	// Direction Message direction
	Direction() MessageDirection
	// Create Date and time begin of decode message
	Create() time.Time
	// Command If the message contained a command, this function returns it
	Command() string
	// ServiceCentreAddress Return service centre address
	ServiceCentreAddress() string
	// ServiceCentreType Return service centre address type
	ServiceCentreType() NumberType
	// ServiceCentreNumericPlan Return service centre numbering plan identifier
	ServiceCentreNumericPlan() NumberNumericPlan
	// ServiceCentreTime Service centre time stamp
	ServiceCentreTime() time.Time
	// Type Return message type indicator (MTI)
	Type() SmsType
	// IsStatusReport Status report indication (TP-SRI)
	IsStatusReport() bool
	// Reply path (TP-RP) if =true-A response is requested
	IsReplyPath() bool
	// OriginatingAddress Originating address
	OriginatingAddress() string
	// OriginatingAddressType Originating address type
	OriginatingAddressType() NumberType
	// OriginatingAddressNumericPlan Originating address numbering plan identifier
	OriginatingAddressNumericPlan() NumberNumericPlan
	// ProtocolIdentifier Protocol identifier
	ProtocolIdentifier() uint8
	// IsFlash Message is flash
	IsFlash() bool
	// IsEncode7Bit Message encoded as 7bit asci
	IsEncode7Bit() bool
	// IsEncodeUSC2 Message encoded as UCS2 (UTF-16)
	IsEncodeUSC2() bool
	// Data Decoded message data
	Data() string
	// DataParts The number of SMS (parts)
	DataParts() int
	// DischargeTime Status report field TP-DT - Discharge Time
	DischargeTime() time.Time
	// ReportStatus Status report field TP-ST
	ReportStatus() StatusReport

	//MessageReference in report message reference
	MessageReference() uint8
}

// Decoded sms message
type message struct {
	Dir                    MessageDirection  // Message direction
	End                    bool              // Decoding of message completed
	CreateTime             time.Time         // Date and time begin of decode message
	Err                    error             // Last error
	Lp                     int               // Last position
	Cmd                    string            // The command
	CmdStat                int64             // Command stat value
	CmdAlpha               string            // Command alpha value
	CmdLength              int64             // Command lengh value
	DataSource             []byte            // Source pdu data
	TpScaLen               uint8             // Length of the SMSC information
	TpScaTypeSource        uint8             // Type of address SMSC
	TpScaType              NumberType        // Type of SMSC number
	TpScaNumericPlan       NumberNumericPlan // SMSC Numbering plan identifier
	TpScaNumber            string            // Service Centre Address number
	MtiSource              byte              // Message Type indicator (MTI)
	MtiSmsType             SmsType           // MTI bits number 0, 1 - Message Type indicator (TP-MTI)
	MtiNoMoreMessageToSend bool              // MTI bit number 2 - More messages to send (TP-MMS). =true - No more messages to send in SC
	MtiStatusReport        bool              // MTI bit number 5 - Status report indication (TP-SRI) SRR
	MtiUdhiFound           bool              // MTI bit number 6 - TP-UDHI present. =true - User Data include User Data Header
	MtiReplyPath           bool              // MTI bit number 7 - Reply path (TP-RP). =true - A response is requested.
	TpMr                   uint8             // Status report field TP-MR
	TpRdRejectDuplicates   bool              // if true - reject duplicates (outgoing)
	TpOaLen                uint8             // Length of the Originating Address
	TpOaTypeSource         uint8             // Type of Originating Address
	TpOaType               NumberType        // Originating Address type
	TpOaNumericPlan        NumberNumericPlan // Originating Address numbering plan identifier
	TpOaNumber             string            // Originating Address
	Pid                    uint8             // Protocol identifier (TP-PID)
	DcsSource              uint8             // Data coding scheme
	DscFlash               bool              // Message is flash
	DscUSC2                bool              // Message encoded as UCS2
	ServiceCentreTimeStamp time.Time         // Service centre time stamp
	TpDischargeTime        time.Time         // Status report field TP-DT - Discharge Time
	SmsDataSourceLength    uint8             // User data length
	SmsDataSource          []byte            // User data source
	SmsDataLength          int               // Data length
	SmsData                string            // data body in string
	UdhiLength             uint8             // User data header length
	UdhiSource             []byte            // User data header as is
	UdhiIei                uint8             // User data header information element identifier
	UdhiIedl               uint8             // User data header information element length of the data
	UdhiIedID              uint16            // User data header message ID
	UdhiNumberParts        uint8             // User data header. Number of parts in the message
	UdhiSequenceID         uint8             // User data header. The sequence number of the message
	TpSt                   uint8             // Status report field TP-ST
	TpStType               StatusReport      // Status report field TP-ST as type
}

// FnDecoder Function call after new message decoded
type FnDecoder func(Message)

// MessageDirection Message direction
type MessageDirection string

// String Convert to string
func (md MessageDirection) String() string { return string(md) }

// SmsType Message Type indicator
type SmsType byte

// String Convert to string
func (st SmsType) String() (ret string) {
	switch st {
	case TypeSmsDeliver:
		ret = `Incomming SMS: SMS-DELIVER`
	case TypeSmsStatusReport:
		ret = `Incomming SMS: SMS-STATUS REPORT`
	case TypeSmsSubmitReport:
		ret = `Incomming SMS: SMS-SUBMIT REPORT`
	case TypeSmsReserved:
		ret = `Incomming and outgoing SMS: RESERVED`
	case TypeSmsDeliverReport:
		ret = `Outgoing SMS: SMS-DELIVER REPORT`
	case TypeSmsCommand:
		ret = `Outgoing SMS: SMS-COMMAND`
	case TypeSmsSubmit:
		ret = `Outgoing SMS: SMS-SUBMIT`
	}
	return
}

// NumberType Type of number
type NumberType string

// NumberNumericPlan Numbering plan identifier
type NumberNumericPlan string

type countParts struct {
	NumberParts uint8
	Count       uint8
}

// StatusReport description
type StatusReport string

// String Convert to string
func (sr StatusReport) String() string { return string(sr) }

// Encode Data structure to encoder message
type Encode struct {
	// Ucs2 Encode type. =true-Encode UTF16. =false-Encode 7bit
	Ucs2 bool
	// Flash =true-SMS is flash
	Flash bool
	// Smsc Service Centre Address number
	Sca string
	// Originating Address number
	Address string
	// Message data
	Message string
	// RejectDuplicates if true - reject duplicates
	RejectDuplicates bool
	// StatusReportRequest Status report request
	StatusReportRequest bool
}
