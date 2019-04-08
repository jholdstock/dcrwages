# dcrwages

[![Build Status](https://travis-ci.org/jholdstock/dcrwages.png?branch=master)](https://travis-ci.org/jholdstock/dcrwages)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

dcrwages is a web application which provides average monthly USDT/DCR price.
It starts a web server with a HTML page displaying the price information,
as well as providing the information over a RESTful API.

## How dcrwages works

dcrwages calculates the average monthly price of Decred in USDT using price data
retrieved from [Poloniex](https://poloniex.com). The weighted average price over
15 minute intervals on the USDT/BTC and BTC/USD markets are used to find
the monthly average USDT_DCR price.

dcrwages collects price data starting with the current
month and working back to June 2016.
Data for DCR/BTC is not available on Poloniex before this time.

dcrwages was written with Go 1.12.

## How to use dcrwages

Build and run locally:

```bash
# Build the executable
env GO111MODULE=on go install .

# Run the executable
$(go env GOPATH)/bin/dcrwages
```

Or build and run in docker:

```bash
# Build the container
docker build -t jholdstock/dcrwages .

# Run the container
docker run -d -p 3000:3000 jholdstock/dcrwages
```

The process will begin contacting Poloniex and downloading price information.
Poloniex only allows 6 HTTP requests on it's API per second, so collecting
all of the data takes around 30 seconds.

The web server will start listening on port 3000. You can open the homepage
in your browser
<http://localhost:3000/>

## REST API

| HTTP Request                | Response                 |
|-----------------------------|--------------------------|
| `/api/prices`               | All available prices     |
| `/api/prices/{year}`        | Prices for a single year |
| `/api/prices/{year}/{month}`| Price for a single month |

Months are handled as integers. 1 = Jan, 2 = Feb, etc.

Errors are indicated using HTTP status codes and an error description in the response body.
For example, a request for data which is unavailble will give a `404` HTTP status and the following body:

```json
{"error":"No data for year 1966"}
```
