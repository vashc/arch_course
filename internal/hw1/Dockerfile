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

RUN go build -o $WD/hw1 arch_course/cmd/hw1

FROM alpine:3.13 AS worker

RUN apk --no-cache add ca-certificates

COPY --from=builder /tmp/bin/hw1 /app/hw1

EXPOSE 8080

CMD ["/app/hw1"]