// SPDX-FileCopyrightText: 2025 Chen Linxuan <me@black-desk.cn>
//
// SPDX-License-Identifier: MIT

package consts

var (
	DefaultFormat = "" +
		`{{printf "%-38s " .timestamp}}` +

		`{{if eq .PRIORITY "0"}}` + "\033[1;31mEMERGENCY\033[0m" +
		`{{else if eq .PRIORITY "1"}}` + "\033[1;31mALERT    \033[0m" +
		`{{else if eq .PRIORITY "2"}}` + "\033[1;31mCRITICAL \033[0m" +
		`{{else if eq .PRIORITY "3"}}` + "\033[1;31mERROR    \033[0m" +
		`{{else if eq .PRIORITY "4"}}` + "\033[1;33mWARNING  \033[0m" +
		`{{else if eq .PRIORITY "5"}}` + "\033[1;34mNOTICE   \033[0m" +
		`{{else if eq .PRIORITY "6"}}` + "\033[1;34mINFO     \033[0m" +
		`{{else if eq .PRIORITY "7"}}` + "\033[1;90mDEBUG    \033[0m" +
		"{{end}}" +

		`{{if ne .SYSLOG_IDENTIFIER nil}}` +
		("" +
			`{{printf "%-8s " .SYSLOG_IDENTIFIER}}`) +
		`{{else}}` +
		("" +
			`{{printf "%-8s " ._COMM}}`) +
		`{{end}}` +
		`{{if ne ._PID nil}}` +
		("" +
			`{{$pid := printf "[%s]" ._PID}}` +
			`{{printf "%-6s" $pid}}`) +
		`{{end}}` +
		"\n" +

		`{{if ne .CODE_FILE nil}}` +
		("" +
			"{{.CODE_FILE}}" +
			`{{if ne .CODE_LINE nil}}` +
			("" +
				":{{.CODE_LINE}}" +
				`{{if ne .CODE_FUNC nil}}` +
				("" +
					" ({{.CODE_FUNC}})") +
				"{{end}}") +
			"{{end}}") +
		"\n{{end}}" +

		"\t" +
		`{{$message := print .MESSAGE }}` +
		`{{if eq .PRIORITY "0"}}` +
		("" +
			"\033[1;31m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "1"}}` +
		("" +
			"\033[1;31m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "2"}}` +
		("" +
			"\033[1;31m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "3"}}` +
		("" +
			"\033[1;31m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "4"}}` +
		("" +
			"\033[1;33m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "5"}}` +
		("" +
			"\033[1;34m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "6"}}` +
		("" +
			"\033[1;34m{{indent 1 $message}}\033[0m") +
		`{{else if eq .PRIORITY "7"}}` +
		("" +
			"\033[1;90m{{indent 1 $message}}\033[0m") +
		"{{else}}" +
		("" +
			"{{$message}}") +
		"{{end}}\n" +
		"{{if ne .extra nil}}" +
		("" +
			"{{range $k, $v := .extra}}" +
			("" +
				"{{$vStr := print $v}}" +
				"\t{{$k}}=\n" +
				"\t\t{{indent 2 $vStr}}\n") +
			"{{end}}") +
		"{{end}}" +
		"\n"
)
