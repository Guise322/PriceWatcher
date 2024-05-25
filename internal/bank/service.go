package bank

import (
	priceTime "PriceWatcher/internal/common/time"
	"PriceWatcher/internal/config"
	"PriceWatcher/internal/extractor"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	req  BankRequester
	ext  extractor.Extractor
	conf config.Config
}

func NewService(
	req BankRequester,
	ext extractor.Extractor,
	conf config.Config) Service {
	return Service{
		req:  req,
		ext:  ext,
		conf: conf,
	}
}

var bankUrl = "https://investzoloto.ru/gold-sber-oms/"

func (s Service) ServePrice() (message, subject string, err error) {
	logrus.Infof("Start processing a price")

	response, err := s.req.RequestPage(bankUrl)
	if err != nil {
		return "", "", fmt.Errorf("Cannot get a page with the current price: %w", err)
	}

	price, err := s.ext.ExtractPrice(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("Cannot extract the price from the body: %w", err)
	}

	msg := fmt.Sprintf("Курс золота. Продажа: %.2fр", price)
	sub := "Че по золоту?"

	return msg, sub, nil
}

func (s Service) GetWaitTime(now time.Time) time.Duration {
	variation := 1800
	randDur := priceTime.RandomSec(variation)
	callTime := getCallTime(now, s.conf.SendingHours)

	return getWaitDurWithRandomComp(now, callTime, randDur)
}

func (Service) PerStartDur(now time.Time) time.Duration {
	return priceTime.PerStartDur(now)
}