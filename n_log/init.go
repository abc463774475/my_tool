package nlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	_log  *loginfo
	_lock sync.Mutex
)

const DEFAULTMAXLINE = -1

type CompressType int

const (
	// 全日志
	Full CompressType = 0 + iota
	// 精简
	Easy
	// 全速
	Quick
)

type flushData struct {
	sBack string
	lName string
}

type loginfo struct {
	file *os.File

	curDir        string
	curTotalLines int

	l sync.Mutex

	oneFileMaxLines int
	isWriteLog      bool

	comressType CompressType

	// 是否异步
	// 暂时不推荐 异步模式，会牵涉日志 shutdown的问题
	isAsyn bool

	c chan *flushData
}

func init() {
	_log = &loginfo{
		file:            nil,
		curDir:          "",
		curTotalLines:   0,
		oneFileMaxLines: DEFAULTMAXLINE,
		l:               sync.Mutex{},
		comressType:     Full,
		isWriteLog:      true,
		isAsyn:          false,
	}
	str := filepath.Base(os.Args[0])

	processName := str[0 : len(str)-len(filepath.Ext(str))]

	dir := fmt.Sprintf("log_h/%s", processName)
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("log init erro %v", err))
		}
	}

	_log.curDir = dir
}

// 不想用默认的log，则改用这种方式
func InitLog(options ...Option) {
	ptemp := &loginfo{
		file:            nil,
		curDir:          "",
		curTotalLines:   0,
		l:               sync.Mutex{},
		oneFileMaxLines: DEFAULTMAXLINE,
		isWriteLog:      true,
		comressType:     Full,
		isAsyn:          false,
	}

	for _, v := range options {
		v.apply(ptemp)
	}

	if ptemp.isAsyn {
		ptemp.c = make(chan *flushData, 10000)

		// 暂时不推荐 异步模式，会牵涉日志 shutdown的问题
		go func() {
			for {
				d, ok := <-ptemp.c
				if !ok {
					return
				}
				flushLog(d.sBack, d.lName)
			}
		}()
	}

	_lock.Lock()
	defer _lock.Unlock()
	if _log.file != nil {
		_log.file.Close()
	}
	if _log.c != nil {
		close(_log.c)
	}

	ptemp.curDir = _log.curDir

	_log = ptemp
}
