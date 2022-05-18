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

package tree

import (
	"fmt"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	d "github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/trace"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/logger"
)

const DefaultPageSize = 15
const KeyMap = " Keymap "
const Detail = " Detail "
const Quit = "<C-c>"

func DisplayList(ctx *cli.Context, displayable *d.Displayable) error {
	data := displayable.Data.(api.TraceBrief)
	condition := displayable.Condition.(*api.TraceQueryCondition)
	if err := ui.Init(); err != nil {
		logger.Log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	list := widgets.NewList()
	list.TitleStyle.Fg = ui.ColorRed
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false
	list.SelectedRowStyle = ui.Style{
		Fg:       ui.ColorBlack,
		Bg:       ui.ColorWhite,
		Modifier: ui.ModifierBold,
	}
	tree := widgets.NewTree()
	tree.TextStyle = ui.NewStyle(ui.ColorYellow)
	tree.WrapText = false
	tree.TitleStyle.Fg = ui.ColorRed
	tree.SelectedRowStyle = ui.Style{
		Fg:       ui.ColorBlack,
		Bg:       ui.ColorWhite,
		Modifier: ui.ModifierBold,
	}

	help := widgets.NewParagraph()
	help.WrapText = false
	help.Title = KeyMap
	help.Text = `[<Left>     ](fg:red,mod:bold) list activated
		[<Right>    ](fg:red,mod:bold) tree activated
		[K or <Up>  ](fg:red,mod:bold) list or tree Scroll Up
		[j or <Down>](fg:red,mod:bold) list or tree Scroll Down
		[<Ctr-b>    ](fg:red,mod:bold) list Page Up
		[<Ctr-f>    ](fg:red,mod:bold) list Page Down
		[p          ](fg:red,mod:bold) list Page Up
		[n          ](fg:red,mod:bold) list Page Down
		[<Home>     ](fg:red,mod:bold) list or tree Scroll to Top
		[<End>      ](fg:red,mod:bold) list or tree Scroll to Bottom
		[q or <Ctr-c>](fg:red,mod:bold) Quit
        `
	detail := widgets.NewParagraph()
	detail.Title = Detail
	detail.WrapText = false

	draw(list, tree, detail, help, data, ctx, condition)
	listenTracesKeyboard(list, tree, data, ctx, detail, help, condition)

	return nil
}

func draw(list *widgets.List, tree *widgets.Tree, detail, help *widgets.Paragraph, data api.TraceBrief,
	ctx *cli.Context, _ *api.TraceQueryCondition) {
	x, y := ui.TerminalDimensions()

	if len(data.Traces) != 0 {
		showIndex := list.SelectedRow
		var traceID = data.Traces[showIndex].TraceIds[0]
		list.Title = fmt.Sprintf("[%s]", traceID)
		nodes, serviceNames := getNodeData(ctx, traceID)
		tree.Title = fmt.Sprintf("[%s]", strings.Join(serviceNames, "->"))
		tree.SetNodes(nodes)
		list.Rows = rows(data, x/4)
		selected := extra[tree.SelectedNode()]
		detail.Text = selected.Detail
	} else {
		noData := "no data"
		list.Title = noData
		tree.Title = noData
		detail.Title = noData
	}

	list.SetRect(0, 0, x, y)
	tree.SetRect(x/5, 0, x, y)
	detail.SetRect(x-x/5, 0, x, y/2)
	help.SetRect(x-x/5, y/2, x, y)
	tree.ExpandAll()
	ui.Render(list, tree, detail, help)
}

func listenTracesKeyboard(list *widgets.List, tree *widgets.Tree, data api.TraceBrief, ctx *cli.Context,
	detail, help *widgets.Paragraph, condition *api.TraceQueryCondition) {
	uiEvents := ui.PollEvents()
	listActive := true
	var err error
	for {
		e := <-uiEvents

		switch e.ID {
		case "q", Quit:
			return
		case "<C-b>", "p":
			pageNum := *condition.Paging.PageNum
			if pageNum != 1 {
				pageNum--
				condition.Paging.PageNum = &pageNum
				data, err = trace.Traces(ctx, condition)
				if err != nil {
					logger.Log.Fatalln(err)
				}
			}
			tree.SelectedRow = 0
		case "<C-f>", "n":
			pageNum := *condition.Paging.PageNum
			pageNum++
			condition.Paging.PageNum = &pageNum
			data, err = trace.Traces(ctx, condition)
			if err != nil {
				logger.Log.Fatalln(err)
			}
			tree.SelectedRow = 0
		case "<Right>":
			listActive = false
		case "<Left>":
			listActive = true
		default:
			if action := listActions(e.ID, list, tree, listActive); action != nil {
				action()
			}
		}

		draw(list, tree, detail, help, data, ctx, condition)
	}
}

func listActions(key string, list *widgets.List, tree *widgets.Tree, listActive bool) func() {
	var f func()
	switch key {
	case "k", "<Up>":
		if listActive {
			f = list.ScrollUp
			tree.SelectedRow = 0
		} else {
			f = tree.ScrollUp
		}
	case "j", "<Down>":
		if listActive {
			tree.SelectedRow = 0
			f = list.ScrollDown
		} else {
			f = tree.ScrollDown
		}
	case "<Home>":
		if listActive {
			f = list.ScrollTop
		} else {
			f = tree.ScrollTop
		}
	case "<End>":
		if listActive {
			f = list.ScrollBottom
		} else {
			f = tree.ScrollBottom
		}
	}

	return f
}

func getNodeData(ctx *cli.Context, traceID string) (nodes []*widgets.TreeNode, serviceNames []string) {
	data, err := trace.Trace(ctx, traceID)

	if err != nil {
		logger.Log.Fatalln(err)
	}

	var roots []*Node
	roots, serviceNames = Adapt(data)

	nodes = make([]*widgets.TreeNode, len(roots))
	for i := range nodes {
		nodes[i] = &widgets.TreeNode{}
	}

	for i, root := range roots {
		adapt(root, nodes[i])
	}
	return nodes, serviceNames
}

func rows(data api.TraceBrief, subLen int) []string {
	var rows []string

	for _, t := range data.Traces {
		endpointName := t.EndpointNames[0]
		if len(endpointName) > subLen-3 {
			endpointName = endpointName[0:subLen-3] + "..."
		}

		rows = append(rows, fmt.Sprintf("[%s](mod:bold,fg:green) ", endpointName))
	}
	return rows
}
