FROM golang:1.18 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/user-syncer ./cmd

FROM alpine:3.16
RUN apk add --update mysql-client

COPY ./tmp/docker/migrate /usr/local/bin/migrate
COPY ./tmp/docker/table ./migate/table

COPY --from=builder /usr/src/app/cmd /usr/local/bin/
CMD ["sh"]