// +build ignore

package main

//import "github.com/webdeskltd/debug"
//import "github.com/webdeskltd/log"
import (
	"fmt"
	"runtime"
	"time"

	"github.com/webdeskltd/pdu"
)

var (
	sms = []pdu.Encode{
		pdu.Encode{
			Sca:                 "+79043490000", // Tele2
			Ucs2:                false,
			Flash:               true,
			Address:             "+00000000000",
			Message:             "Hello world!",
			StatusReportRequest: true,
		},

		//		pdu.Encode{
		//			Sca:                 "+79168999100",
		//			Ucs2:                true,
		//			Flash:               false,
		//			Address:             "+00000000000",
		//			Message:             "Что делать, если хотел вытереть своей девушке слезы, но случайно стер брови?",
		//			StatusReportRequest: false,
		//		},
		//
		//		pdu.Encode{
		//			//			Sca:                 "+79168999100",
		//			Ucs2:    true,
		//			Flash:   false,
		//			Address:             "+00000000000",
		//			Message: `- Девушка, а почему у Вас ноги такие кривые?
		//- А это чтобы ты уши не натёр.`,
		//			StatusReportRequest: false,
		//		},

		//		pdu.Encode{
		//			Sca:     "+79043490000", // Tele2
		//			Ucs2:    false,
		//			Flash:   false,
		//			Address: "+00000000000",
		//			Message: `Just like a cold noreaster
		//			At first she'll sting,
		//			And then a single salty tear
		//			The heart will wring.
		//
		//			The evil heart will pity
		//			Something and then regret.
		//			But this light-headed sadness
		//			It will not forget.
		//
		//			I only sow. To harvest.
		//			Others will come. And yes!
		//			The lovely group of harvesters
		//			May true God bless.
		//
		//			And that more perfectly I could
		//			Give to you gratitude,
		//			Allow me to give the world
		//			Love incorruptible.`,
		//		},

		//		pdu.Encode{
		//			Sca:     "+79168999100",
		//			Ucs2:    true,
		//			Flash:   true,
		//			Address:             "+00000000000",
		//			Message: `Ночь, улица, фонарь, аптека,
		//				Бессмысленный и тусклый свет.
		//				Живи еще хоть четверть века -
		//				Все будет так. Исхода нет.
		//
		//				Умрешь - начнешь опять сначала
		//				И повторится все, как встарь:
		//				Ночь, ледяная рябь канала,
		//				Аптека, улица, фонарь.`,
		//		},

		//		pdu.Encode{
		//			Sca:     "+79168999100",
		//			Ucs2:    true,
		//			Flash:   false,
		//			Address:             "+00000000000",
		//			Message: `Ночь, улица, фонарь, аптека,
		//		Я покупаю вазелин.
		//		За мной стоят 2 человека:
		//		Армян и сумрачный грузин.
		//
		//		Вот скрипнула в подъезд пружина.
		//		И повторилось все как встарь:
		//		Пустая банка вазелина,
		//		аптека, улица, "фонарь".`,
		//		},

		//		pdu.Encode{
		//			//			Sca:     "+79043490000",
		//			Ucs2:    false,
		//			Flash:   false,
		//			Address:             "+00000000000",
		//			Message: "Noch'. Ulica. Fonar'. Apteka. Ja pokupaju vazelin. Za mnoj stojat dva cheloveka: armjan i sumrachnyj gruzin. Vot skripnula v pod#ezd pruzhina i povtorilos' vse kak vstar': pustaja banka vazelina, apteka, ulica, fonar'.",
		//		},

		//		pdu.Encode{
		//			Address: "+00000000000",
		//			Message: "Краткость, как известно, сестра таланта. Посему возьмем это на вооружение и освоим более краткий формат отправки SMS. А именно: из формулы SMS = SCA + TPDU исключается SCA.",
		//		},
	}

	incomming = []string{
		`+CMGR: 1,,25
07919740430900F302170B910000000000F0610130910435216101309104552100`,
		`+CMGR: 1,,159
07919740430900F3440B910000000000F00008610130915483218C050003210302044200200441043E043C043D0435043D044C044F0020043C043E0438002C000A04180020043D0435002004340430044E04420020043C043D04350020043704300441043D04430442044C000A041E043F044F0442044C00200442044004350432043E04330438002004320020043D043E04470438002C000A042F00200437043D0430044E0020`,
		`+CMGR: 1,,95
07919740430900F3440B910000000000F00008610130915493214C05000321030320130020044D0442043E0020043C043E04390020043F04430442044C002C000A04180020043F0440043E0434043E043B04360430044E00200438043404420438002E002E002E`,
		`+CMGR: 1,,159
07919740430900F3440B910000000000F00008610130915483218C050003210301042F0020043F0440043E0434043E043B04360430044E002004410432043E04390020043F04430442044C002C000A042F0020043F0440043E0434043E043B04360430044E00200438043404420438002C000A041800200445043E0442044C0020043F043E0440043E044E00200441043204350440043D04430442044C000A0425043E0442044F`,
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
	var err error
	var pduCoder pdu.Interface
	var messages []string
	var i int

	pduCoder = pdu.New().Decoder(messageReceiver)
	defer pduCoder.Done()

	// Encode SMS
	for i = range sms {
		var enc []string
		enc, err = pduCoder.Encoder(sms[i])
		//log.Infof("l: %d", len(sms[i].Message))
		if err != nil {
			//log.Errorf("Error encode message: %s", err.Error())
			continue
		}
		messages = append(messages, enc...)
	}
	for i = range messages {
		println(messages[i])
	}

	// Decode SMS
	for i = range incomming {
		pduCoder.Writer().Write([]byte(incomming[i]))
	}

	println()

}

// Receive new messages
func messageReceiver(msg pdu.Message) {
	var out string

	out += "New message found\n"
	if msg.Error() != nil {
		out += "Message error\n"
		out += msg.Error().Error() + "\n"
		return
	}

	out += " SMSC:"
	out += msg.ServiceCentreAddress()
	if msg.Type() == pdu.TypeSmsStatusReport {
		out += " (status report: ["
		out += msg.DischargeTime().String()
		out += "] '"
		out += msg.ReportStatus().String()
		out += "')"
	}
	out += "\n"

	out += " From:"
	out += msg.OriginatingAddress() + "\n"
	out += " SMS ("
	out += fmt.Sprintf("%d", msg.DataParts())
	out += "):"
	out += msg.Data() + "\n"

	println(out)
}
