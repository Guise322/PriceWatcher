package main

import (
	"PriceWatcher/internal/app"
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/domain/message"
	"PriceWatcher/internal/infrastructure/configer"
	"PriceWatcher/internal/infrastructure/sender"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	serv, err := newService()
	if err != nil {
		logrus.Errorf("Got the error: %v", err)
	}

	logrus.Infoln("Start the application")

	serv.Watch(ctx.Done(), cancel, clock.RealClock{})

	logrus.Infoln("The application is done")
}

func newService() (*app.PriceService, error) {
	sen := sender.Sender{}
	val := message.MessageHourVal{}
	conf := configer.Configer{}

	return app.NewPriceService(sen, val, conf)
}
