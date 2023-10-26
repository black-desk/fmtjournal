package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"

	"github.com/black-desk/fmtjournal/consts"
	"github.com/black-desk/fmtjournal/types"
)

var flags types.Flags

var rootCmd = &cobra.Command{
	Use:   "fmtjournal",
	Short: "A tool format journalctl json stream from stdin.",
	RunE: func(_ *cobra.Command, args []string) (err error) {
		return rootCmdRun(flags)
	},
}

func rootCmdRun(flags types.Flags) (err error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	var tmpl *template.Template
	tmpl = template.New("fmt").Option("missingkey=zero").
		Funcs(map[string]any{
			"indent": func(indent int, str string) (ret string) {
				indentStr := strings.Repeat("\t", indent)
				ret = strings.ReplaceAll(str, "\n", "\n"+indentStr)
				return
			},
		})

	var templateStr string
	templateStr, err = loadFmtFile(flags.FmtFile)
	if errors.Is(err, os.ErrNotExist) {
		templateStr = consts.DefaultFormat
	} else if err != nil {
		return
	}

	tmpl, err = tmpl.Parse(templateStr)
	if err != nil {
		return
	}

	go func() {
		dec := json.NewDecoder(os.Stdin)
		for {
			journalEntryMap := map[string]any{}

			if err := dec.Decode(&journalEntryMap); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			err := printWithMap(tmpl, journalEntryMap)
			if err != nil {
				panic(err)
			}
		}
		os.Exit(0)
	}()

	<-sigChan

	return
}

func loadFmtFile(path string) (ret string, err error) {
	if !filepath.IsAbs(path) {
		path, err = xdg.SearchConfigFile(
			filepath.Join("fmtjournal", path),
		)
		if err != nil {
			err = os.ErrNotExist
			return
		}
	}
	var file *os.File
	file, err = os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}

	var bs []byte
	bs, err = io.ReadAll(file)
	if err != nil {
		return
	}

	if len(bs) > 1 && bs[len(bs)-1] == '\n' {
		bs = bs[:len(bs)-1]
	}

	ret = string(bs)
	return
}

func printWithMap(tmpl *template.Template, journalEntryMap map[string]any) (err error) {
	journalExtraMap := map[string]any{}

	for k := range journalEntryMap {
		if isWellKnowField(k) {
			continue
		}

		journalExtraMap[k] = journalEntryMap[k]
	}

	if _, ok := journalEntryMap["PRIORITY"]; !ok {
		journalEntryMap["PRIORITY"] = strconv.FormatInt(
			int64(syslog.LOG_INFO),
			10,
		)
	}

	timestamp, err := strconv.ParseInt(
		journalEntryMap["__REALTIME_TIMESTAMP"].(string),
		10, 64)
	if err != nil {
		journalEntryMap["timestamp"] = journalEntryMap["__REALTIME_TIMESTAMP"].(string)
	} else {
		journalEntryMap["timestamp"] = time.UnixMicro(timestamp).String()
	}

	if len(journalExtraMap) == 0 {
		journalExtraMap = nil
	}

	journalEntryMap["extra"] = journalExtraMap

	err = tmpl.Execute(os.Stdout, journalEntryMap)
	if err != nil {
		return
	}

	return err
}

var wellKnowFields = map[string]struct{}{
	"MESSAGE":            {},
	"MESSAGE_ID":         {},
	"PRIORITY":           {},
	"CODE_FILE":          {},
	"CODE_LINE":          {},
	"CODE_FUNC":          {},
	"ERRNO":              {},
	"INVOCATION_ID":      {},
	"USER_INVOCATION_ID": {},
	"SYSLOG_FACILITY":    {},
	"SYSLOG_IDENTIFIER":  {},
	"SYSLOG_PID":         {},
	"SYSLOG_TIMESTAMP":   {},
	"SYSLOG_RAW":         {},
	"DOCUMENTATION":      {},
	"TID":                {},
	"UNIT":               {},
	"USER_UNIT":          {},

	"_PID":                       {},
	"_UID":                       {},
	"_GID":                       {},
	"_COMM":                      {},
	"_EXE":                       {},
	"_CMDLINE":                   {},
	"_CAP_EFFECTIVE":             {},
	"_AUDIT_SESSION":             {},
	"_AUDIT_LOGINUID":            {},
	"_SYSTEMD_CGROUP":            {},
	"_SYSTEMD_SLICE":             {},
	"_SYSTEMD_UNIT":              {},
	"_SYSTEMD_USER_UNIT":         {},
	"_SYSTEMD_USER_SLICE":        {},
	"_SYSTEMD_SESSION":           {},
	"_SYSTEMD_OWNER_UID":         {},
	"_SELINUX_CONTEXT":           {},
	"_SOURCE_REALTIME_TIMESTAMP": {},
	"_BOOT_ID":                   {},
	"_MACHINE_ID":                {},
	"_SYSTEMD_INVOCATION_ID":     {},
	"_HOSTNAME":                  {},
	"_TRANSPORT":                 {},
	"_STREAM_ID":                 {},
	"_LINE_BREAK":                {},
	"_NAMESPACE":                 {},
	"_RUNTIME_SCOPE":             {},

	"_KERNEL_DEVICE":    {},
	"_KERNEL_SUBSYSTEM": {},
	"_UDEV_SYSNAME":     {},
	"_UDEV_DEVNODE":     {},
	"_UDEV_DEVLINK":     {},

	"COREDUMP_UNIT":            {},
	"COREDUMP_USER_UNIT":       {},
	"OBJECT_PID":               {},
	"OBJECT_UID":               {},
	"OBJECT_GID":               {},
	"OBJECT_COMM":              {},
	"OBJECT_EXE":               {},
	"OBJECT_CMDLINE":           {},
	"OBJECT_AUDIT_SESSION":     {},
	"OBJECT_AUDIT_LOGINUID":    {},
	"OBJECT_SYSTEMD_CGROUP":    {},
	"OBJECT_SYSTEMD_SESSION":   {},
	"OBJECT_SYSTEMD_OWNER_UID": {},
	"OBJECT_SYSTEMD_UNIT":      {},
	"OBJECT_SYSTEMD_USER_UNIT": {},

	"__CURSOR":              {},
	"__REALTIME_TIMESTAMP":  {},
	"__MONOTONIC_TIMESTAMP": {},
}

func isWellKnowField(key string) bool {
	if strings.HasPrefix(key, "_") {
		return true
	}
	_, ok := wellKnowFields[key]
	return ok
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(
		&flags.FmtFile,
		"format-file", "c",
		"default",
		""+
			"Format string of journal logs store in a file, "+
			"last \\n ignored. "+
			"If path is relative, it will be searched in "+
			"`${XDG_CONFIG_HOME}/fmtjournal/`.\n"+
			"If that file not existed, "+
			"fallback to the builtin format.",
	)
}
