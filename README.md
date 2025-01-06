# dcrwages

[![Build Status](https://github.com/jholdstock/dcrwages/workflows/Build%20and%20Test/badge.svg)](https://github.com/jholdstock/dcrwages/actions)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

dcrwages is a web application which provides average monthly USDT/DCR price.
It starts a web server with a HTML page displaying the price information, as
well as providing the information over an API.

<http://dcrwages.jholdstock.uk>

This rate is used when Decred project contractors submit invoices denominated in
US Dollars and receive payment in DCR.

## How dcrwages works

### Before April 2019

Price data retrieved from [Poloniex](https://poloniex.com) was used to calculate
the monthly price of Decred in USDT.
The weighted average price over 15 minute intervals on the USDT/BTC and BTC/DCR
markets were used to find monthly average USDT/DCR prices.

This historic data is now hard-coded in dcrwages - the Poloniex API is no longer
used by this project.

### April 2019 - November 2024

The price calculation was updated to use [Binance](https://binance.com) instead
of Poloniex, and the interval was changed from 15 minutes to 1 hour.

### Since December 2024

The [Binance](https://binance.com) DCR/USDT market is now used to find the
monthly average price instead of the DCR/BTC and BTC/USDT markets.

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
docker run -d -p 3000:3000 jholdstock/dcrwages:latest
```

The process will begin contacting Binance and downloading price information.

The web server will start listening on port 3000.
You can open the homepage in your browser <http://localhost:3000/>.

## REST API

| HTTP Request                | Response                 |
|-----------------------------|--------------------------|
| `/api/prices`               | All available prices     |
| `/api/prices/{year}`        | Prices for a single year |
| `/api/prices/{year}/{month}`| Price for a single month |

Months are handled as integers. 1 = Jan, 2 = Feb, etc.

Errors are indicated using HTTP status codes and an error description in the
response body.
For example, a request for data which is unavailble will give a `404` HTTP
status and the following body:

```json
{ "error": "No data for year 1966" }
```
