FROM golang:1.7-alpine

RUN apk add --update git
RUN go get github.com/honeybadger-io/honeybadger-go

COPY . /sensebox-sketches

WORKDIR /sensebox-sketches

ENV GOBIN=$GOPATH/bin

RUN go install

RUN apk del git && rm -rf /var/cache/apk/*

EXPOSE 3924

CMD ["/go/bin/sensebox-sketches"]