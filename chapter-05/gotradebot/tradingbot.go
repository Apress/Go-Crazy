package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getClient() *binance.Client {

	fileName := fmt.Sprintf("%s/.config/binance/jun.key", os.Getenv("HOME"))
	file, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(file), "\n")

	// first line is binance API key
	// second line is binance secret key

	binanceAPIKey := lines[0]
	binanceSecretKey := lines[1]
	return binance.NewClient(binanceAPIKey, binanceSecretKey)
}

func StartBot(ctx context.Context, symbol string, interval string, capital float64) {

	client := getClient()

	select {
	case <-ctx.Done():
		fmt.Println("has just been canceled")
	default:
		time.Sleep(100 * time.Millisecond)
		runStrategy(client, symbol, interval, capital)
	}

}

func runStrategy(client *binance.Client, symbol string, interval string, capital float64) bool {

	// setup strategy to run
	GCStrategy := &GoldenCrossMaStrategy{}
	GCStrategy.SetShortLookback(50)
	GCStrategy.SetShortLookback(200)

	// Fetch klines data for the specified symbol and interval
	klines, err := client.NewKlinesService().Symbol(symbol).Interval(interval).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Extract the OHLC prices from the klines data
	var opens, highs, lows, closes []float64
	for _, kline := range klines {

		close, _ := strconv.ParseFloat(kline.Close, 64)
		open, _ := strconv.ParseFloat(kline.Open, 64)
		high, _ := strconv.ParseFloat(kline.High, 64)
		low, _ := strconv.ParseFloat(kline.Low, 64)

		closes = append(closes, close)
		opens = append(opens, open)
		highs = append(highs, high)
		lows = append(lows, low)

	}

	var indexes []int
	for i := 0; i < len(closes); i++ {
		indexes = append(indexes, i)
	}

	// create all indicator data from price action data
	indicatorData, err := buildIndicators(5, 9, indexes, opens, highs, lows, closes)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return true
	}

	positionOpen := false

	checkIndex := len(closes) - 1

	if !positionOpen && GCStrategy.ShouldEnterMarket(indicatorData, checkIndex) {

		_, err := checkFunds(client, capital)
		if err {
			return true
		}
		buy(symbol, capital, client)
		positionOpen = true
		fmt.Println("Entry price: ", closes[checkIndex])
	} else if positionOpen && GCStrategy.ShouldExitMarket(indicatorData, checkIndex) {

		// selling position
		positionOpen = false
		sell(symbol, capital, client)
		fmt.Println("Exit price: ", closes[checkIndex])

	}

	return false

}

func buy(symbol string, capital float64, client *binance.Client) {
	// Place a market buy order for the specified symbol and capital
	log.Printf("Buying %s with %f USDT\n", symbol, capital)

	order, err := client.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(strconv.FormatFloat(capital, 'f', 2, 64)).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Market buy order %s executed at price %s\n", order.OrderID, order.Price)

}

func checkFunds(client *binance.Client, capital float64) (error, bool) {
	// Check if there are available funds to buy
	balance, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var availableBalance float64
	for _, b := range balance.Balances {
		if b.Asset == "USDT" {
			availableBalance, _ = strconv.ParseFloat(b.Free, 64)
			break
		}
	}
	log.Printf("Current balance (%f) to buy with capital %f", availableBalance, capital)

	if availableBalance < capital {
		log.Printf("Not enough available balance (%f) to buy with capital %f", availableBalance, capital)
		return nil, true
	}
	return err, false
}

func sell(symbol string, capital float64, client *binance.Client) {
	// Place a market sell order for the specified symbol and quantity
	log.Printf("Selling %s with %f USDT\n", symbol, capital)
	order, err := client.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).QuoteOrderQty(strconv.FormatFloat(capital, 'f', 2, 64)).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Market sell order %s executed at price %s\n", order.OrderID, order.Price)

}
