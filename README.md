# dcrwages

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

`dcrwages` calculates average monthly USDT/DCR price and writes it to the terminal.
It can also serve the price information over a RESTful HTTP interface.

## How dcrwages works

`dcrwages` calculates the average monthly price of Decred in USDT using price data
collected from [Poloniex](https://poloniex.com). USDT/BTC and BTC/USD price history
is downloaded, and the weighted average price over 15 minute periods is used to find
the monthly average USDT_DCR price.

`dcrwages` is currently hardcoded to collect price data starting with the current
month and working back to June 2016.
Data for DCR/BTC is not available before this time.

dcrwages was written with Go 1.12.

## How to use dcrwages

```bash
# Build the executable
env GO111MODULE=on go install .

# Run the executable
$(go env GOPATH)/bin/dcrwages
```

The process will begin contacting Poloniex and downloading price information.
It will print the price information out to the console as it downloads.
Poloniex only allows 6 HTTP requests on it's API per second, so collecting
all of the data takes around 30 seconds.

Once the data is obtained, the API server will start listening on port 3000

<http://localhost:3000>
