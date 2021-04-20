FROM golang:1.15 AS build

RUN apt-get update
RUN apt-get upgrade
RUN apt-get install git make build-essential libc6 -y

WORKDIR /src

COPY . /src

#RUN go mod download

RUN cd src/; GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -extldflags "-static"' -installsuffix cgo -o faceit-users .

FROM alpine AS runtime

EXPOSE 8888

COPY src/ /app/

RUN chmod +x /app/faceit-users

CMD ["/app/faceit-users"]
