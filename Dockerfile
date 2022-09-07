FROM golang:1.19-alpine

# Install migrate
RUN apk update && \
    apk add curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/bin/migrate && \
    go install github.com/swaggo/swag/cmd/swag@latest

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /usr/bin/wait
RUN chmod +x /usr/bin/wait

WORKDIR /app

COPY . .

RUN touch .env && \
    go build -o backend cmd/project_template/* && \
    chmod +x scripts/entrypoint.sh

ENTRYPOINT ["/app/scripts/entrypoint.sh"]