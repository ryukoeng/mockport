FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /mockport ./cmd/mockport

FROM gcr.io/distroless/static-debian12

COPY --from=builder /mockport /mockport

EXPOSE 43101

ENTRYPOINT ["/mockport"]
CMD ["run", "--config", "/etc/mockport/mockport.yml"]
