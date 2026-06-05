FROM golang:1.26 AS deps

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

# As build stage we reuse the go module cache
FROM golang:1.26 AS build

WORKDIR /src

COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/e2b-server ./cmd/server

# Run layer
FROM gcr.io/distroless/static-debian12:nonroot AS run

WORKDIR /app

ENV HOME=/home/nonroot
ENV PORT=8080

COPY --from=build /out/e2b-server /app/e2b-server

EXPOSE 8080

ENTRYPOINT ["/app/e2b-server"]
