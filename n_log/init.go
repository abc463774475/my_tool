package n_log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type loginfo struct {
	file *os.File

	curDir   string
	curTotalLines int
	param    string

	lock sync.Mutex
}

var (
	g_log             *loginfo
	oneFileMaxLines int = 8192
	isWriteLog		bool = true
)

func ChangeWriteLog (b bool)  {
	isWriteLog = b
}

func init() {
	g_log = &loginfo{}
	str := filepath.Base(os.Args[0])

	processName := str[0 : len(str)-len(filepath.Ext(str))]

	//if len(os.Args) > 3 {
	//	processName = os.Args[1]
	//}
	if len(os.Args) >= 3 {
		isWriteLog = os.Args[2] == "1"
	}

	isWriteLog = true

	tdir := fmt.Sprintf("log_h/%s", processName)
	if _, err := os.Stat(tdir); err != nil {
		err = os.MkdirAll(tdir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("log init erro %v", err))
		}
	}

	g_log.curDir = tdir
}
