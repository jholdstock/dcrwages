# builder image
FROM golang:1.20

LABEL description="dcrwages"
LABEL version="1.0"
LABEL maintainer "holdstockjamie@gmail.com"

USER root
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build .

ENV GIN_MODE release
EXPOSE 3000
CMD ["/app/dcrwages"]