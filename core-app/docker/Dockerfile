FROM golang:1.22.0-alpine as builder

WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o zex-app .

FROM alpine as runner

WORKDIR /app
COPY --from=builder /app/zex-app /app
COPY ./frontend /app/frontend
COPY ./sqlc /app/sqlc
COPY ./static /app/static

ENTRYPOINT ["/app/zex-app"]