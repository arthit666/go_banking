FROM golang:1.21-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o ./out/go-app .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-app /app/go-app

EXPOSE 8000

CMD ["/app/go-app"]