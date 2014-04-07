package gslog

import "testing"


func TestLogger(tst *testing.T) {
    tst.Log("Logger")
    //l := GetLogger("first")
    l := GetLogger("")
    l.Debug("haha")
}
