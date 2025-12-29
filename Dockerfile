FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN --mount=type=bind,target=. CGO_ENABLED=0 go build -ldflags="-s -w" -o /server .

FROM scratch

COPY --from=builder /server /server

EXPOSE 8421

ENTRYPOINT ["/server"]
