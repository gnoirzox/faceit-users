FROM golang:1.15-alpine AS build

RUN apk add --no-cache git

WORKDIR /src

COPY . /src

#RUN go mod download

RUN cd src/; go build -o faceit-users

FROM alpine AS runtime

COPY src/faceit-users /app/faceit-users

CMD ["/app/faceit-users"]
