FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/app ./app
COPY --from=build:develop /app/.env.dev ./.env.dev

ENV APP_ENV=dev

CMD ["/app/app"]
