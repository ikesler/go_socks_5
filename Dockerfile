FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /go_socks_5

EXPOSE 10800
CMD [ "/go_socks_5" ]
