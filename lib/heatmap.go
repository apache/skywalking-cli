// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lib

import (
	"fmt"
	im "image"

	ui "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"
)

type HeatMap struct {
	ui.Block
	XLabelStyles []ui.Style
	CellColors   [][]ui.Color
	NumStyles    [][]ui.Style
	Formatter    func(float64) string
	Data         [][]float64
	XLabels      []string
	YLabels      []string
	CellWidth    int
	CellGap      int
}

func NewHeatMap() *HeatMap {
	return &HeatMap{
		Block:        *ui.NewBlock(),
		CellColors:   [][]ui.Color{ui.StandardColors, ui.StandardColors},
		NumStyles:    [][]ui.Style{ui.StandardStyles, ui.StandardStyles},
		Formatter:    func(n float64) string { return fmt.Sprint(n) },
		XLabelStyles: ui.StandardStyles,
		CellGap:      1,
		CellWidth:    3,
	}
}

func (hm *HeatMap) Draw(buffer *ui.Buffer) {
	hm.Block.Draw(buffer)

	cellX := hm.Inner.Min.X

	for i, column := range hm.Data {
		cellY := 0
		for j, datum := range column {
			buffer.SetString(
				hm.YLabels[j],
				ui.StyleClear,
				im.Pt(hm.Inner.Min.X, (hm.Inner.Max.Y-2)-cellY),
			)
			for x := cellX + 5; x < ui.MinInt(cellX+hm.CellWidth, hm.Inner.Max.X)+5; x++ {
				for y := (hm.Inner.Max.Y - 2) - cellY; y > (hm.Inner.Max.Y-2)-cellY-1; y-- {
					cell := ui.NewCell(' ', ui.NewStyle(ui.ColorClear, color(hm.CellColors, i, j)))
					buffer.SetCell(cell, im.Pt(x, y))
				}
			}

			if hm.Formatter != nil {
				hm.drawNumber(buffer, datum, i, j, cellX+5, cellY)
			}

			cellY++
		}

		if i < len(hm.XLabels) {
			hm.drawLabel(buffer, cellX+5, i)
		}

		cellX += hm.CellWidth + hm.CellGap
	}
}

func (hm *HeatMap) drawLabel(buffer *ui.Buffer, cellX, i int) {
	labelX := cellX + ui.MaxInt(int(float64(hm.CellWidth)/2)-int(float64(rw.StringWidth(hm.XLabels[i]))/2), 0)
	buffer.SetString(
		ui.TrimString(hm.XLabels[i], hm.CellWidth),
		ui.SelectStyle(hm.XLabelStyles, i),
		im.Pt(labelX, hm.Inner.Max.Y-1),
	)
}

func (hm *HeatMap) drawNumber(buffer *ui.Buffer, datum float64, i, j, cellX, cellY int) {
	x := cellX + int(float64(hm.CellWidth)/2) - 1
	numberStyle := style(hm.NumStyles, i, j)
	cellColor := color(hm.CellColors, i, j)
	buffer.SetString(
		hm.Formatter(datum),
		ui.NewStyle(numberStyle.Fg, cellColor, numberStyle.Modifier),
		im.Pt(x, (hm.Inner.Max.Y-2)-cellY),
	)
}

func color(colors [][]ui.Color, i, j int) ui.Color {
	return colors[i%len(colors)][j%len(colors)]
}

func style(styles [][]ui.Style, i, j int) ui.Style {
	return styles[i%len(styles)][j%len(styles)]
}
