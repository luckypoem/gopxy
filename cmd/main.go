package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"gopxy"
	"os"
)

var cfg = flag.String("config", "./config.json", "config file")
var loglv = flag.Uint("enable_all_log", 4, "log level")

func main() {
	flag.Parse()
	logger := &logrus.Logger{}
	logger.SetLevel(logrus.Level(*loglv))
	logger.SetOutput(os.Stdout)
	txtfmt := &logrus.TextFormatter{}
	txtfmt.TimestampFormat = "2006-01-02 15:04:05"
	logger.SetFormatter(txtfmt)

	cfg, err := gopxy.Parse(*cfg)
	if err != nil {
		panic(err)
	}

	dir, _ := os.Getwd()
	logger.Infof("Load config finish:%+v, current pwd:%s", *cfg, dir)
	svr := gopxy.New(cfg, logger)
	logger.Info("Server start...")
	err = svr.Start()
	if err != nil {
		panic(err)
	}
}
