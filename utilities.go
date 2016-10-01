package pdu // import "github.com/webdeskltd/pdu"

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import ()

// Decode number type and return numeric plan and number type
func decodeNumberType(nSrc int) (np NumberNumericPlan, nt NumberType) {
	var nPlan, nType = nSrc, nSrc >> 4
	nPlan, nType = nPlan&0x0F, nType&0x07
	switch nPlan {
	case 0x0:
		np = NumericPlanAlphanumeric
	case 0x1:
		np = NumericPlanInternational
	default:
		np = NumericPlanUnknown
	}
	switch nType {
	case 0x0:
		nt = NumberTypeUnknown
	case 0x1:
		nt = NumberTypeInternational
	case 0x2:
		nt = NumberTypeInternal
	case 0x3:
		nt = NumberTypeService
	case 0x4:
		nt = NumberTypeSubscriber
	case 0x5:
		nt = NumberTypeAlphanumeric
	case 0x6:
		nt = NumberTypeReduced
	case 0x7:
		nt = NumberTypeReserved
	}
	return
}
