package bank

import (
	"PriceWatcher/internal/entities/page"
	"fmt"
	"net/http"
)

type BankRequester struct{}

func (r BankRequester) RequestPage(url string) (page.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return page.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", err)
	}

	return page.Response{Body: resp.Body}, nil
}
