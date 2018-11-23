package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/luizalabs/rey/checker"
	"github.com/luizalabs/rey/component"
	"github.com/luizalabs/rey/metric"
	"github.com/luizalabs/rey/notifier"
	"github.com/luizalabs/rey/runner"
)

type Config struct {
	Timeout           int    `envconfig:"checker_timeout" default:"5"`
	MaxRetry          int    `envconfig:"checker_max_retry" default:"3"`
	CircleInterval    int    `envconfig:"runner_circle_interval" default:"10"`
	ComponentsPath    string `envconfig:"components_path" default:"/etc/rey/components.json"`
	MetricsServerPort string `envconfig:"metrics_server_port" default:"5000"`
	NotifierToken     string `envconfig:"notifier_token"`
	NotifierUsername  string `envconfig:"notifier_username" default:"Doc. Rey"`
	NotifierAvatar    string `envconfig:"notifier_avatar" default:"https://bit.ly/2Sbf65n"`
	NotifierChannel   string `envconfig:"notifier_channel"`
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

	nt := notifier.New(conf.NotifierToken, conf.NotifierUsername, conf.NotifierAvatar, conf.NotifierChannel)
	cc := checker.New(conf.Timeout, conf.MaxRetry)
	r := runner.New(conf.CircleInterval, cc, nt)

	ms := metric.NewServer(conf.MetricsServerPort)

	go func() {
		log.Println("starting gauge metrics server")
		if err := ms.Run(); err != nil {
			log.Fatal("gauge metrics server error:", err)
		}
	}()

	go func() {
		log.Println("starting rey")
		if err := r.Run(ctx, compList); err != nil {
			log.Fatal("running error:", err)
		}
	}()

	<-exitChan
	r.Stop()
	ms.Stop(ctx)
	ctxCancel()
}
