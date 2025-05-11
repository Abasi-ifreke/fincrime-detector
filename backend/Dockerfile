FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


FROM golang:1.24-alpine 

WORKDIR /app

COPY --from=builder /app . 

COPY .env .

RUN mkdir -p /var/log/app

EXPOSE 8080

CMD ["go", "run", "."]