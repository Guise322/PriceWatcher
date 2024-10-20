package bank

import (
	"PriceWatcher/internal/entities/bank"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

type BankRequester struct{}

func (r BankRequester) RequestPage() (bank.Response, error) {
	u := launcher.New().Bin("/usr/bin/chromium").MustLaunch()
	browser := rod.New().ControlURL(u).Timeout(time.Minute).MustConnect()
	browser.MustIgnoreCertErrors(true)
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.sberbank.ru/ru/quotes/metalbeznal")
	time.Sleep(20 * time.Second)

	html, err := page.HTML()
	if err != nil {
		return bank.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", err)
	}

	respReader := strings.NewReader(html)

	return bank.Response{Body: respReader}, nil
}
