FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /yeti ./cmd/yeti

FROM alpine:3.21
COPY --from=build /yeti /usr/local/bin/yeti
EXPOSE 8080
ENTRYPOINT ["yeti"]
