package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type result struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
type bian struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           string `json:"openTime"`
	CloseTime          string `json:"closeTime"`
	FirstId            string `json:"firstId"`
	LastId             string `json:"lastId"`
	Count              string `json:"count"`
}

func main() {

	runtime.LockOSThread()
	cocoa.TerminateAfterWindowsClose = false
	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()

		nextClicked := make(chan bool)
		go func() {
			for {
				select {
				case <-time.After(1 * time.Second):
				case <-nextClicked:
				}
				//获取最新价格
				resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
				if err != nil {
					return
				}
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				var res result
				_ = json.Unmarshal(body, &res)
				price := res.Price
				price = price[:8]

				//获取24h涨跌
				resp, err = http.Get("https://api.binance.com/api/v3/ticker/24hr?symbol=BTCUSDT")
				if err != nil {
					return
				}
				defer resp.Body.Close()
				body, _ = ioutil.ReadAll(resp.Body)
				var rep bian
				_ = json.Unmarshal(body, &rep)
				change := rep.PriceChangePercent

				core.Dispatch(func() {
					obj.Button().SetTitle(fmt.Sprintf("BTC   $%v   %v％", price, change))

				})
			}
		}()
		nextClicked <- true

		itemQuit := cocoa.NSMenuItem_New()
		itemQuit.SetTitle("Quit")
		itemQuit.SetAction(objc.Sel("terminate:"))

		menu := cocoa.NSMenu_New()
		menu.AddItem(itemQuit)
		obj.SetMenu(menu)
	})
	app.Run()
}
