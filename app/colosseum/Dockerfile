### Builder
FROM golang:1.11-alpine as builder

RUN apk update && apk add git

WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .

ENV GO11MODULE on

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o colosseum app/colosseum/main.go 


### Runtime
FROM ubuntu:18.04

RUN useradd --create-home -s /bin/bash app
WORKDIR /home/app

RUN apt update
RUN apt install -y build-essential

# python
RUN apt install -y python3.6 python3.6-dev python3-pip
RUN ln -s $(which python3) /usr/bin/python
RUN ln -s $(which pip3) /usr/bin/pip

COPY gamer gamer
RUN pip install -r gamer/python3/requirements.txt

RUN chown -R app:app /home/app
USER app

COPY --from=builder /usr/src/app/colosseum /home/app/colosseum
RUN mkdir -p /home/app/source

ENV COLOSSEUM_BASE_PATH /home/app
ENV PYTHONPATH /home/app

CMD ["./colosseum"]
