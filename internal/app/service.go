package app

import (
	"PriceWatcher/internal/app/clock"
	"PriceWatcher/internal/app/interrupt"
	"PriceWatcher/internal/domain/hour"
	"PriceWatcher/internal/domain/page"
	"PriceWatcher/internal/interfaces/configer"
	interReq "PriceWatcher/internal/interfaces/requester"
	interSend "PriceWatcher/internal/interfaces/sender"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type GoldPriceService struct {
	req    interReq.Requester
	sender interSend.Sender
	ext    page.Extractor
	val    hour.HourValidator
	conf   configer.Configer
}

func NewGoldPriceService(
	req interReq.Requester,
	sender interSend.Sender,
	ext page.Extractor,
	val hour.HourValidator,
	conf configer.Configer) *GoldPriceService {

	serv := GoldPriceService{
		req:    req,
		sender: sender,
		ext:    ext,
		val:    val,
		conf:   conf,
	}

	return &serv
}

func (s *GoldPriceService) serve(clock clock.Clock) error {
	curHour := clock.Now().Hour()

	conf, err := s.conf.GetConfig()
	if err != nil {
		return fmt.Errorf("on getting the config an error occurs: %w", err)
	}

	logrus.Infof("Check time for processing a gold price. The time value: %v", curHour)

	if !s.val.Validate(curHour, conf.SendingHours) {
		logrus.Info("It is not appropriate time for getting a price")

		return nil
	}

	logrus.Info("Start processing a gold price")

	response, err := s.req.RequestPage()
	if err != nil {
		return fmt.Errorf("cannot get a page with the current price of gold: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return fmt.Errorf("cannot extract the gold price from the body: %w", err)
	}

	err = s.sender.Send(price, conf.Email)
	if err != nil {
		return fmt.Errorf("cannot send the gold price: %w", err)
	}

	logrus.Info("The gold price is processed")

	return nil
}

func (s *GoldPriceService) Watch(done <-chan struct{}, cancel context.CancelFunc, clock clock.Clock) {
	interrupt.WatchForInterruption(cancel)

	errMes := "An error occurs while serving a gold price: %v"
	if err := s.serve(clock); err != nil {
		logrus.Errorf(errMes, err)
	}

	err := waitHourStart(clock.Now())
	if err != nil {
		logrus.Errorf("An error occurs while waiting when the next hour begins: %v", err)
	}

	t := time.NewTicker(1 * time.Hour)
	defer t.Stop()

	for {
		select {
		case <-done:
			logrus.Info("Shut down the application")
			return
		case <-t.C:
			if err := s.serve(clock); err != nil {
				logrus.Errorf(errMes, err)
			}
		}
	}
}
