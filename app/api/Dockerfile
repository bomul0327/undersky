### Builder
FROM golang:1.11-alpine as builder

RUN apk update && apk add git

WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .

ENV GO11MODULE on

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o api app/api/*.go


### Make executable image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/src/app/api /api

CMD [ "/api" ]
