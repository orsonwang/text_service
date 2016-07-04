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

	doc, err := goquery.NewDocument("http://rate.bot.com.tw/Pages/Static/UIP003.zh-TW.htm")
	if err != nil {
		log.Println(err)
	}

	doc.Find("[class=\"titleLeft\"]").Each(func(i int, s *goquery.Selection) {
		currency := strings.TrimSpace(s.Text())
		pos := strings.Index(currency, " ")
		currCut := currency[0:pos]
		//		fmt.Printf("%s: ", currCut)
		inCashRate := s.Next().Text()
		outCashRate := s.Next().Next().Text()
		inRate := s.Next().Next().Next().Text()
		outRate := s.Next().Next().Next().Next().Text()
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
