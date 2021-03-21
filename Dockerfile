FROM golang:1.15
LABEL maintainer="oleksandr.shuienko@gmail.com"

WORKDIR /go/src/purpleair-go
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["purpleair-go"]