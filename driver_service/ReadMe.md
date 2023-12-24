### Как запустить докер-контейнер с данным приложением и всей инфрой

docker-compose --env-file .env.development up -d

### Mongo
* Для просмотра БД можно использовать mongo-express. Креды: "user" и "pass".  
* Или можно использовать DataGrip. Креды для SCRAM-SHA-1: "root" и "root", база аутентификации - "admin".

### Kafka
Для создания топиков использовал Kafdrop, топики называются "trip-inbound" и "trip-outbound". Там же смотрел все сообщения.

### Location-Service
Мой коллега запустил свой сервис на удаленном хосте и я подключаюсь к нему через http://85.193.91.23:8081. К сожалению, 
если вы будете запускать локально мой сервис, то функционал get-trips не будет работать. Сможем показать работспособность 
только на защите.

**Как тестировал?**  
Использовал kcat (https://github.com/edenhill/kcat) для создания сообщений, чтобы их читал консьюмер.
* kcat -P -b localhost:29092 -t trip-outbound -D "\e"  
* Пишем тело эвента, ниже приведу примеры для создания трипа и изменения его статуса:  
```
{
    "id": "284655d6-0190-49e7-34e9-9b4060acc123",
    "source": "/trip",
    "type": "trip.event.created",
    "datacontenttype": "application/json",
    "time": "2023-11-09T17:31:00Z",
    "data": {
        "trip_id": "e82c42d6-b86f-4e2a-93a2-858413acb123",
        "offer_id": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImZyb20iOnsibGF0IjowLCJsbmciOjB9LCJ0byI6eyJsYXQiOjAsImxuZyI6MH0sImNsaWVudF9pZCI6InN0cmluZyIsInByaWNlIjp7ImFtb3VudCI6OTkuOTUsImN1cnJlbmN5IjoiUlVCIn19.fg0Bv2ONjT4r8OgFqJ2tpv67ar7pUih2LhDRCRhWW3c",
        "price": {
        "currency": "EUR",
        "amount": 200
        },
        "status": "DRIVER_SEARCH",
        "from": {
            "lat": 13.123,
            "lng": 13.123
        },
        "to": {
            "lat": 13.123,
            "lng": 13.123
        }
    }
}
  ```
* Control+D на маке (На винде не знаю, сори)

### Полезные ссылки
MongoExpress: http://localhost:8081/
Хост Mongo для подключения в DataGrip (или еще где-нибудь): http://localhost:27017/
Kafdrop: http://localhost:9080/
Prometheus: http://localhost:9090/  
Jaeger: http://localhost:16686/  
Само приложение: http://localhost:8080/  