FROM golang:1.22 AS builder
MAINTAINER ZeroHash <contact@zerohash.com>

#filesys setup
RUN mkdir /app
ADD . /app
WORKDIR /app

#build and install the binary and deps
RUN make install

FROM alpine:3.19.1 AS runner

ARG target

RUN echo $target

COPY --from=builder /app/$target/ .


ENTRYPOINT ["./vwap-engine"]
