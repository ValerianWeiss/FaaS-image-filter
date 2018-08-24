FROM golang:1.9.7-alpine

ENTRYPOINT []

WORKDIR /fsrc
RUN git clone https://github.com/ValerianWeiss/FaaS-image-filter.git
RUN cp FaaS-image-filter/src/imageFilters/blurfilter/blurfilter.go /go/src/blurfilter
RUN cp FaaS-image-filter/src/imageFilters/utils /go/src/blurfilter

WORKDIR /go/src/blurfilter
RUN rm -f /fsrc
RUN go install

ADD https://github.com/openfaas/faas/releases/download/0.8.0/fwatchdog /usr/bin
RUN chmod +x /usr/bin/fwatchdog

ENV fprocess "/go/bin/golang"
HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD [ "/usr/bin/fwatchdog"]