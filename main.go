// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SubscriptionMessage struct {
	Did   string
	Type  string
	Data  string
	Time  time.Time
	Topic string
	Read  bool
}
type File struct {
	Cid  string
	Size int
}

func getFiles() []File {
	response, err := http.Get("http://localhost:4242/files")
	if err != nil {
		return nil
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	var f []File
	json.Unmarshal(responseData, &f)
	return f
}

func getFileCids(files []File) []string {
	us := make([]string, len(files))
	for i := range files {
		us[i] = files[i].Cid + "    " + fmt.Sprintf("%v", files[i].Size)
	}
	return us
}

func getFilesTotalSize(files []File) int {
	s := 0
	for i := range files {
		s += files[i].Size
	}
	return s
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = "Text"
	p.Text = "FLUBBER DEMO\nPress Q to quit."
	p.SetRect(0, 0, 50, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	updateParagraph := func(count int) {
		if count%2 == 0 {
			p.BorderStyle.Fg = ui.ColorCyan
			p.TextStyle.Fg = ui.ColorCyan
		} else if count%3 == 0 {
			p.BorderStyle.Fg = ui.ColorBlue
			p.TextStyle.Fg = ui.ColorBlue
		} else if count%5 == 0 {
			p.BorderStyle.Fg = ui.ColorRed
			p.TextStyle.Fg = ui.ColorRed
		} else if count%7 == 0 {
			p.BorderStyle.Fg = ui.ColorYellow
			p.TextStyle.Fg = ui.ColorYellow
		}
	}

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = getFileCids(getFiles())
	l.SetRect(0, 5, 75, 12)
	l.TextStyle.Fg = ui.ColorYellow

	g := widgets.NewGauge()
	g.Title = "Gauge"
	g.Percent = (getFilesTotalSize(getFiles()) / 1000000) * 100
	g.SetRect(0, 12, 75, 15)
	g.BarColor = ui.ColorRed
	g.BorderStyle.Fg = ui.ColorWhite
	g.TitleStyle.Fg = ui.ColorCyan

	sinData := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	lc := widgets.NewPlot()
	lc.Title = "dot-marker Line Chart"
	lc.Data = make([][]float64, 1)
	lc.Data[0] = sinData
	lc.SetRect(0, 15, 75, 25)
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorRed
	lc.Marker = widgets.MarkerDot

	barchartData := []float64{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}

	bc := widgets.NewBarChart()
	bc.Title = "Bar Chart"
	bc.SetRect(50, 0, 75, 10)
	bc.Labels = []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.BarColors[0] = ui.ColorGreen
	bc.NumStyles[0] = ui.NewStyle(ui.ColorBlack)

	lc2 := widgets.NewPlot()
	lc2.Title = "braille-mode Line Chart"
	lc2.Data = make([][]float64, 1)
	lc2.Data[0] = sinData
	lc2.SetRect(50, 15, 75, 25)
	lc2.AxesColor = ui.ColorWhite
	lc2.LineColors[0] = ui.ColorYellow

	p2 := widgets.NewParagraph()
	p2.Text = "Hey!\nI am a borderless block!"
	p2.Border = false
	p2.SetRect(50, 10, 75, 10)
	p2.TextStyle.Fg = ui.ColorMagenta

	draw := func(count int) {
		g.Percent = (getFilesTotalSize(getFiles()) / 1000000) * 100
		l.Rows = getFileCids(getFiles())
		lc.Data[0] = sinData[count/2%220:]
		lc2.Data[0] = sinData[2*count%220:]
		bc.Data = barchartData[count/2%10:]

		ui.Render(p, l, g, lc, bc, lc2, p2)
	}

	tickerCount := 1
	draw(tickerCount)
	tickerCount++
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			updateParagraph(tickerCount)
			draw(tickerCount)
			tickerCount++
		}
	}
}
