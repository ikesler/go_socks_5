FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY go-socks5/go.mod ./go-socks5/
RUN go mod download
COPY . .
RUN go build -o /go_socks_5

EXPOSE 10800
CMD [ "/go_socks_5" ]
