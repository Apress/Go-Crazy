package main

import (
    "encoding/csv"
    "fmt"
    "github.com/markcheno/go-talib"
    "log"
    "os"
    "strconv"
    "time"
)

type MarketDataPriceOnly struct {
    Date  time.Time
    Price float64
    High  float64
    Low   float64
}

const (
	factor = 3.0
	period = 10
)

type MarketData struct {
    Date          time.Time
    Price         float64
    High          float64
    Low           float64
    Rsi           float64
    Sma50         float64
    Sma200        float64
    EmaShort      float64
    EmaLong       float64
    Adx           float64
}

func ReadHistoricalDataFromCsvFile(csvFile string) ([]MarketDataPriceOnly, error) {
    file, err := os.Open(csvFile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    dateIndex, closeIndex, highIndex, lowIndex, openIndex := -1, -1, -1, -1, -1
    for i, column := range records[0] {
        if column == "datetime" {
            dateIndex = i
        }
        if column == "close" {
            closeIndex = i
        }
        if column == "open" {
            openIndex = i
        }
        if column == "high" {
            highIndex = i
        }
        if column == "low" {
            lowIndex = i
        }
        if dateIndex != -1 && closeIndex != -1 && highIndex != -1 && lowIndex != -1 && openIndex != -1 {
            break
        }
    }

    if dateIndex == -1 {
        log.Fatal("The 'datetime' column was not found in the CSV file.")
    }
    if closeIndex == -1 {
        log.Fatal("The 'close' column was not found in the CSV file.")
    }

    var dataPriceOnly []MarketDataPriceOnly
    
    for _, record := range records[1:] {
        //date, err := time.Parse(time.RFC3339, record[dateIndex])
        date, err := time.Parse("2006-01-02 15:04:05", record[dateIndex])
        if err != nil {
            log.Fatal(err)
        }

        price, err := strconv.ParseFloat(record[closeIndex], 64)
        if err != nil {
            log.Fatal(err)
        }

        high, err := strconv.ParseFloat(record[highIndex], 64)
        if err != nil {
            log.Fatal(err)
        }
        
        low, err := strconv.ParseFloat(record[lowIndex], 64)
        if err != nil {
            log.Fatal(err)
        }
        
        dataPriceOnly = append(dataPriceOnly, MarketDataPriceOnly{Date: date, Price: price, High: high, Low: low})
    }

    return dataPriceOnly, nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func BuildIndicatorsFromMarketData (slb int, llb int, histData []MarketDataPriceOnly) ([]MarketData, error) {

    close := make([]float64, len(histData))
    high  := make([]float64, len(histData))
    low   := make([]float64, len(histData))

    for i := 0; i < len(histData); i++ {
        close[i] = histData[i].Price
        high[i]  = histData[i].High
        low[i]   = histData[i].Low
    }

    rsi          := talib.Rsi(close, 14)
    emaShort     := talib.Ema(close, slb)
    emaLong      := talib.Ema(close, llb)
    sma50        := talib.Sma(close, 50)
    sma200       := talib.Sma(close, 200)
    adx          := talib.Adx(high, low, close, 14)

    var data []MarketData
    for i := 0; i < len(histData); i++ {

        data = append(
            data,
            MarketData{
                Date: histData[i].Date,
                Price: histData[i].Price,
                High: histData[i].High,
                Low: histData[i].Low,
                Rsi: rsi[i],
                Sma50: sma50[i],
                Sma200: sma200[i],
                EmaShort: emaShort[i],
                EmaLong: emaLong[i],
                Adx: adx[i],
            })
    }

    fmt.Printf("Finished preparing %d data points.\n\n", len(histData))
    
    return data, nil
}

func buildIndicators(slb int, llb int) ([]MarketData, []MarketData, error) {

    fmt.Println("File name to read: ", config.HourliesFile)
    historicalData, err := ReadHistoricalDataFromCsvFile(config.HourliesFile)
    if err != nil {
        return nil, nil, err
    }

    indicatorData, error := BuildIndicatorsFromMarketData(slb,llb,historicalData)
    if error != nil {
        return nil, nil, error
    }
    
    fmt.Println("File name to read: ", config.DailiesFile)    
    dailyHistoricalData, err := ReadHistoricalDataFromCsvFile(config.DailiesFile)
    if err != nil {
        return nil, nil, err
    }

    dailiesIndicatorData, error := BuildIndicatorsFromMarketData(slb,llb,dailyHistoricalData)
    if error != nil {
        return nil, nil, error
    }
    
    return indicatorData, dailiesIndicatorData, error
}

func saveIndicatorsToCsv (indicators []MarketData, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write the header
    header := []string {"Datetime", "Price", "Sma50", "Sma200", "EmaShort", "EmaLong", "Adx", "Rsi"}
    
    writer.Write(header)

    // Write the trade data
    for _, ind := range indicators {
        record := []string{
            ind.Date.Format(time.RFC3339),
            fmt.Sprintf("%.2f", ind.Price),
            fmt.Sprintf("%.2f", ind.Sma50),
            fmt.Sprintf("%.2f", ind.Sma200),
            fmt.Sprintf("%.2f", ind.EmaShort),
            fmt.Sprintf("%.2f", ind.EmaLong),
            fmt.Sprintf("%.2f", ind.Adx),
            fmt.Sprintf("%.2f", ind.Rsi),
        }
        writer.Write(record)
    }

    return nil
}



