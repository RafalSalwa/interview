FROM golang:1.22-alpine3.19
WORKDIR /interview
ENV TZ=Europe/Warsaw

RUN apk add --no-cache graphviz icu-libs tzdata && go install github.com/cosmtrek/air@latest \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
    
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["air", "-c", "./build/.air.gateway.toml"]
EXPOSE 8021 8022



#FROM golang:1.21-alpine3.18 AS builder
#ARG APP_NAME
#WORKDIR /app
#
#RUN apk add --no-cache graphviz && go install github.com/cosmtrek/air@latest
#
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -o $APP_NAME ./cmd/$APP_NAME/main.go
#
#FROM golang:1.21-alpine3.18
#ARG APP_NAME 
#ARG APP_ENV
#
#ENV APP_PATH=/app/$APP_NAME
#
#WORKDIR /app
#COPY --from=builder /app/$APP_NAME /app/$APP_NAME
#COPY docker-entrypoint.sh /app/docker-entrypoint.sh
#RUN chmod +x /app/docker-entrypoint.sh
#
#COPY --from=builder /app/cmd/$APP_NAME/config/config.$APP_ENV.yaml /app/$APP_NAME.$APP_ENV.yaml
#
#ENTRYPOINT [ "/app/docker-entrypoint.sh" ]
#