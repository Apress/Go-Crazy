package main

import (
	"fmt"
    "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

// Config holds the application configuration
type Config struct {
    HourliesFile      string  `yaml:"hourliesfile"`
    DailiesFile       string  `yaml:"dailiesfile"`
    HourliesLeverage  float64 `yaml:"hourliesleverage"`
    HourliesTradeSize float64 `yaml:"hourliestradesize"`
    DailiesLeverage   float64 `yaml:"dailiesleverage"`
    DailiesTradeSize  float64 `yaml:"dailiestradesize"`
    CashBalance       float64 `yaml:"cashbalance"`
}

// struct to save trades
type Trade struct {
	Datetime     time.Time
	Indicator    string
	Price        float64
	Quantity     float64
    Poslen       int
    PnL          float64
}

// Global variable to hold the configuration
var config Config

// Read and parse the configuration file
func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
}


func generate_report (stratName string, profits []float64, trades []Trade) (err error) {

	//fmt.Printf("Profit from backtesting: %.2f\n", profits)
    fmt.Println("\n\nStrategy: ", stratName)
    
	// PnL
	totalPnL := PnLCalcs(trades)
	fmt.Printf("PnL: %.2f\n", totalPnL)

	// Calculate the Sharpe Ratio (assuming a risk-free rate of 0.02)
	sharpeRatio := SharpeRatio(profits, 0.02)
	fmt.Printf("Sharpe Ratio: %.2f\n", sharpeRatio)

	// Calculate the MAR Ratio
	marRatio := MARRatio(profits)
	fmt.Printf("MAR Ratio: %.2f\n", marRatio)

    return nil

}

func main() {

    // create all indicator data from price action data
    // this is done loading a CSV file with Open High Low Close
    // two variable length moving averages are set as arguments
    indicatorData, dailiesIndicatorData, err := buildIndicators(5, 9)
    if err != nil {
        fmt.Println("Error fetching data:", err)
        return
    }

    leverage := config.HourliesLeverage
    trade_size := config.HourliesTradeSize
    ema_cash_balance := config.CashBalance

    fmt.Print("Starting strategy calculations with:\n")
    fmt.Printf("Starting Cash Balance: %.2f\n", config.CashBalance)
    fmt.Printf("    Hourlies Leverage: %.2f\n", config.HourliesLeverage)
    fmt.Printf("     Dailies Leverage: %.2f\n", config.DailiesLeverage)
    fmt.Printf("  Hourlies Trade Size: %.2f\n", config.HourliesTradeSize)
    fmt.Printf("   Dailies Trade Size: %.2f\n\n", config.DailiesTradeSize)

    // below two strategies are run sequentially
	EmaStrategy := &MovingAverageCrossoverStrategy{}
	EmaStrategy.SetShortLookback(5)
	EmaStrategy.SetLongLookback(9)
    GCStrategy := &GoldenCrossMaStrategy{}
    GCStrategy.SetShortLookback(50)
    GCStrategy.SetShortLookback(200)

    emaprofit := 0.0
    gcprofit := 0.0

	// Assuming returns is a slice of float64 representing the strategy's returns
	var emaprofits []float64
	// Assuming returns is a slice of float64 representing the strategy's trades
	var ematrades []Trade
	// Fill the returns slice with your strategy's returns
	emaprofits, ematrades, emaprofit = Backtest(EmaStrategy, indicatorData, leverage, trade_size, ema_cash_balance)

    fmt.Printf("Ema Total PnL: %.2f\n", emaprofit)

    // output statistics to STDOUT
    err = generate_report("Ema100Strategy", emaprofits, ematrades)
    if err != nil {
        panic (err)
    }

    // storing all executed trades into a CSV file
    // this is particularly useful when comparing against a chart
    filename := "ema_trades.csv"
    err = saveTradesToCSV(ematrades, filename)
    if err != nil {
        fmt.Println("Error saving trades to CSV:", err)
    } else {
        fmt.Printf("Trades saved to %s\n\n\n", filename)
    }

    // Assuming returns is a slice of float64 representing the strategy's returns
	var gcprofits []float64
    var gctrades []Trade
	// Fill the returns slice with your strategy's returns
	gcprofits, gctrades, gcprofit = BacktestGoldenCross(GCStrategy, indicatorData, leverage, trade_size, ema_cash_balance)

    fmt.Printf("Golden Cross Total PnL: %.2f\n", gcprofit)
    
    // output statistics to STDOUT    
    err = generate_report("GoldenCross", gcprofits, gctrades)
    if err != nil {
        panic(err)
    }
    
    // storing all executed trades into a CSV file
    // this is particularly useful when comparing against a chart
    filename = "golden_cross_trades.csv"
    err = saveTradesToCSV(gctrades, filename)
    if err != nil {
        fmt.Println("Error saving trades to CSV:", err)
    } else {
        fmt.Printf("Trades saved to %s\n", filename)
    }
    
    // export forensic analysis
    filename = "indicators.csv"
    err = saveIndicatorsToCsv(indicatorData, filename)
    if err != nil {
        fmt.Printf("Error saving indicators to CSV:", err)
    } else {
        fmt.Printf("Indicators saved to %s\n\n", filename)
    }

    // export forensic analysis
    filename = "dailies_indicators.csv"
    err = saveIndicatorsToCsv(dailiesIndicatorData, filename)
    if err != nil {
        fmt.Printf("Error saving daily indicators to CSV:", err)
    } else {
        fmt.Printf("dailyIndicators saved to %s\n\n", filename)
    }
    
}
