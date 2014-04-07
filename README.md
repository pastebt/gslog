#Gslog

Golang Simple Logger

##Install

```bash
$ go get github.com/pastebt/gslog
```

##Usage

###Simple one

```go
import "github.com/pastebt/gslog"

logger := gslog.GetLogger("")
logger.Debug("haha")
```
Will get a logger log message into os.Stderr with default format
```
2014-04-07 14:24:02 -- DEBUG   - haha
```

###Different Log name
```go
import "github.com/pastebt/gslog"

logger := gslog.GetLogger("second")
logger.Debug("this will be logged")
```
log message in os.Stderr with default format
```
2014-04-07 14:24:02 -second- INFO   - "this will be logged"
```

###Different Log level
```go
import "github.com/pastebt/gslog"

logger := gslog.GetLogger("").SetLevel(gslog.INFO)
logger.Debug("this will not be logged")
logger.Info("this will be logged")
```
log message in os.Stderr with default format
```
2014-04-07 14:24:02 -- INFO   - "this will be logged"
```

###Log into file
```go
import "github.com/pastebt/gslog"

logger := gslog.GetLogger("").SetWriter(WriterNew("/tmp/file.log"))
logger.Debug("haha")
```
then run ```cat /tmp/file.log``` you will see same log message

###Log into file, with custermized file size and number of keeped log file
```go
import "github.com/pastebt/gslog"

logger := gslog.GetLogger("").SetWriter(WriterNew("/tmp/file.log").SetSize(1000000).SetNum(5))
logger.Debug("haha")
```
default file size is 10M (10000000), number is 9, and 10k < size < 1G, 0 < num < 100
If you want the log file increase forever, SetNum(0)

###Customize log format
```go
import "time"
import "github.com/pastebt/gslog"

func cfmt(name string, level string, msg string) string {
    return time.Now().Format("2006/01/02 15:04:05") +
           " -" + name + "- " + lv + " : " + msg + "\n"
}
logger := gslog.GetLogger("").SetFmt(cfmt)
logger.Debug("haha")
```
You will see different log format from default
```
2014/04/07 14:24:02 -- DEBUG   : haha
```

