FROM node:22-alpine AS css
WORKDIR /src
COPY static/css/input.css static/css/
COPY . .
RUN npx -y @tailwindcss/cli -i static/css/input.css -o static/css/tailwind.css --minify

FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=css /src/static/css/tailwind.css static/css/tailwind.css
RUN go build -o /yeti ./cmd/yeti

FROM alpine:3.21
COPY --from=build /yeti /usr/local/bin/yeti
EXPOSE 8080
ENTRYPOINT ["yeti"]
