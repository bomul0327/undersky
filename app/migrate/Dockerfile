### Builder
FROM golang:1.11-alpine as builder

RUN apk update && apk add git

WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .

ENV GO11MODULE on

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o migrate app/migrate/*.go


### Make executable image
FROM scratch

COPY --from=builder /usr/src/app/migrate /migrate

CMD [ "/migrate", "up" ]
