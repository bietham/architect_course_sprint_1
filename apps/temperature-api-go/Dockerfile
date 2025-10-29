# ---------- Build stage ----------
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/temperature-api

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /out/temperature-api /temperature-api
EXPOSE 8081
ENV PORT=8081
USER 65532:65532
ENTRYPOINT ["/temperature-api"]
