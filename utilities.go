package pdu

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import ()

// Decode number type and return numeric plan and number type
func decodeNumberType(nSrc uint8) (np NumberNumericPlan, nt NumberType) {
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

// Encode number type and return numeric plan and number type
func encodeNumberType(nt NumberType, np NumberNumericPlan) (ret uint8) {
	switch nt {
	case NumberTypeUnknown:
		ret = 0x8F
	case NumberTypeInternational:
		ret = 0x9F
	case NumberTypeInternal:
		ret = 0xAF
	case NumberTypeService:
		ret = 0xBF
	case NumberTypeSubscriber:
		ret = 0xCF
	case NumberTypeAlphanumeric:
		ret = 0xDF
	case NumberTypeReduced:
		ret = 0xEF
	case NumberTypeReserved:
		ret = 0xFF
	default:
		ret = 0x8F
	}
	switch np {
	case NumericPlanAlphanumeric:
		ret = ret & 0xF0
	case NumericPlanInternational:
		ret = ret & 0xF1
	default:
		ret = ret & 0xF0
	}
	return
}
