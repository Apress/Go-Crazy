package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func PnLCalcs(my_trades []Trade) float64 {

	var totalPnL float64
	for _, r := range my_trades {
		totalPnL += r.PnL
	}

	return totalPnL
}

func saveTradesToCSV(trades []Trade, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Datetime", "Indicator", "Price", "Quantity", "PositionLength", "PnL"}
	writer.Write(header)

	// Write the trade data
	for _, trade := range trades {
		record := []string{
			trade.Datetime.Format(time.RFC3339),
			trade.Indicator,
			fmt.Sprintf("%.2f", trade.Price),
			fmt.Sprintf("%.2f", trade.Quantity),
            fmt.Sprintf("%d", trade.Poslen),
            fmt.Sprintf("%.2f", trade.PnL),
		}
		writer.Write(record)
	}

	return nil
}
