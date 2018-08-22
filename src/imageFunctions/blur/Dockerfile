FROM golang:1.9.7-alpine

ENTRYPOINT []

WORKDIR /go/src/imagefilter/openfaas/faas/function
RUN mkdir ./utils

COPY ./handler.go .
COPY ./utils/utils.go ./utils

RUN go install

ADD https://github.com/openfaas/faas/releases/download/0.8.0/fwatchdog /usr/bin
RUN chmod +x /usr/bin/fwatchdog

ENV fprocess "/go/bin/golang"
HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["/usr/bin/fwatchdog"]