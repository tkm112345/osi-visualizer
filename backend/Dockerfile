# --- build stage ---
FROM golang:1.24-alpine AS build
WORKDIR /src

# 依存は標準ライブラリのみ（go.sum は無い）。まず go.mod だけコピーしてキャッシュを効かせる。
COPY go.mod ./
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

# --- run stage ---
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
