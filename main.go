package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/go-chat/gochat/config"
	"github.com/go-chat/gochat/handler"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cfg        *config.Config
	configFile = kingpin.Flag("config-file", "Path to config file").Short('c').Default("").String()
)

func main() {
	var err error
	kingpin.Parse()

	cfg, err = config.New(*configFile)
	if err != nil {
		logrus.WithError(err).Error("Load config file failed, use default config")
		cfg = config.GetDefault()
	}

	logrus.SetLevel(logrus.InfoLevel)

	handler.NewServer(cfg)

	handler.Serve(cfg)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	handler.StopServer()
}
