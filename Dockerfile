# Build Geth in a stock Go builder container
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-ethereum
WORKDIR /go-ethereum
RUN make geth

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["geth"]

#------------
# PYTHON

RUN apk --update add python3 \
    python3-dev

RUN apk --update add py-pip
RUN pip install --upgrade pip

#------------

RUN mkdir -p /home/gf/data

#------------

COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/