FROM golang:1.17-alpine AS builder

RUN apk add --no-cache git

ENV CGO_ENABLED=0
ENV GO111MODULE on
ENV GOSUMDB off
ENV WD=/tmp/bin

WORKDIR $WD

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o $WD/hw5/auth-service arch_course/cmd/hw5/auth-service

FROM alpine:3.13 AS worker

RUN apk --no-cache add ca-certificates

COPY --from=builder /tmp/bin/hw5/auth-service /app/hw5/auth-service

EXPOSE 8000

CMD ["/app/hw5/auth-service"]