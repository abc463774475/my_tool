package n_log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
	"encoding/json"
)

func getStr(s string, v ...interface{}) string {
	stmp := fmt.Sprintf(s, v...)
	stmp = fmt.Sprintln(stmp)
	return stmp
}

func Info(str string, v ...interface{}) {
	appendLog("Info", getStr(str, v...), 2)
}

func Info_cengci(cengci int,str string, v ...interface{}) {
	appendLog("Info", getStr(str, v...), cengci)
}

func Debug(str string, v ...interface{}) {
	appendLog("Debug", getStr(str, v...), 2)
}

func Debug_cengci(cengci int,str string, v ...interface{}) {
	appendLog("Debug", getStr(str, v...), cengci)
}

func Erro(str string, v ...interface{}) {
	appendLog("Erro", getStr(str, v...), 2)
}

func Erro_cengci(cengci int,str string, v ...interface{}) {
	appendLog("Erro", getStr(str, v...), cengci)
}

func Panic (str string,v ...interface{}){
	stmp := getStr(str, v...)
	appendLog("Erro",stmp , 2)
	panic(stmp)
}

func Panic_cengci (cengci int,str string,v ...interface{}){
	stmp := getStr(str, v...)
	appendLog("Erro",stmp , cengci)
	panic(stmp)
}


func ErroBack(str string, v ...interface{}) error {
	str = getStr(str, v...)
	appendLog("Erro", str, 2)
	return errors.New(str)
}

func ErroBack_cengci(cengci int,str string, v ...interface{}) error {
	str = getStr(str, v...)
	appendLog("Erro", str, cengci)
	return errors.New(str)
}


func OnlyFile(str string, v ...interface{}) {
	appendLog("OnlyFile", getStr(str, v...), 2)
}

func OnlyFile_cengci(cengci int,str string, v ...interface{}) {
	appendLog("OnlyFile", getStr(str, v...), cengci)
}


func PrintJsonDebug (str string, v ...interface{}){
	d, err := json.Marshal(v)
	if err != nil {
		Erro_special(3,"json marsharl not right  %v  %v",err,v)
		return
	}

	appendLog("Debug", getStr(str, string(d)), 2)
}

func PrintJsonDebug_cengci (cengci int,str string, v ...interface{}){
	d, err := json.Marshal(v)
	if err != nil {
		Erro_special(cengci+1,"json marsharl not right  %v  %v",err,v)
		return
	}

	appendLog("Debug", getStr(str, string(d)), cengci)
}

func PrintJsonInfo (str string, v ...interface{}){
	d, err := json.Marshal(v)
	if err != nil {
		Erro_special(3,"json marsharl not right  %v  %v",err,v)
		return
	}

	appendLog("Info", getStr(str, string(d)), 2)
}

func PrintJsonInfo_cengci (cengci int,str string, v ...interface{}){
	d, err := json.Marshal(v)
	if err != nil {
		Erro_special(cengci+1,"json marsharl not right  %v  %v",err,v)
		return
	}

	appendLog("Info", getStr(str, string(d)), cengci)
}

func PrintJsonErro (str string, v ...interface{}){
	d, err := json.Marshal(v)
	if err != nil {
		Erro_special(2,"json marsharl not right  %v  %v",err,v)
		return
	}

	appendLog("Erro", getStr(str, string(d)), 2)
}

func PrintJsonErro_cengci (cengci int,str string, v ...interface{}){
	d, err := json.Marshal(v)
	if err != nil {
		Erro_special(cengci + 1,"json marsharl not right  %v  %v",err,v)
		return
	}

	appendLog("Erro", getStr(str, string(d)), cengci)
}

func Erro_special(cengci int,str string, v ...interface{}) {
	str = getStr(str, v...)
	str = fmt.Sprintf("%v\nstack print  %v", str, string(debug.Stack()))

	appendLog("Erro", str, cengci)
}


//func Info_cengci(cengci int,str string, v ...interface{})  {
//	str = getStr(str,v...)
//	appendLog("Info",str,cengci)
//}

func appendLog(lName, s string, cengci int) {
	g_log.lock.Lock()
	defer g_log.lock.Unlock()

	if g_log.file == nil {
		ltmp := time.Unix(time.Now().Unix(), 0).Format("2006-01-02_15_04_05")
		ltmp = fmt.Sprintf("./%s/%s%v.log", g_log.curDir, ltmp, time.Now().Nanosecond())

		lf, err := os.Create(ltmp)
		if err != nil {
			str := fmt.Sprintf("logFile Init Erro %v %v ", ltmp, err)
			log.Println(str)
			return
			//panic(str)
		}
		g_log.file = lf
	}

	var fileName string
	var lines int
	var bret bool
	var f1 uintptr

	curTimeStr := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	f1, fileName, lines, bret = runtime.Caller(cengci)

	if !bret {
		fileName = "???"
		lines = 0
	}

	var sBack string

	if len(g_log.param) == 0 {
		sBack = fmt.Sprintf("%s:%d : [%s]: %s func[%s]: %s",
			fileName, lines,
			lName, curTimeStr, runtime.FuncForPC(f1).Name(), s)
	} else {
		sBack = fmt.Sprintf("%s:%d : [%s]: %s func[%s]: %s last log timeOut \n%v\n",
			fileName, lines,
			lName, curTimeStr, runtime.FuncForPC(f1).Name(), s, g_log.param)

		g_log.param = ""
	}

	if isWriteLog || lName == "Erro"{
		g_log.file.Write([]byte(sBack))
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

	g_log.curTotalLines++

	if g_log.curTotalLines >= oneFileMaxLines {
		g_log.file.Close()
		g_log.curTotalLines = 0
		g_log.file = nil
	}
}
