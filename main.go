package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/luizalabs/rey/aggregator"
	"github.com/luizalabs/rey/checker"
	"github.com/luizalabs/rey/component"
	"github.com/luizalabs/rey/runner"
)

type Config struct {
	ApiID          string `envconfig:"aggregator_api_id"`
	ApiKey         string `envconfig:"aggregator_api_key"`
	Timeout        int    `envconfig:"checker_timeout" default:"5"`
	MaxRetry       int    `envconfig:"checker_max_retry" default:"3"`
	CircleInterval int    `envconfig:"runner_circle_interval" default:"10"`
	ComponentsPath string `envconfig:"components_path" default:"/etc/rey/components.json"`
}

func main() {
	conf := new(Config)
	if err := envconfig.Process("rey", conf); err != nil {
		log.Fatal("processing configs:", err)
	}

	compList, err := component.GetList(conf.ComponentsPath)
	if err != nil {
		log.Fatal("getting component list:", err)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	defer close(exitChan)

	ctx, ctxCancel := context.WithCancel(context.Background())

	cc := checker.New(conf.Timeout, conf.MaxRetry)
	ag := aggregator.New(conf.ApiID, conf.ApiKey)
	r := runner.New(conf.CircleInterval, cc, ag)

	go func() {
		log.Println("starting rey")
		if err := r.Run(ctx, compList); err != nil {
			log.Fatal("running error:", err)
		}
	}()

	<-exitChan
	r.Stop()
	ctxCancel()
}
