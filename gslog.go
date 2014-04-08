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


func WriterNew(fn string) (w *Writer) {
    w = &Writer{file:os.Stderr, m:&(sync.Mutex{})}
    if fn == "" { return }
    var err error
    w.file, err = os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    if err == nil {
        w.path = fn
        w.size = 10000000
        w.num = 9
    }
    return
}

// limit one log file size between 10k to 1G
func (w *Writer)SetSize(i int) *Writer {
    if 10000 < i && i < 1000000000 { w.size = i; }
    return w
}

// limit total log file < 100
func (w *Writer)SetNum(i int) *Writer {
    if 0 <= i && i < 100 { w.num = i;}
    return w
}


func (w *Writer)rotate(s int, d int) (err error) {
    src := w.path
    if s > 0 { src = fmt.Sprintf("%s.%d", src, s) }
    if s >= w.num {
        err = os.Remove(src)
        return
    }
    dst := fmt.Sprintf("%s.%d", w.path, d)
    if _, err = os.Stat(dst); err == nil {
        if err = w.rotate(d, d + 1); err != nil {
            err = os.Rename(src, dst)
        }
    }
    return
}


func (w *Writer)log(msg string) {
    w.m.Lock()
    defer w.m.Unlock()
    // if has path and num > 0, means we need file and rotate
    if len(w.path) > 0 && w.num > 0 {
        s, e := w.file.Stat()
        if e == nil &&  int64(len(msg)) + s.Size() > int64(w.size) {
            if err :=w.rotate(0, 1); err == nil {
                w.file.Close()
                w.file, _ = os.OpenFile(w.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
            }
        }
    }
    w.file.WriteString(msg)
    w.file.Sync()
}


type Logger struct {
    name   string
    w      *Writer
    lvl    int
    fmt    func(name string, lv string, msg string) string
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


func (l *Logger)SetWriter(w *Writer) *Logger {
    l.w = w;
    return l
}

func (l *Logger)SetLevel(i int) *Logger {
    if DEBUG <= i && i <= FATAL { l.lvl = i }
    return l
}

func (l *Logger)SetFmt(f func(name string, lv string, msg string) string) *Logger {
    l.fmt = f
    return l
}


//var loggerIdx = make(map[string]*Logger)
var loggerIdx = struct {
            data map[string]*Logger
            m *sync.Mutex
           } {make(map[string]*Logger),
              &sync.Mutex{}}


// get a new logger with name, write to stdout, with DEBUG level
func GetLogger(name string) (l *Logger) {
    loggerIdx.m.Lock()
    defer loggerIdx.m.Unlock()
    l, ok := loggerIdx.data[name]
    if ok { return }
    w := WriterNew("")
    l = &Logger{name:name, w:w, lvl:DEBUG}
    loggerIdx.data[name] = l
    l.fmt = func (name string, lv string, msg string) (ret string) {
        return time.Now().Format("2006-01-02 15:04:05") +
               " -" + name + "- " + lv + " - " + msg + "\n"
    }
    l.Debug  = l.getFunc(DEBUG,   "DEBUG  ")
    l.Info   = l.getFunc(INFO,    "INFO   ")
    l.Warn   = l.getFunc(WARNING, "WARNING")
    l.Error  = l.getFunc(ERROR,   "ERROR  ")
    l.Fatal  = l.getFunc(FATAL,   "FATAL  ")
    l.Debugf = l.getFunf(DEBUG,   "DEBUG  ")
    l.Infof  = l.getFunf(INFO,    "INFO   ")
    l.Warnf  = l.getFunf(WARNING, "WARNING")
    l.Errorf = l.getFunf(ERROR,   "ERROR  ")
    l.Fatalf = l.getFunf(FATAL,   "FATAL  ")
    return
}


func (l *Logger)getFunc(li int, lv string) func (v ...interface{}) {
    return func (v ...interface{}) {
        if li < l.lvl { return }
        msg := l.fmt(l.name, lv, fmt.Sprint(v...))
        l.w.log(msg)
    }
}


func (l *Logger)getFunf(li int, lv string) func (f string, v ...interface{}) {
    return func (f string, v ...interface{}) {
        if li < l.lvl { return }
        msg := l.fmt(l.name, lv, fmt.Sprintf(f, v...))
        l.w.log(msg)
    }
}
