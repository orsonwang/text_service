package main

import (
	//	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type exrate struct {
	inCashRate, outCashRate, inRate, outRate string
}

var exRates map[string]exrate

func crawer() {

	doc, err := goquery.NewDocument("http://rate.bot.com.tw/xrt?Lang=zh-TW")
	if err != nil {
		log.Println(err)
	}

	doc.Find("[class=\"hidden-phone print_show\"]").Each(func(i int, s *goquery.Selection) {
		currency := strings.TrimSpace(s.Text())
		pos := strings.Index(currency, " ")
		currCut := currency[0:pos]
		//		fmt.Printf("%s: ", currCut)
		inCashRate := s.Parent().Parent().Next().Text()
		//		inCashRate := s.NextUntil("[class=\"rate-content-cash text-right print_hide\"]").Text()
		outCashRate := s.Parent().Parent().Next().Next().Text()
		inRate := s.Parent().Parent().Next().Next().Next().Text()
		outRate := s.Parent().Parent().Next().Next().Next().Next().Text()
		var rate exrate
		rate.inCashRate = inCashRate
		rate.outCashRate = outCashRate
		rate.inRate = inRate
		rate.outRate = outRate
		exRates[currCut] = rate
		//		fmt.Printf("%s %s %s %s\n", inCashRate, outCashRate, inRate, outRate)
		//		fmt.Printf("%s %s %s %s\n", exRates["美金"].inCashRate, exRates["美金"].outCashRate, exRates["美金"].inRate, exRates["美金"].outRate)
	})

}
