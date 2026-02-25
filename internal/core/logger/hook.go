package logger

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"server/internal/core/shared/types"
	"strconv"
	"strings"
)

type hook struct {
	file        string
	name        string
	writer      io.Writer
	entriesChan chan *log.Entry
}

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

func newHook(name string) log.Hook {
	file := filepath.Join("logs", fmt.Sprintf("%s.log", name))
	hookObj := &hook{
		name: name,
		file: file,
		writer: &lumberjack.Logger{
			Filename:   file,
			MaxSize:    50, // MB
			MaxBackups: 1000,
			MaxAge:     365, // days
			Compress:   true,
		},
		entriesChan: make(chan *log.Entry, 1000000),
	}
	go func(hookObj *hook) {
		for {
			entry := <-hookObj.entriesChan
			if hookObj.name != "http" {
				hookObj.print(entry)
			}
			hookObj.write(entry)
		}
	}(hookObj)
	return hookObj
}

func (h *hook) Levels() []log.Level {
	return log.AllLevels
}

func (h *hook) Fire(entry *log.Entry) error {
	if logLevel == types.LogLevelOff {
		return nil
	}
	dir := filepath.Dir(h.file)
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if entry.Level == log.FatalLevel || entry.Level == log.PanicLevel {
		h.print(entry)
		h.write(entry)
	} else {
		h.entriesChan <- entry
	}
	return nil
}

func (h *hook) format(entry *log.Entry, colors bool) string {
	tme := entry.Time.Format("2006-01-02 03:04:05 MST")
	pfx := fmt.Sprintf("%-"+strconv.Itoa(9)+"s", strings.ToUpper(entry.Level.String()))
	msg := regexp.MustCompile(`\s+`).ReplaceAllString(strings.Trim(entry.Message, "\n"), " ")
	if colors {
		var lvc int
		switch entry.Level {
		case log.InfoLevel:
			lvc = blue
		case log.WarnLevel:
			lvc = yellow
		case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
			lvc = red
		default:
			lvc = gray
		}
		llb := bytes.NewBuffer([]byte(""))
		if _, err := fmt.Fprintf(llb, " \x1b[%dm%s\x1b[0m", lvc, pfx); err != nil {
			return msg
		}
		pfx = llb.String()
	}
	return fmt.Sprintf("%s %s %s\n", tme, pfx, msg)
}

func (h *hook) print(entry *log.Entry) {
	fmt.Print(h.format(entry, true))
}

func (h *hook) write(entry *log.Entry) {
	_, _ = h.writer.Write([]byte(h.format(entry, false)))
}
