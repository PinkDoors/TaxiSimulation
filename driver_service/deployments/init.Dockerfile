# Используйте базовый образ Go
FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/app ./app
COPY --from=build:develop /app/.env.development ./.env.development

ENV APP_ENV=development

CMD ["/app/app"]
