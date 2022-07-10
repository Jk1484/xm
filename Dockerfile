FROM golang:1.17-alpine as build

RUN mkdir /xm

ADD . /xm

WORKDIR /xm

RUN go build -o main ./cmd

FROM alpine:latest
COPY --from=build /xm /xm

WORKDIR /xm

CMD ["/xm/main"]