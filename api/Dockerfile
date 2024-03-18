FROM golang:1.21-alpine AS build

WORKDIR /app
COPY . /app
RUN go build -o backend

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/backend .
EXPOSE 8080
CMD ["./backend"]