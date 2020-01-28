FROM golang:1.13-alpine

WORKDIR /go/src/albion-status-checker

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["albion-status-checker"]