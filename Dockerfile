FROM golang:1.18 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -v -o ./user-syncer ./cmd

FROM alpine:3.16
RUN apk add --update mysql-client

COPY ./tmp/docker/migrate /usr/local/bin/migrate
COPY ./tmp/docker/table ./migrate/table

COPY --from=builder /usr/src/app/user-syncer /usr/local/bin/
CMD ["sh"]