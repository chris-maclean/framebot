FROM golang:1.17-alpine

WORKDIR /usr/src/app

RUN apk update && apk add bash
RUN apk add --no-cache ffmpeg


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

# Run framebot once per minute
RUN echo "*/5 * * * * /usr/local/bin/framebot" >> /etc/crontabs/root

CMD ["crond", "-f"]

# Start docker container:
# docker run -v d:/git/framebot/dev/start-job.json:/opt/framebot/framebot-state.json -v d:/git/framebot/framebot.cron:/etc/cron.d/framebot -ti -v e:/My\ Videos/Movies/2001\ A\ Space\ Odyssey.mkv:/opt/framebot/movie --env-file ./dev/.env --rm cmaclean/framebot