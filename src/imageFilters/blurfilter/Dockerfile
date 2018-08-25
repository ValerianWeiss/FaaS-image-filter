FROM golang:1.9.7-alpine

ENTRYPOINT []

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /fsrc

ADD https://api.github.com/repos/ValerianWeiss/FaaS-image-filter/git/refs/heads/master version.json
RUN git clone -b master https://github.com/ValerianWeiss/FaaS-image-filter.git

WORKDIR /go/src/FaaS-image-filter/src/imageFilters/
RUN cp -r /fsrc/FaaS-image-filter/src/imageFilters/blurfilter ./blurfilter
RUN cp -r /fsrc/FaaS-image-filter/src/imageFilters/utils ./utils
RUN rm -r /fsrc

WORKDIR /go/src/FaaS-image-filter/src/imageFilters/blurfilter
RUN go install


FROM alpine

ADD https://github.com/openfaas/faas/releases/download/0.8.0/fwatchdog /usr/bin
RUN chmod +x /usr/bin/fwatchdog

COPY --from=0 /go/bin /go/bin

ENV fprocess "/go/bin/blurfilter"
HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["/usr/bin/fwatchdog"]
