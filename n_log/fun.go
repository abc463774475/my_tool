package nlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func getStr(s string, v ...interface{}) string {
	stmp := fmt.Sprintf(s, v...)
	stmp = fmt.Sprintln(stmp)
	return stmp
}

func Info(str string, v ...interface{}) {
	appendLog("Info", getStr(str, v...), 2)
}

// layer 本身底层有2层 ， 比如 1 +2 = 3 调用函数的上层打印地址
func InfoWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	appendLog("Info", getStr(str, v...), layer)
}

func Debug(str string, v ...interface{}) {
	appendLog("Debug", getStr(str, v...), 2)
}

func DebugWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	appendLog("Debug", getStr(str, v...), layer)
}

func Panic(str string, v ...interface{}) {
	stmp := getStr(str, v...)
	appendLog("Erro", stmp, 2)
	panic(stmp)
}

func PanicWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	stmp := getStr(str, v...)
	appendLog("Erro", stmp, layer)
	panic(stmp)
}

func Erro(str string, v ...interface{}) {
	str = getStr(str, v...)
	appendLog("Erro", str, 2)
}

func ErroWithBack(str string, v ...interface{}) error {
	str = getStr(str, v...)
	appendLog("Erro", str, 2)
	return errors.New(str)
}

func ErroWithLayer(layer int, str string, v ...interface{}) error {
	layer += 2
	str = getStr(str, v...)
	appendLog("Erro", str, layer)
	return errors.New(str)
}

func OnlyFile(str string, v ...interface{}) {
	appendLog("OnlyFile", getStr(str, v...), 2)
}

func OnlyFileWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	appendLog("OnlyFile", getStr(str, v...), layer)
}

func JsonDebug(str string, v ...interface{}) {
	d, err := json.Marshal(v)
	if err != nil {
		ErroWithStack(3, "json marsharl not right  %v  %v", err, v)
		return
	}

	appendLog("Debug", getStr(str, string(d)), 2)
}

func JsonDebugWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	d, err := json.Marshal(v)
	if err != nil {
		ErroWithStack(layer+1, "json marsharl not right  %v  %v", err, v)
		return
	}

	appendLog("Debug", getStr(str, string(d)), layer)
}

func JsonInfo(str string, v ...interface{}) {
	d, err := json.Marshal(v)
	if err != nil {
		ErroWithStack(3, "json marsharl not right  %v  %v", err, v)
		return
	}

	appendLog("Info", getStr(str, string(d)), 2)
}

func JsonInfoWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	d, err := json.Marshal(v)
	if err != nil {
		ErroWithStack(layer+1, "json marsharl not right  %v  %v", err, v)
		return
	}

	appendLog("Info", getStr(str, string(d)), layer)
}

func JsonErro(str string, v ...interface{}) {
	d, err := json.Marshal(v)
	if err != nil {
		ErroWithStack(2, "json marsharl not right  %v  %v", err, v)
		return
	}

	appendLog("Erro", getStr(str, string(d)), 2)
}

func JsonErroWithLayer(layer int, str string, v ...interface{}) {
	layer += 2
	d, err := json.Marshal(v)
	if err != nil {
		ErroWithStack(layer+1, "json marsharl not right  %v  %v", err, v)
		return
	}

	appendLog("Erro", getStr(str, string(d)), layer)
}

func ErroWithStack(layer int, str string, v ...interface{}) {
	layer += 2
	str = getStr(str, v...)
	str = fmt.Sprintf("%v\nstack print  %v", str, string(debug.Stack()))

	appendLog("Erro", str, layer)
}

func appendLog(lName, s string, layer int) {
	var fileName string
	var lines int
	var bret bool
	var f1 uintptr

	curTimeStr := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	f1, fileName, lines, bret = runtime.Caller(layer)

	if !bret {
		fileName = "???"
		lines = 0
	}

	var sBack string

	funName := runtime.FuncForPC(f1).Name()

	compressType := _log.comressType

	if compressType == Quick {
		strs := strings.Split(fileName, "/")
		if len(strs) != 0 {
			fileName = strs[len(strs)-1]
		}
	}

	if compressType != Full {
		strs := strings.Split(funName, "/")
		if len(strs) != 0 {
			funName = strs[len(strs)-1]
		}
	}

	sBack = fmt.Sprintf("%s:%d : [%s]: %s func[%s]: %s",
		fileName, lines,
		lName, curTimeStr, funName, s)

	if !_log.isAsyn {
		flushLog(sBack, lName)
	} else {
		_log.c <- &flushData{
			sBack: sBack,
			lName: lName,
		}
	}
}

func flushLog(sBack string, lName string) {
	_lock.Lock()
	defer _lock.Unlock()

	if _log.file == nil {
		ltmp := time.Unix(time.Now().Unix(), 0).Format("2006-01-02_15_04_05")
		ltmp = fmt.Sprintf("./%s/%s%v.log", _log.curDir, ltmp, time.Now().Nanosecond())

		lf, err := os.Create(ltmp)
		if err != nil {
			str := fmt.Sprintf("logFile Init Erro %v %v ", ltmp, err)
			log.Println(str)
			return
		}
		_log.file = lf
	}

	if _log.isWriteLog || lName == "Erro" {
		_, _ = _log.file.Write([]byte(sBack))
	}

	// console color
	if lName == "Info" {
		sBack = fmt.Sprintf("\033[35m%s\033[0m", sBack)
	} else if lName == "Erro" {
		sBack = fmt.Sprintf("\033[31m%s\033[0m", sBack)
	} else if lName == "Debug" {
		sBack = fmt.Sprintf("\033[34m%s\033[0m", sBack)
	}

	if lName != "OnlyFile" {
		fmt.Print(sBack)
	}

	_log.curTotalLines++

	if _log.oneFileMaxLines != -1 && _log.curTotalLines >= _log.oneFileMaxLines {
		_log.file.Close()
		_log.curTotalLines = 0
		_log.file = nil
	}
}
