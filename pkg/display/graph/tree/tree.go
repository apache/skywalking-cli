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

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/apache/skywalking-cli/internal/logger"
)

type Node struct {
	Children []*Node
	Detail   string
	Value    fmt.Stringer
}

var extra = make(map[*widgets.TreeNode]*Node)

func Display(roots []*Node, serviceNames []string) error {
	if err := ui.Init(); err != nil {
		logger.Log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	nodes := make([]*widgets.TreeNode, len(roots))
	for i := range nodes {
		nodes[i] = &widgets.TreeNode{}
	}

	for i, root := range roots {
		adapt(root, nodes[i])
	}

	tree := widgets.NewTree()
	tree.TextStyle = ui.Style{
		Fg:       ui.ColorWhite,
		Bg:       ui.ColorClear,
		Modifier: 0,
	}
	tree.SelectedRowStyle = ui.Style{
		Fg:       ui.ColorBlack,
		Bg:       ui.ColorWhite,
		Modifier: ui.ModifierBold,
	}
	tree.WrapText = false
	tree.SetNodes(nodes)
	tree.Title = fmt.Sprintf("[ %s ]        [%s]", strings.Join(serviceNames, "->"), " Press ? to show help ")
	tree.TitleStyle.Modifier = ui.ModifierBold
	tree.TitleStyle.Fg = ui.ColorRed

	x, y := ui.TerminalDimensions()

	tree.SetRect(0, 0, x, y)

	detail := widgets.NewParagraph()
	detail.Title = Detail
	detail.WrapText = false
	detail.SetRect(x, 0, x, y)

	help := widgets.NewParagraph()
	help.WrapText = false
	help.SetRect(x, 0, x, y)
	help.Title = KeyMap
	help.Text = `
		[?          ](fg:red,mod:bold) Toggle this help
		[k          ](fg:red,mod:bold) Scroll Up
		[<Up>       ](fg:red,mod:bold) Scroll Up
		[j          ](fg:red,mod:bold) Scroll Down
		[<Down>     ](fg:red,mod:bold) Scroll Down
		[<Ctr-b>    ](fg:red,mod:bold) Scroll Page Up
		[<Ctr-u>    ](fg:red,mod:bold) Scroll Half Page Up
		[<Ctr-f>    ](fg:red,mod:bold) Scroll Page Down
		[<Ctr-d>    ](fg:red,mod:bold) Scroll Half Page Down
		[<Home>     ](fg:red,mod:bold) Scroll to Top
		[gg         ](fg:red,mod:bold) Scroll to Top
		[<Enter>    ](fg:red,mod:bold) Show Trace Detail
		[<Space>    ](fg:red,mod:bold) Show Trace Detail
		[o          ](fg:red,mod:bold) Toggle Expand
		[G          ](fg:red,mod:bold) Scroll to Bottom
		[<End>      ](fg:red,mod:bold) Scroll to Bottom
		[E          ](fg:red,mod:bold) Expand All
		[C          ](fg:red,mod:bold) Collapse All
		[q          ](fg:red,mod:bold) Quit
		[<Ctr-c>    ](fg:red,mod:bold) Quit
	`

	ui.Render(tree, detail, help)

	listenKeyboard(tree, detail, help)

	return nil
}

func adapt(n1 *Node, n2 *widgets.TreeNode) {
	if n1 == nil || n2 == nil {
		return
	}

	n2.Expanded = true
	n2.Value = n1.Value
	n2.Nodes = []*widgets.TreeNode{}
	extra[n2] = n1

	for _, child := range n1.Children {
		node := &widgets.TreeNode{}
		n2.Nodes = append(n2.Nodes, node)
		adapt(child, node)
	}
}

func actions(key string, tree *widgets.Tree) func() {
	// mostly vim style
	actions := map[string]func(){
		"k":      tree.ScrollUp,
		"<Up>":   tree.ScrollUp,
		"j":      tree.ScrollDown,
		"<Down>": tree.ScrollDown,
		"<C-b>":  tree.ScrollPageUp,
		"<C-u>":  tree.ScrollHalfPageUp,
		"<C-f>":  tree.ScrollPageDown,
		"<C-d>":  tree.ScrollHalfPageDown,
		"<Home>": tree.ScrollTop,
		"o":      tree.ToggleExpand,
		"G":      tree.ScrollBottom,
		"<End>":  tree.ScrollBottom,
		"E":      tree.ExpandAll,
		"C":      tree.CollapseAll,
		"<Resize>": func() {
			x, y := ui.TerminalDimensions()
			tree.SetRect(0, 0, x, y)
		},
	}

	return actions[key]
}

func listenKeyboard(tree *widgets.Tree, detail, help *widgets.Paragraph) {
	var previousKey string
	var previousSelected *Node

	visibilities := make(map[interface{}]bool)

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents

		switch e.ID {
		case "q", Quit:
			return
		case "g":
			if previousKey == "g" {
				tree.ScrollTop()
			}
		case "<Enter>", "<Space>":
			selected := extra[tree.SelectedNode()]
			detail.Text = selected.Detail

			selectionChanged := previousSelected != selected
			visibilities[detail] = selectionChanged || !visibilities[detail]

			previousSelected = selected
		case "?":
			visibilities[help] = !visibilities[help]
		default:
			if action := actions(e.ID, tree); action != nil {
				action()
			}
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		redraw(visibilities, tree, detail, help)
	}
}

func redraw(shouldShow map[interface{}]bool, tree *widgets.Tree, detail, help *widgets.Paragraph) {
	x, y := ui.TerminalDimensions()

	shouldDisplaySideBar := shouldShow[detail] || shouldShow[help]
	if shouldDisplaySideBar {
		tree.SetRect(0, 0, x*2/3, y)
	} else {
		tree.SetRect(0, 0, x, y)
	}

	if shouldShow[detail] && shouldShow[help] {
		detail.SetRect(x*2/3, 0, x, y/2)
		help.SetRect(x*2/3, y/2+1, x, y)
	} else if shouldShow[help] {
		detail.SetRect(0, 0, 0, 0)
		help.SetRect(x*2/3, 0, x, y)
	} else if shouldShow[detail] {
		detail.SetRect(x*2/3, 0, x, y)
		help.SetRect(0, 0, 0, 0)
	} else {
		help.SetRect(0, 0, 0, 0)
		detail.SetRect(0, 0, 0, 0)
	}

	ui.Render(tree, detail, help)
}
