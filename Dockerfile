FROM golang:1.19-alpine AS BUILD
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:latest
RUN apk --no-cache add tzdata
COPY --from=BUILD /app/main /app/
CMD ["/app/main"]