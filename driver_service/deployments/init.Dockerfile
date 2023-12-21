# Используйте базовый образ Go
FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/app ./app
COPY --from=build:develop /app/.env ./.env
#
#ENV APP_ENV=Development

CMD ["/app/app"]
