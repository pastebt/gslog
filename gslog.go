package gslog


import (
    "os"
    "fmt"
    "sync"
    "time"
)


const (
    DEBUG = iota + 1
    INFO
    WARNING
    ERROR
    FATAL
)


type Writer struct {
    path string
    file *os.File
    size int
    num  int
    m    *sync.Mutex
}


type Logger struct {
    name string
    w    *Writer
    lvl  int
    fmt    func(lv string, msg string) string
    Debug  func(v ...interface{})
    Info   func(v ...interface{})
    Warn   func(v ...interface{})
    Error  func(v ...interface{})
    Fatal  func(v ...interface{})
    Debugf func(fmt string, v ...interface{})
    Infof  func(fmt string, v ...interface{})
    Warnf  func(fmt string, v ...interface{})
    Errorf func(fmt string, v ...interface{})
    Fatalf func(fmt string, v ...interface{})
}


var loggerIdx = make(map[string]*Logger)


func WriterNew(fn string) (w *Writer) {
    w = &Writer{file:os.Stderr, m:&(sync.Mutex{})}
    if fn == "" { return }
    var err error
    w.file, err = os.OpenFile(fn, os.O_APPEND|os.O_CREATE, 0600)
    if err == nil {
        w.path = fn
        w.size = 1000000
        w.num = 9
    }
    return
}


func (w *Writer)SetSize(i int) *Writer { w.size = i; return w }
func (w *Writer)SetNum(i int) *Writer { w.num = i; return w }


func (l *Logger)SetWriter(w *Writer) *Logger { l.w = w; return l }
func (l *Logger)SetLeverl(i int) *Logger {
    if DEBUG <= i && i <= FATAL { l.lvl = i }
    return l
}
func (l *Logger)SetFmt(f func(lv string, msg string) string) *Logger {
    l.fmt = f
    return l
}


// get a new logger with name, write to stdout, with DEBUG level
func GetLogger(name string) (l *Logger) {
    l, ok := loggerIdx[name]
    if ok { return }
    w := WriterNew("")
    l = &Logger{name:name, w:w, lvl:DEBUG}
    loggerIdx[name] = l
    l.fmt = func (lv string, msg string) (ret string) {
        return time.Now().Format("2006-01-02 15:04:05") +
               " -" + l.name + "- " + lv + " - " + msg + "\n"
    }
    l.Debug  = l.getFunc(DEBUG,   "DEBUG")
    l.Info   = l.getFunc(INFO,    "INFO ")
    l.Warn   = l.getFunc(WARNING, "WARN ")
    l.Error  = l.getFunc(ERROR,   "ERROR")
    l.Fatal  = l.getFunc(FATAL,   "FATAL")
    l.Debugf = l.getFunf(DEBUG,   "DEBUG")
    l.Infof  = l.getFunf(INFO,    "INFO ")
    l.Warnf  = l.getFunf(WARNING, "WARN ")
    l.Errorf = l.getFunf(ERROR,   "ERROR")
    l.Fatalf = l.getFunf(FATAL,   "FATAL")
    return
}


func (l *Logger)getFunc(li int, lv string) func (v ...interface{}) {
    return func (v ...interface{}) {
        if li < l.lvl { return }
        msg := l.fmt(lv, fmt.Sprint(v...))
        l.w.log(msg)
    }
}


func (l *Logger)getFunf(li int, lv string) func (f string, v ...interface{}) {
    return func (f string, v ...interface{}) {
        if li < l.lvl { return }
        msg := l.fmt(lv, fmt.Sprintf(f, v...))
        l.w.log(msg)
    }
}


func (w *Writer)log(msg string) {
    // require lock
    w.m.Lock()
    // release lock
    defer w.m.Unlock()
    w.file.WriteString(msg)
}
