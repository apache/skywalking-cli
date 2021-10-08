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

package util

// TextFormat defines variables of style and colors.
const TextFormat = `
{{- $bold	:="\x1b[1m" -}}
{{- $black	:="\x1b[0;30m" -}}	{{- $Black	:="\x1b[1;30m" -}}	
{{- $red	:="\x1b[0;31m" -}}	{{- $Red	:="\x1b[1;31m" -}}
{{- $green	:="\x1b[0;32m" -}}	{{- $Green	:="\x1b[1;32m" -}}
{{- $yellow	:="\x1b[0;33m" -}}	{{- $Yellow	:="\x1b[1;33m" -}}
{{- $blue	:="\x1b[0;34m" -}}	{{- $Blue	:="\x1b[1;34m" -}}
{{- $purple	:="\x1b[0;35m" -}}	{{- $purple	:="\x1b[1;35m" -}}
{{- $cyan	:="\x1b[0;36m" -}}	{{- $Cyan	:="\x1b[1;36m" -}}
{{- $white	:="\x1b[0;37m" -}}	{{- $White	:="\x1b[1;37m" -}}
{{- $plain	:="\x1b[0m" -}}
`

const AppHelpTemplate = TextFormat + `{{$Green -}}
NAME:
	{{$green}}{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}{{$Green}}

USAGE:
	{{$green}}{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}
	{{- if .VisibleFlags}} [global options]{{end}}{{if .Commands}} command [command options]{{end}}
	{{- if .ArgsUsage}} {{.ArgsUsage}}{{else}} [arguments...]{{end}}{{end}}{{$Yellow}}{{if .Version}}{{if not .HideVersion}}

VERSION:
	{{$yellow}}{{.Version}}{{end}}{{end}}{{$Yellow}}{{if .Description}}

DESCRIPTION:
	{{$yellow}}{{.Description}}{{end}}{{$Yellow}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{$yellow}}{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{$Blue}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}

	{{.Name}}:{{range .VisibleCommands}}
	{{$Blue}}{{join .Names ", "}}{{"\t\t\t\t"}}{{$blue}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
	{{$Blue}}{{join .Names ", "}}{{"\t\t\t\t"}}{{$blue}}{{.Usage}}{{$plain}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}{{$Cyan}}

GLOBAL OPTIONS:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$cyan}}{{$option}}{{$plain}}{{end}}{{end}}{{$Green}}{{if .Copyright}}

COPYRIGHT:
	{{$green}}{{.Copyright}}{{end}}{{$plain}}
`

const CommandHelpTemplate = TextFormat + `{{$Green -}}
NAME:
	{{$green}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}}{{- if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{end}} - {{.Usage}}{{$Green}}

USAGE:
	{{$green}}{{if .UsageText}}{{.UsageText | nindent 2 | trim}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} 
	{{- if .ArgsUsage}} {{.ArgsUsage}}{{end}}{{end}}{{$Green}}{{if .Category}}

CATEGORY:
	{{$green}}{{.Category}}{{end}}{{$Blue}}{{if .Description}}

DESCRIPTION:
	{{$blue}}{{.Description}}{{end}}{{$Cyan}}{{if .VisibleFlags}}

OPTIONS:
	{{$cyan}}{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{$plain}}
`

const SubcommandHelpTemplate = TextFormat + `{{$Green -}}
NAME:
	{{$green}}{{.HelpName}} - {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}{{$Green}}

USAGE:
	{{$green}}{{if .UsageText}}{{.UsageText | nindent 2 | trim}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}}
	{{- if .ArgsUsage}} {{.ArgsUsage}}{{else}} [arguments...]{{end}}{{end}}{{$Blue}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}

	{{.Name}}:{{range .VisibleCommands}}
	{{$Blue}}{{join .Names ", "}}{{$blue}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
	{{$Blue}}{{join .Names ", "}}{{$blue}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}{{$Cyan}}

OPTIONS:
	{{$cyan}}{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{$plain}}
`
