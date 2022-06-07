FROM golang:1.18-buster as builder
WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o output

FROM scratch
COPY --from=builder /app/output/ /gcp-resource-manager
# COPY serviceaccount.json /serviceaccount.json
# ENV GOOGLE_APPLICATION_CREDENTIALS serviceaccount.json
# ENV PROJECT_ID $projectId

CMD ["/gcp-resource-manager"]