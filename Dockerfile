FROM golang:1.21-alpine AS builder

# install git required for fetching go dependencies
# virtual package mechanism to keep it small
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

COPY --from=builder /app/main /main

EXPOSE 8001

CMD [ "./main" ]
