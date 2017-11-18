package glog

import "testing"

func TestLog(t *testing.T) {
	Init()
	//go run logging.go -log_dir=.
	log.Debugf("debug %s", Password("secret"))
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
}
