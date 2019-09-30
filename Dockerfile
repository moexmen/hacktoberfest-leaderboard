FROM golang:1.13-alpine as builder

COPY . /app
WORKDIR /app
RUN go build

FROM alpine
RUN apk add tzdata
COPY --from=builder /app /app
WORKDIR /app
ENTRYPOINT ["/app/hacktoberfest-leaderboard"]
EXPOSE 4000
