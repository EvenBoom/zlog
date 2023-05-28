package zlog

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

var debug = false
var logPath = "./log"
var appPath = logPath + "/app/"
var infoPath = logPath + "/info/"
var warnPath = logPath + "/warn/"
var errorPath = logPath + "/error/"
var panicPath = logPath + "/panic/"

// 当前时间
var curTime struct {
	dateFile string
	dateTime string
	l        sync.RWMutex
}

// SetPath 设置路径
func SetPath(path string) {
	logPath = path
}

var digitals = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var infoLogger struct {
	w io.Writer
}

var warnLogger struct {
	w io.Writer
}

var errorLogger struct {
	w io.Writer
}

var panicLogger struct {
	w io.Writer
}

// Default 不会打印日志到标准输出
func Default() {
	initZlog()
}

// Debug 打印日志到标准输出
func Debug() {
	debug = true
	initZlog()
}

func initZlog() {
	now := time.Now()
	curTime.dateTime = now.Format("2006/01/02 15:04:05")
	curTime.dateFile = now.Format("20060102.log")
	initDir()
	initInfo()
	initWarn()
	initError()
	initPanic()
	go logTimer()
}

func logTimer() {
	ticker := time.NewTicker(time.Second)
	for t := range ticker.C {
		curTime.l.Lock()
		curTime.dateTime = t.Format("2006/01/02 15:04:05")
		d := t.Format("20060102.log")
		if curTime.dateFile != d {
			curTime.dateFile = d

			appOut, err := os.OpenFile(appPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			infoOut, err := os.OpenFile(infoPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			warnOut, err := os.OpenFile(warnPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			errorOut, err := os.OpenFile(errorPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			panicOut, err := os.OpenFile(panicPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			if debug {
				infoLogger.w = io.MultiWriter(appOut, infoOut, os.Stdout)
				warnLogger.w = io.MultiWriter(appOut, warnOut, os.Stdout)
				errorLogger.w = io.MultiWriter(appOut, errorOut, os.Stdout)
				panicLogger.w = io.MultiWriter(appOut, panicOut, os.Stdout)
			} else {
				infoLogger.w = io.MultiWriter(appOut, infoOut)
				warnLogger.w = io.MultiWriter(appOut, warnOut)
				errorLogger.w = io.MultiWriter(appOut, errorOut)
				panicLogger.w = io.MultiWriter(appOut, panicOut)
			}

		}
		curTime.l.Unlock()
	}
}

func initDir() {

	err := os.MkdirAll(appPath, os.ModeDir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(infoPath, os.ModeDir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(warnPath, os.ModeDir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(errorPath, os.ModeDir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(panicPath, os.ModeDir)
	if err != nil {
		panic(err)
	}

}

func initInfo() {

	appOut, err := os.OpenFile(appPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(infoPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	if debug {
		infoLogger.w = io.MultiWriter(appOut, out, os.Stdout)
	} else {
		infoLogger.w = io.MultiWriter(appOut, out)
	}

}

func initWarn() {

	appOut, err := os.OpenFile(appPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(warnPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	if debug {
		warnLogger.w = io.MultiWriter(appOut, out, os.Stdout)
	} else {
		warnLogger.w = io.MultiWriter(appOut, out)
	}

}

func initError() {

	appOut, err := os.OpenFile(appPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(errorPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	if debug {
		errorLogger.w = io.MultiWriter(appOut, out, os.Stdout)
	} else {
		errorLogger.w = io.MultiWriter(appOut, out)
	}

}

func initPanic() {

	appOut, err := os.OpenFile(appPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(panicPath+curTime.dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	if debug {
		panicLogger.w = io.MultiWriter(appOut, out, os.Stdout)
	} else {
		panicLogger.w = io.MultiWriter(appOut, out)
	}

}

// Info 信息
func Info(msg string) {

	curTime.l.RLock()
	dateTime := curTime.dateTime
	w := infoLogger.w
	curTime.l.RUnlock()

	var callers []logCaller
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = append(callers, logCaller{file: file, line: line})
	}
	// 异步
	go func() {
		buf := bytes.Buffer{}
		buf.WriteString("[INFO] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')
		for _, caller := range callers {
			buf.WriteString("[LINE] ")
			buf.WriteString(caller.file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for caller.line > 0 {
				b = append(b, digitals[caller.line%10])
				caller.line = caller.line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
		}
		w.Write(buf.Bytes())
	}()

}

// Warn 警告
func Warn(msg string) {

	curTime.l.RLock()
	dateTime := curTime.dateTime
	w := warnLogger.w
	curTime.l.RUnlock()

	var callers []logCaller
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = append(callers, logCaller{file: file, line: line})
	}
	// 异步
	go func() {
		buf := bytes.Buffer{}
		buf.WriteString("[WARN] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')
		for _, caller := range callers {
			buf.WriteString("[LINE] ")
			buf.WriteString(caller.file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for caller.line > 0 {
				b = append(b, digitals[caller.line%10])
				caller.line = caller.line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
		}

		w.Write(buf.Bytes())
	}()
}

// Error 错误
func Error(msg string) {

	curTime.l.RLock()
	dateTime := curTime.dateTime
	w := errorLogger.w
	curTime.l.RUnlock()

	var callers []logCaller
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = append(callers, logCaller{file: file, line: line})
	}
	// 异步
	go func() {
		buf := bytes.Buffer{}
		buf.WriteString("[EROR] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')
		for _, caller := range callers {
			buf.WriteString("[LINE] ")
			buf.WriteString(caller.file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for caller.line > 0 {
				b = append(b, digitals[caller.line%10])
				caller.line = caller.line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
		}
		w.Write(buf.Bytes())
	}()
}

// Panic 恐慌
func Panic(msg string) {

	curTime.l.RLock()
	dateTime := curTime.dateTime
	w := panicLogger.w
	curTime.l.RUnlock()

	var callers []logCaller
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = append(callers, logCaller{file: file, line: line})
	}
	// 异步
	go func() {
		buf := bytes.Buffer{}
		buf.WriteString("[PANC] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')
		for _, caller := range callers {
			buf.WriteString("[LINE] ")
			buf.WriteString(caller.file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for caller.line > 0 {
				b = append(b, digitals[caller.line%10])
				caller.line = caller.line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
		}
		w.Write(buf.Bytes())
		panic(msg)
	}()
}

// 日志调用函数
type logCaller struct {
	file string
	line int
}
