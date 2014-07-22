package gslog

import (
    "testing"
    "time"
    "fmt"
)


func dbg(v ...interface{}) {
    fmt.Println(v...)
}


func TestLogger(tst *testing.T) {
    tst.Log("Logger")
    l := GetLogger("first")
    //l := GetLogger("")
    l.Debug("haha")
    l.Info("you can see this")
    l.SetLevel(INFO)
    l.Debug("you should not see this")
    l.Info("you can see this too")
}


func TestLoggerFile(tst *testing.T) {
    l := GetLogger("").SetWriter(WriterNew("/tmp/gslog.log").SetNum(3).SetSize(10002))
    //l.w.size = 10
    for i := 0; i < 10000; i++ {
        l.Debugf("test log %d", i * 1000000)
        l.Debug("test log")
        time.Sleep(30 * time.Millisecond)
    }
}
