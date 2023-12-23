FROM arm64v8/postgres:latest

# Обновление списка пакетов
RUN apt-get update

# Установка необходимых пакетов
RUN apt-get install -y postgresql-postgis

# Удаление ненужных файлов
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*


ENV POSTGRES_USER usr
ENV POSTGRES_PASSWORD pass1234
ENV POSTGRES_DB location

COPY ./init-postgis.sql /docker-entrypoint-initdb.d/