# Building the binary of the App
FROM golang:1.17.2-alpine3.14 AS build

WORKDIR /go/src/micro-order

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest

WORKDIR /app

COPY --from=build /go/src/micro-order/app .

## THE LIFE SAVER
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait
RUN chmod +x /wait

# Exposes port 3000 because our program listens on that port
EXPOSE 3006

CMD /wait && ./app