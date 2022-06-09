FROM golang:1.18-buster as builder
WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o output

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=builder /app/output/ /gcp-resource-manager
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/gcp-resource-manager"]