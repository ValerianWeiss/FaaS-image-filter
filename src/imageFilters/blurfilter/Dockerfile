FROM golang:1.9.7-alpine

ENTRYPOINT []

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

ADD https://github.com/openfaas/faas/releases/download/0.8.0/fwatchdog /usr/bin
RUN chmod +x /usr/bin/fwatchdog

WORKDIR /fsrc


ADD https://api.github.com/repos/ValerianWeiss/FaaS-image-filter/git/refs/heads/master version.json
RUN git clone https://github.com/ValerianWeiss/FaaS-image-filter.git

WORKDIR /go/src/blurfilter

RUN cp /fsrc/FaaS-image-filter/src/imageFilters/blurfilter/handler.go .
RUN cp -r /fsrc/FaaS-image-filter/src/imageFilters/utils /go/src/blurfilter/utils
RUN rm -r /fsrc
RUN go install


ENV fprocess "/go/bin/blurfilter"
HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD [ "/usr/bin/fwatchdog"]
