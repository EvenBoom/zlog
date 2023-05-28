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
var dateFile = ""
var dateTime = ""
var logPath = "./log"
var appPath = logPath + "/app/"
var infoPath = logPath + "/info/"
var warnPath = logPath + "/warn/"
var errorPath = logPath + "/error/"
var panicPath = logPath + "/panic/"

// SetPath 设置路径
func SetPath(path string) {
	logPath = path
}

var digitals = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var infoLogger struct {
	w io.Writer
	l sync.Mutex
}

var warnLogger struct {
	w io.Writer
	l sync.Mutex
}

var errorLogger struct {
	w io.Writer
	l sync.Mutex
}

var panicLogger struct {
	w io.Writer
	l sync.Mutex
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
	dateTime = now.Format("2006/01/02 15:04:05")
	dateFile = now.Format("20060102.log")
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
		dateTime = t.Format("2006/01/02 15:04:05")
		d := t.Format("20060102.log")
		if dateFile != d {
			dateFile = d

			appOut, err := os.OpenFile(appPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			infoOut, err := os.OpenFile(infoPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			warnOut, err := os.OpenFile(warnPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			errorOut, err := os.OpenFile(errorPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}

			panicOut, err := os.OpenFile(panicPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	appOut, err := os.OpenFile(appPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(infoPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	appOut, err := os.OpenFile(appPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(warnPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	appOut, err := os.OpenFile(appPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(errorPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	appOut, err := os.OpenFile(appPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(panicPath+dateFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	// 异步
	go func() {
		infoLogger.l.Lock()

		buf := bytes.Buffer{}
		buf.WriteString("[INFO] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')

		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			buf.WriteString("[LINE] ")
			buf.WriteString(file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for line > 0 {
				b = append(b, digitals[line%10])
				line = line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
			if !ok {
				break
			}
		}
		infoLogger.w.Write(buf.Bytes())
		infoLogger.l.Unlock()
	}()

}

// Warn 警告
func Warn(msg string) {
	// 异步
	go func() {
		warnLogger.l.Lock()

		buf := bytes.Buffer{}
		buf.WriteString("[WARN] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')

		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			buf.WriteString("[LINE] ")
			buf.WriteString(file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for line > 0 {
				b = append(b, digitals[line%10])
				line = line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
			if !ok {
				break
			}
		}
		warnLogger.w.Write(buf.Bytes())
		warnLogger.l.Unlock()
	}()
}

// Error 错误
func Error(msg string) {
	// 异步
	go func() {
		errorLogger.l.Lock()

		buf := bytes.Buffer{}
		buf.WriteString("[EROR] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')

		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			buf.WriteString("[LINE] ")
			buf.WriteString(file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for line > 0 {
				b = append(b, digitals[line%10])
				line = line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
			if !ok {
				break
			}
		}
		errorLogger.w.Write(buf.Bytes())
		errorLogger.l.Unlock()
	}()
}

// Panic 恐慌
func Panic(msg string) {
	// 异步
	go func() {
		panicLogger.l.Lock()

		buf := bytes.Buffer{}
		buf.WriteString("[PANC] ")
		buf.WriteString(dateTime)
		buf.WriteByte(' ')
		buf.WriteString(msg)
		buf.WriteByte('\n')

		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			buf.WriteString("[LINE] ")
			buf.WriteString(file)
			buf.WriteByte(':')
			b := make([]byte, 0, 9)
			for line > 0 {
				b = append(b, digitals[line%10])
				line = line / 10
			}

			for i := len(b) - 1; i >= 0; i-- {
				buf.WriteByte(b[i])
			}

			buf.WriteByte('\n')
			if !ok {
				break
			}
		}
		panicLogger.w.Write(buf.Bytes())
		panicLogger.l.Unlock()
		panic(msg)
	}()
}
