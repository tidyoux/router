package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tidyoux/goutils"
	"github.com/tidyoux/goutils/cmd"
	"github.com/tidyoux/router/common"
	"github.com/tidyoux/router/common/db"
	"github.com/tidyoux/router/server"
	"github.com/tidyoux/router/server/config"
	"github.com/tidyoux/router/server/handler"
	"github.com/tidyoux/router/server/model"
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
		"routersrv",
		"routersrv is router server",
		"",
		run)
	c.Flags().StringVarP(&cfgFile, "config", "c", "app.yml", "config file (default is app.yml)")

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(*cmd.Command) error {
	serviceName := "routersvr"

	defer goutils.DeferRecover(serviceName, nil)()

	err := goutils.InitDaysJSONRotationLogger("./log/", serviceName+".log", 60)
	if err != nil {
		panic(err)
	}

	log.Infof("%s service start", serviceName)

	// Initial config.
	err = common.ReadConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig()

	// Initial db.
	dbInst, err := db.New(cfg.DSN)
	if err != nil {
		panic(err)
	}
	defer dbInst.Close()
	err = model.Init(dbInst)
	if err != nil {
		panic(fmt.Errorf("init model failed, %v", err))
	}

	// Initial server.
	err = server.Init(cfg)
	if err != nil {
		panic(fmt.Errorf("init server failed, %v", err))
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = goutils.NewLogWriter(log.Info)
	gin.DefaultErrorWriter = goutils.NewLogWriter(log.Error)

	r := gin.Default()
	v1 := r.Group("v1")
	err = handler.Init(v1)
	if err != nil {
		panic(fmt.Errorf("init handler failed, %v", err))
	}

	v1.Static("/", "./static")

	err = r.Run(cfg.Port)
	if err != nil {
		panic(err)
	}

	return nil
}
