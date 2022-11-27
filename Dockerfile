FROM --platform=amd64 golang:1.19.3-alpine3.16

WORKDIR /app

COPY . .

RUN go get

CMD go run .