# dcrwages

[![Build Status](https://travis-ci.org/jholdstock/dcrwages.png?branch=master)](https://travis-ci.org/jholdstock/dcrwages)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

dcrwages calculates average monthly USDT/DCR price and writes it to the terminal.
It then starts a web server with a HTML page displaying all of the price information,
as well as providing the information over a RESTful HTTP interface.

## How dcrwages works

dcrwages calculates the average monthly price of Decred in USDT using price data
collected from [Poloniex](https://poloniex.com). USDT/BTC and BTC/USD price history
is downloaded, and the weighted average price over 15 minute periods is used to find
the monthly average USDT_DCR price.

dcrwages is currently hardcoded to collect price data starting with the current
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

Once the data is obtained, the web server will start listening on port 3000:
<http://localhost:3000/>

## REST API

| HTTP Request            | Response                 |
|-------------------------|--------------------------|
| `/api/prices`               | All available prices     |
| `/api/prices/{year}`        | Prices for a single year |
| `/api/prices/{year}/{month}`| Price for a single month |

Months are represented by numbers. 1 = Jan, 2 = Feb...

Errors are indicated using HTTP status codes and an error description in the response body.
For example, a request for data which is unavailble will give a `400` HTTP status with the following body:

```json
{"error":"No data for year 1966"}
```
