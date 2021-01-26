package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fasthttp/router"
	awpapi "github.com/kAbramenko/AnonWallPoster/api/awp"
	awpconf "github.com/kAbramenko/AnonWallPoster/internal/awp"
	"github.com/valyala/fasthttp"
)

func main() {
	if err := awpconf.Parse(); err != nil {
		log.Fatal(err)
	}
	configure_logger()
	if pwd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(pwd)
	}
	log.Debug("Config:", awpconf.Cfg)

	r := router.New()

	r.GET("/", awpapi.Index)
	r.GET("/posts/{id:[0-9]+}", awpapi.Index)

	r.POST("/api/post", awpapi.Post)

	r.ANY("/", awpapi.BadRequest)

	r.ServeFiles("/static/{filepath:*}", "./web/static")
	srv := &fasthttp.Server{
		Handler:            r.Handler,
		Logger:             log.StandardLogger(),
		MaxRequestBodySize: awpconf.Cfg.BodySize * 1024 * 1024,
	}
	srv.Name = awpconf.Cfg.Name

	log.Println("Start listen: ", awpconf.Cfg.Address)
	if err := srv.ListenAndServe(awpconf.Cfg.Address); err != nil {
		log.Fatal(err)
	}
}

func configure_logger() {
	switch awpconf.Cfg.Log {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Fatal("Unknown log level:", awpconf.Cfg.Log)
	}
	log.Info("Log level:", log.GetLevel())
}
