FROM golang:1.13-alpine

COPY . /app
WORKDIR /app
RUN go build
ENTRYPOINT ["/app/hacktoberfest-leaderboard"]
EXPOSE 4000
