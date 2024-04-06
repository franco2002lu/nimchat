FROM golang:1.22 as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /out/main ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates && \
    adduser -D user

USER user

COPY --from=builder /out/main /app/main

EXPOSE 8080

CMD ["/app/main"]