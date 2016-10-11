package main

import (
	logger "log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/nats-io/nats"
)

var log *logger.Logger

func matchString(pattern, text string) (result bool) {
	result, _ = regexp.MatchString(pattern, text)
	return
}

// OnTextMessage ...
func onTextMessage(text string) (strResult string) {
	log.Printf("Receive Text: %s\n", text)
	strAfterCut := strings.ToUpper(text)
	strResult = "目前系統功能如下\n" +
		"利率(預設為台幣),外幣利率\n" +
		"匯率(預設為總表),美元,日圓與人民幣匯率\n" +
		"台外幣各類存款餘額\n" +
		"及信用卡功能"

	switch {
	// 利匯率服務
	case matchString("外幣+.*利率+.*", strAfterCut):
		strResult = "常用外幣利率表\n 美元 定存 2.3% 活存 1.8% \n 日圓 定存 0.1% 活存 0.1%"
		break
	case matchString("台幣+.*利率+.*", strAfterCut):
	case matchString("利率.*", strAfterCut):
		strResult = "台幣活存利率表 \n 活存 0.5% 活儲 0.6% \n 定存\n 三個月 0.76% 六個月 0.78% 一年 0.80% 三年 0.80%\n https://www.skbank.com.tw/RAT/RAT2_TWSaving.aspx"
		break
	case matchString("(美元|美金|USD)+.*歷史+.*匯率+.*", strAfterCut):
		strResult = "美元歷史匯率參考 http://tw.exchange-rateorg/history/TWD/USD/T"
		break
	case matchString("(日圓|日元|日幣|JPY)+.*歷史+.*匯率+.*", strAfterCut):
		strResult = "日元歷史匯率參考 http://tw.exchange-rateorg/history/TWD/JPY/T"
		break
	case matchString("(人民幣|RMB)+.*歷史+.*匯率+.*", strAfterCut):
		strResult = "人民幣歷史匯率參考 http://tw.exchange-rateorg/history/TWD/CNY/T"
		break
	case matchString("歷史+.*匯率+.*", strAfterCut):
		strResult = "歷史匯率參考 http://tw.exchange-rateorg/history/TWD/USD/T"
		break
	case matchString("(美元|美金|USD)+.*匯率+.*", strAfterCut):
		strResult = "美金匯率\n" +
			"現金買入 " + exRates["美金"].inCashRate + "\n" +
			"現金賣出 " + exRates["美金"].outCashRate + "\n" +
			"即期買入 " + exRates["美金"].inRate + "\n" +
			"即期賣出 " + exRates["美金"].outRate
		break
	case matchString("(日圓|日元|日幣|JPY)+.*匯率+.*", strAfterCut):
		strResult = "日圓匯率\n" +
			"現金買入 " + exRates["日圓"].inCashRate + "\n" +
			"現金賣出 " + exRates["日圓"].outCashRate + "\n" +
			"即期買入 " + exRates["日圓"].inRate + "\n" +
			"即期賣出 " + exRates["日圓"].outRate
		break
	case matchString("(人民幣|RMB)+.*匯率+.*", strAfterCut):
		strResult = "人民幣匯率\n" +
			"現金買入 " + exRates["人民幣"].inCashRate + "\n" +
			"現金賣出 " + exRates["人民幣"].outCashRate + "\n" +
			"即期買入 " + exRates["人民幣"].inRate + "\n" +
			"即期賣出 " + exRates["人民幣"].outRate
		break
	case matchString("(港幣|HKD)+.*匯率+.*", strAfterCut):
		strResult = "港幣匯率\n" +
			"現金買入 " + exRates["港幣"].inCashRate + "\n" +
			"現金賣出 " + exRates["港幣"].outCashRate + "\n" +
			"即期買入 " + exRates["港幣"].inRate + "\n" +
			"即期賣出 " + exRates["港幣"].outRate
		break
	case matchString("(英鎊|GBP)+.*匯率+.*", strAfterCut):
		strResult = "英鎊匯率\n" +
			"現金買入 " + exRates["英鎊"].inCashRate + "\n" +
			"現金賣出 " + exRates["英鎊"].outCashRate + "\n" +
			"即期買入 " + exRates["英鎊"].inRate + "\n" +
			"即期賣出 " + exRates["英鎊"].outRate
		break
	case matchString("(歐元|EUR)+.*匯率+.*", strAfterCut):
		strResult = "歐元匯率\n" +
			"現金買入 " + exRates["歐元"].inCashRate + "\n" +
			"現金賣出 " + exRates["歐元"].outCashRate + "\n" +
			"即期買入 " + exRates["歐元"].inRate + "\n" +
			"即期賣出 " + exRates["歐元"].outRate
		break
	case matchString("匯率+.*(美元|美金|USD)+.*", strAfterCut):
		strResult = "美金匯率\n" +
			"現金買入 " + exRates["美金"].inCashRate + "\n" +
			"現金賣出 " + exRates["美金"].outCashRate + "\n" +
			"即期買入 " + exRates["美金"].inRate + "\n" +
			"即期賣出 " + exRates["美金"].outRate
		break
	case matchString("匯率+.*(日圓|日元|日幣|JPY)+.*", strAfterCut):
		strResult = "日圓匯率\n" +
			"現金買入 " + exRates["日圓"].inCashRate + "\n" +
			"現金賣出 " + exRates["日圓"].outCashRate + "\n" +
			"即期買入 " + exRates["日圓"].inRate + "\n" +
			"即期賣出 " + exRates["日圓"].outRate
		break
	case matchString("匯率+.*(人民幣|RMB)+.*", strAfterCut):
		strResult = "人民幣匯率\n" +
			"現金買入 " + exRates["人民幣"].inCashRate + "\n" +
			"現金賣出 " + exRates["人民幣"].outCashRate + "\n" +
			"即期買入 " + exRates["人民幣"].inRate + "\n" +
			"即期賣出 " + exRates["人民幣"].outRate
		break
	case matchString("匯率+.*(港幣|HKD)+.*", strAfterCut):
		strResult = "港幣匯率\n" +
			"現金買入 " + exRates["港幣"].inCashRate + "\n" +
			"現金賣出 " + exRates["港幣"].outCashRate + "\n" +
			"即期買入 " + exRates["港幣"].inRate + "\n" +
			"即期賣出 " + exRates["港幣"].outRate
		break
	case matchString("匯率+.*(英鎊|GBP)+.*", strAfterCut):
		strResult = "英鎊匯率\n" +
			"現金買入 " + exRates["英鎊"].inCashRate + "\n" +
			"現金賣出 " + exRates["英鎊"].outCashRate + "\n" +
			"即期買入 " + exRates["英鎊"].inRate + "\n" +
			"即期賣出 " + exRates["英鎊"].outRate
		break
	case matchString("匯率+.*(歐元|EUR)+.*", strAfterCut):
		strResult = "歐元匯率\n" +
			"現金買入 " + exRates["歐元"].inCashRate + "\n" +
			"現金賣出 " + exRates["歐元"].outCashRate + "\n" +
			"即期買入 " + exRates["歐元"].inRate + "\n" +
			"即期賣出 " + exRates["歐元"].outRate
		break
	case matchString("匯率+.*", strAfterCut):
		strResult = ""
		break
		//存款服務
	case matchString("(美元|美金|USD)+.*(活存|存款)+.*(餘額)?.*", strAfterCut):
		strResult = "您的美元活存帳戶餘額為: 233,188.66 美元"
		break
	case matchString("(日圓|日元|日幣|JPY)+.*(活存|存款)+.*(餘額)?.*", strAfterCut):
		strResult = "您的日元活存帳戶餘額為: 233,188.66 日元"
		break
	case matchString("(人民幣|RMB)+.*(活存|存款)+.*餘額?.*", strAfterCut):
		strResult = "您沒有人民幣帳戶，若要開立請點連結 https://virtual.bank"
		break
	case matchString("(美元|美金|USD)+.*(定存|存單)+.*(餘額)?.*", strAfterCut):
		strResult = "您的美元定存帳戶餘額為: 1,000.00 美元"
		break
	case matchString("(日圓|日元|日幣|JPY)+.*(定存|存單)+.*(餘額)?.*", strAfterCut):
		strResult = "您沒有日元定存帳戶，若要開立請點連結 https://virtual.bank"
		break
	case matchString("(人民幣|RMB)+.*(定存|存單)+.*(餘額)?.*", strAfterCut):
		strResult = "您沒有人民幣帳戶，若要開立請點連結 https://virtual.bank"
		break
	case matchString("(存款|活存|帳戶)+.*(餘額)?.*", strAfterCut):
		strResult = "您的台幣活存帳戶餘額為: 233,188.66 元\n "
		break
	case matchString("(定存|存單)+.*(餘額)?.*", strAfterCut):
		strResult = "您的台幣定存帳戶餘額為: 1,000,000.00 元\n"
		break
		// 只有幣別就直接給匯率
	case matchString("(美元|美金|USD)+.*", strAfterCut):
		strResult = "美金匯率\n" +
			"現金買入 " + exRates["美金"].inCashRate + "\n" +
			"現金賣出 " + exRates["美金"].outCashRate + "\n" +
			"即期買入 " + exRates["美金"].inRate + "\n" +
			"即期賣出 " + exRates["美金"].outRate
		break
	case matchString("(日圓|日元|日幣|JPY)+.*", strAfterCut):
		strResult = "日圓匯率\n" +
			"現金買入 " + exRates["日圓"].inCashRate + "\n" +
			"現金賣出 " + exRates["日圓"].outCashRate + "\n" +
			"即期買入 " + exRates["日圓"].inRate + "\n" +
			"即期賣出 " + exRates["日圓"].outRate
		break
	case matchString("(人民幣|RMB)+.*", strAfterCut):
		strResult = "人民幣匯率\n" +
			"現金買入 " + exRates["人民幣"].inCashRate + "\n" +
			"現金賣出 " + exRates["人民幣"].outCashRate + "\n" +
			"即期買入 " + exRates["人民幣"].inRate + "\n" +
			"即期賣出 " + exRates["人民幣"].outRate
		break
	case matchString("(港幣|HKD)+.*", strAfterCut):
		strResult = "港幣匯率\n" +
			"現金買入 " + exRates["港幣"].inCashRate + "\n" +
			"現金賣出 " + exRates["港幣"].outCashRate + "\n" +
			"即期買入 " + exRates["港幣"].inRate + "\n" +
			"即期賣出 " + exRates["港幣"].outRate
		break
	case matchString("(英鎊|GBP)+.*", strAfterCut):
		strResult = "英鎊匯率\n" +
			"現金買入 " + exRates["英鎊"].inCashRate + "\n" +
			"現金賣出 " + exRates["英鎊"].outCashRate + "\n" +
			"即期買入 " + exRates["英鎊"].inRate + "\n" +
			"即期賣出 " + exRates["英鎊"].outRate
		break
	case matchString("(歐元|EUR)+.*", strAfterCut):
		strResult = "歐元匯率\n" +
			"現金買入 " + exRates["歐元"].inCashRate + "\n" +
			"現金賣出 " + exRates["歐元"].outCashRate + "\n" +
			"即期買入 " + exRates["歐元"].inRate + "\n" +
			"即期賣出 " + exRates["歐元"].outRate
		break
		// 信用卡服務
	case matchString("信用卡+.*最低應繳+.*", strAfterCut):
		strResult = "您這個月的信用卡帳單\n" +
			"最低應繳金額： 1,234 元\n" +
			" 繳款截止日： 5月28號"
		break
	case matchString("信用卡+.*應繳+.*", strAfterCut):
		strResult = "您這個月的信用卡帳單\n" +
			"應繳金額： 12,345 元\n" +
			" 繳款截止日： 5月28號"
		break
	case matchString("(補寄|郵寄)?.*(上個月|上月|上一期|上期|前期)+.*(帳單)+.*", strAfterCut):
		strResult = "您申請的上個月信用卡帳單\n" +
			"已經排入郵寄系統，麻煩您注意最近的郵局信件"
		break
	case matchString("補寄+.*帳單+.*", strAfterCut):
		strResult = "如果您要申請補寄帳單可以使用以下指令作業\n" +
			"\"郵寄x月帳單\"系統將會自動郵寄該月帳單\n" +
			"\"傳真x月帳單到0xxxxxxxxx\"系統會將該月帳單傳真至指定號碼\n" +
			"\"電郵x月帳單\"系統會將該月電子帳單發送至設定的電子郵件信箱"
		break
	case matchString("(這個月|當月|這一期|當期)+.*信用卡+.*帳單+.*", strAfterCut):
	case matchString("信用卡+.*帳單+.*", strAfterCut):
		strResult = "您這個月的信用卡帳單\n" +
			"應繳金額： 12,345 元\n" +
			"最低應繳金額： 1,234 元\n" +
			" 繳款截止日： 5月28號"
		break
	case matchString("信用卡+.*餘額+.*", strAfterCut):
		strResult = "您目前信用卡\n" +
			"可用餘額： 13,655 元\n" +
			"可臨時調整額度為： 8,000 元"
		break
	case matchString("信用卡+.*額度+.*", strAfterCut):
		strResult = "您目前信用卡\n" +
			"可用總額度： 50,000 元\n" +
			"可用餘額： 13,655 元"
		break
	case matchString("信用卡+.*(最近|最新)+.*交易+.*", strAfterCut):
		strResult = "您信用卡最近五筆\n" +
			"Paypal: 5.00(USD) , 約 165.03 元\n" +
			"新光三越站前店美食: 180.00 元\n" +
			"新光三越站前店美食: 380.00 元\n" +
			"新光人壽保險費: 24,000.00 元\n" +
			"所得稅: 18,180.00 元"
		break
	}
	log.Printf("Return text: %s\n", strResult)
	return
}

func main() {
	f, err := os.OpenFile("./textservice.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Can't open log file: %v\n", err)
	}
	log = new(logger.Logger)
	log.SetOutput(f)

	urls := "nats://localhost:4222"
	showTime := true

	ticker := time.NewTicker(1 * 60 * time.Second)

	exRates = make(map[string]exrate, 0)
	crawer()
	go func() {
		for range ticker.C {
			crawer()

		}
	}()

	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect to NATS: %v\n", err)
	}
	var subj = "aitc.text.service"

	nc.Subscribe(subj, func(msg *nats.Msg) {
		reply := onTextMessage(string(msg.Data))
		nc.Publish(msg.Reply, []byte(reply))
	})

	log.Printf("Listening on [%s]\n", subj)
	if showTime {
		log.SetFlags(logger.LstdFlags)
	}

	runtime.Goexit()
}
