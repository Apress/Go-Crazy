package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var (
	rsiBuyThreshold  float64
	rsiSellThreshold float64
	fastPeriod       int
	slowPeriod       int
	capital          float64
)

func main() {
	// Create a new fyne application
	app := app.New()

	// Create a new window
	win := app.NewWindow("Binance Trading Bot")

	var form *widget.Form
	var started = false

	var ctx, cancel = context.WithCancel(context.Background())
	// Create a new form with input fields

	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Symbol", Widget: widget.NewEntry(), HintText: "BTCUSDT"},
			{Text: "Interval", Widget: widget.NewSelect([]string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}, func(s string) {

			})},
			{Text: "Capital", Widget: widget.NewEntry()},
		},
		OnSubmit: func() {

			if started == true {
				fmt.Printf("Stopping...\n")
				cancel()
				started = !started
				form.SubmitText = "Start"
				form.Refresh()
			} else {
				fmt.Printf("Starting...\n")
				started = !started
				ctx, cancel = context.WithCancel(context.Background())
				// Get the input values from the form
				symbol := form.Items[0].Widget.(*widget.Entry).Text
				interval := form.Items[1].Widget.(*widget.Select).Selected
				capital, _ = strconv.ParseFloat(form.Items[2].Widget.(*widget.Entry).Text, 64)

				// Start the bot
				go StartBot(ctx, symbol, interval, capital)

				form.SubmitText = "Stop"
				form.Refresh()
			}

		},
	}

	form.Items[0].Widget.(*widget.Entry).Text = "BTCUSDT"
	form.Items[1].Widget.(*widget.Select).Selected = "15m"
	form.Items[2].Widget.(*widget.Entry).Text = "50"

	form.SubmitText = "Start"

	// Add the form and the start button to a container
	//content := container.NewVBox(form, startButton)
	content := container.NewVBox(form)

	// Set the window content
	win.SetContent(content)

	// Show the window
	win.ShowAndRun()
}
