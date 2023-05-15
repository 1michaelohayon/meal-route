FROM golang:1.20.3-alpine


workdir /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o run_api .

EXPOSE 5000


CMD ["./run_api"]