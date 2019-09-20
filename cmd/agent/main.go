package main

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidyoux/goutils"
	"github.com/tidyoux/goutils/cmd"
	"github.com/tidyoux/goutils/service"
	"github.com/tidyoux/router/agent"
	"github.com/tidyoux/router/agent/config"
	"github.com/tidyoux/router/common"
)

var (
	cfgFile string
)

func init() {
	rand.Seed(time.Now().Unix())
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	c := cmd.New(
		"routeragt",
		"routeragt is router agent",
		"",
		run)
	c.Flags().StringVarP(&cfgFile, "config", "c", "app.yml", "config file (default is app.yml)")

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(*cmd.Command) error {
	serviceName := "routeragt"

	defer goutils.DeferRecover(serviceName, nil)()

	err := goutils.InitDaysJSONRotationLogger("./log/", serviceName+".log", 60)
	if err != nil {
		panic(err)
	}

	log.Infof("%s service start", serviceName)

	err = common.ReadConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	svr := service.NewWithInterval(agent.New(cfg), time.Second)
	if err := svr.Start(); err != nil {
		panic(err)
	}

	return nil
}
