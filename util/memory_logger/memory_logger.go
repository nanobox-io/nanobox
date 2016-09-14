package memory_logger

import (
	"bytes"
	"fmt"

	"github.com/jcelliott/lumber"
)

var Logger = &memoryLogger{}

type memoryLogger struct {
	prefix string
	buffer bytes.Buffer
}

func (m *memoryLogger) Dump() string {
	return m.buffer.String()
}

func (m *memoryLogger) Fatal(f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Fatal] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) Error(f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Error] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) Warn(f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Warn] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) Info(f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Info] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) Debug(f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Debug] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) Trace(f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Trace] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) IsFatal() bool {
	return true
}

func (m *memoryLogger) IsError() bool {
	return true
}

func (m *memoryLogger) IsWarn() bool {
	return true
}

func (m *memoryLogger) IsInfo() bool {
	return true
}

func (m *memoryLogger) IsDebug() bool {
	return true
}

func (m *memoryLogger) IsTrace() bool {
	return true
}

func (m *memoryLogger) GetLevel() int {
	return 0
}

func (m *memoryLogger) Print(n int, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Print] %s\n", m.prefix, fmt.Sprint(v...)))
}

func (m *memoryLogger) Printf(n int, f string, v ...interface{}) {
	m.buffer.WriteString(fmt.Sprintf("%s[Trace] %s\n", m.prefix, fmt.Sprintf(f, v...)))
}

func (m *memoryLogger) Level(n int) {
}

func (m *memoryLogger) Prefix(str string) {
	m.prefix = str	
}

func (m *memoryLogger) TimeFormat(string) {
}

func (m *memoryLogger) Close() {

}

func (m *memoryLogger) output(msg *lumber.Message) {

}

