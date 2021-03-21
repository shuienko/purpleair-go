FROM golang:1.15

WORKDIR /go/src/purpleair-go
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["purpleair-go"]