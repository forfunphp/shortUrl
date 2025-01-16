FROM golang:1.20 AS builder

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/metal_shop

RUN go mod tidy && go build -o app

# Use a smaller base image for the final image
FROM golang:1.20

RUN apt-get update && apt-get install -y webp

WORKDIR /app

COPY --from=builder /usr/local/bin/dockerize /usr/local/bin/dockerize

COPY --from=builder /app/cmd/metal_shop/app .

COPY .env .

RUN mkdir /app/temp

EXPOSE 8080

CMD ["dockerize", "-wait", "tcp://postgres:5432", "-timeout", "30s", "./app"]
