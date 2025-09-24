FROM golang:alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.app.toml"]