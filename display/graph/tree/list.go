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

	d "github.com/apache/skywalking-cli/display/displayable"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/graphql/trace"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/logger"
)

const DefaultPageSize = 15
const keymap = " Keymap "
const cc = "<C-c>"

func DisplayList(ctx *cli.Context, displayable *d.Displayable) error {
	data := displayable.Data.(schema.TraceBrief)
	condition := displayable.Condition.(*schema.TraceQueryCondition)
	if err := ui.Init(); err != nil {
		logger.Log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	list := widgets.NewList()
	list.TitleStyle.Fg = ui.ColorRed
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false

	tree := widgets.NewTree()
	tree.TextStyle = ui.NewStyle(ui.ColorYellow)
	tree.WrapText = false
	tree.TitleStyle.Fg = ui.ColorRed

	help := widgets.NewParagraph()
	help.WrapText = false
	help.Title = keymap
	help.Text = `[k          ](fg:red,mod:bold) Scroll Up
		[<Up>       ](fg:red,mod:bold) Scroll Up
		[j          ](fg:red,mod:bold) Scroll Down
		[<Down>     ](fg:red,mod:bold) Scroll Down
		[<Ctr-b>    ](fg:red,mod:bold) list Page Up
		[<Ctr-f>    ](fg:red,mod:bold) list Page Down           
		[p          ](fg:red,mod:bold) list Page Up
		[n          ](fg:red,mod:bold) list Page Down
		[<Ctr-u>    ](fg:red,mod:bold) Scroll Half Page Up
		[<Ctr-d>    ](fg:red,mod:bold) Scroll Half Page Down
		[<Home>     ](fg:red,mod:bold) Scroll to Top
		[<Enter>    ](fg:red,mod:bold) Show Trace
		[<End>      ](fg:red,mod:bold) Scroll to Bottom
		[q          ](fg:red,mod:bold) Quit
		[<Ctr-c>    ](fg:red,mod:bold) Quit
        `
	draw(list, tree, help, data, 0, ctx, condition)
	listenTracesKeyboard(list, tree, data, ctx, help, condition)

	return nil
}

func draw(list *widgets.List, tree *widgets.Tree, help *widgets.Paragraph, data schema.TraceBrief, showIndex int,
	ctx *cli.Context, condition *schema.TraceQueryCondition) {
	x, y := ui.TerminalDimensions()

	if data.Total != 0 {
		var traceID = data.Traces[showIndex].TraceIds[0]
		list.Title = fmt.Sprintf("[ %d/%d  %s]", *condition.Paging.PageNum, totalPages(data.Total), traceID)
		nodes, serviceNames := getNodeData(ctx, traceID)
		tree.Title = fmt.Sprintf("[%s]", strings.Join(serviceNames, "->"))
		tree.SetNodes(nodes)
		list.Rows = rows(data, x/4)
	} else {
		noData := "no data"
		list.Title = noData
		tree.Title = noData
	}

	list.SetRect(0, 0, x, y)
	tree.SetRect(x/4, 0, x, y)
	help.SetRect(x-x/7, 0, x, y)
	tree.ExpandAll()
	ui.Render(list, tree, help)
}
func totalPages(total int) int {
	if total%DefaultPageSize == 0 {
		return total / DefaultPageSize
	}
	return total/DefaultPageSize + 1
}

func listenTracesKeyboard(list *widgets.List, tree *widgets.Tree, data schema.TraceBrief, ctx *cli.Context,
	help *widgets.Paragraph, condition *schema.TraceQueryCondition) {
	uiEvents := ui.PollEvents()
	for {
		showIndex := 0
		e := <-uiEvents

		switch e.ID {
		case "q", cc:
			return
		case "<C-b>", "p":
			pageNum := *condition.Paging.PageNum
			if pageNum != 1 {
				pageNum--
				condition.Paging.PageNum = &pageNum
				data = trace.Traces(ctx, condition)
			}
		case "<C-f>", "n":
			pageNum := *condition.Paging.PageNum
			if pageNum < totalPages(data.Total) {
				pageNum++
				condition.Paging.PageNum = &pageNum
				data = trace.Traces(ctx, condition)
			}
		default:
			if action := listActions(e.ID, list); action != nil {
				action()
			}
			showIndex = list.SelectedRow
		}
		draw(list, tree, help, data, showIndex, ctx, condition)
	}
}
func listActions(key string, list *widgets.List) func() {
	// mostly vim style
	actions := map[string]func(){
		"k":      list.ScrollUp,
		"<Up>":   list.ScrollUp,
		"j":      list.ScrollDown,
		"<Down>": list.ScrollDown,
		"<C-u>":  list.ScrollHalfPageUp,
		"<C-d>":  list.ScrollHalfPageDown,
		"<Home>": list.ScrollTop,
		"G":      list.ScrollBottom,
		"<End>":  list.ScrollBottom,
	}

	return actions[key]
}

func getNodeData(ctx *cli.Context, traceID string) (nodes []*widgets.TreeNode, serviceNames []string) {
	data := trace.Trace(ctx, traceID)
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

func rows(data schema.TraceBrief, subLen int) []string {
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
