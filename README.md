# dcrwages

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

`dcrwages` is a simple application to calculate average monthly USDT/DCR price and
write it to the terminal.

## How dcrwages works

`dcrwages` calculates the average monthly price of Decred in USDT using price data
collected from [Poloniex](https://poloniex.com). USDT/BTC and BTC/USD price history
for a month is downloaded, and the weighted average price over 15 minute periods is
used to find the monthly average USDT_DCR price.

`dcrwages` is currently hardcoded to collect price data between June 2016 and March 2019.
Data for DCR/BTC is not available before this time.

## How to use dcrwages

```bash
go build
./dcrwages
```

dcrwages was written and tested using Go 1.12.

Poloniex only allows 6 HTTP requests on it's API per second, so collecting the data
takes around 30 seconds.
