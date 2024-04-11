FROM golang:1.22-alpine3.19
WORKDIR /interview
ENV TZ=Europe/Warsaw

RUN apk add --no-cache graphviz icu-libs tzdata && go install github.com/cosmtrek/air@latest \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env /interview/bin/

CMD ["air", "-c", "./build/.air.consumer_service.toml"]
EXPOSE 8021 8022
