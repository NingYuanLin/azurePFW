package logger

import (
	"io"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	a := "test"
	Debugf("this is Debug %s", a)
	Infof("this is Info %s", a)
	Warnf("this is Warn %s", a)
	Errorf("this is Error %s", a)

	Debug("this is Debug")
	Info("this is Info")
	Warn("this is Warn")
	Error("this is Error")

	SetDebugOutPutPath("logger_test.log")
	echoInfo := "this is Debug"
	Debug(echoInfo)
	file, err := os.OpenFile("logger_test.log", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal("Echo log to file error: the log file is not created successfully")
	}
	all, err := io.ReadAll(file)
	if err != nil {
		t.Fatal("Echo log to file error: the log file can not be read successfully")
	}
	for index := 0; index < len(echoInfo); index++ {
		b := echoInfo[index]
		if b != all[len(all)-len(echoInfo)+index-1] {
			t.Fatal("Echo log to file err: the file format is not right")
		}
	}

}
