# Revisor

[![Build Status](https://drone.dbeliakov.ru/api/badges/dbeliakov/revisor/status.svg)](https://drone.dbeliakov.ru/dbeliakov/revisor)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbeliakov/revisor)](https://goreportcard.com/report/github.com/dbeliakov/revisor)
[![Coverage Status](https://coveralls.io/repos/github/dbeliakov/revisor/badge.svg)](https://coveralls.io/github/dbeliakov/revisor)

Revisor - легкий и простой сервис для проведения код-ревью без необходимости использовать системы контроля версий, предназначенный для проверки кода на младших курсах института.

Демо доступно на [https://revisor.dbeliakov.ru](https://revisor.dbeliakov.ru).

Возможности:
* Загрузка файла с исходным кодом напрямую в web-интерфейсе
* Просмотр разницы между любыми двумя версиями файла
* Многоуровневые вложенные комментарии
* Добавление комментария к любой строке любой ревизии файла
* Поддержка языка разметки markdown в комментариях

## Установка

```
docker run -d --name revisor \
    -e SECRET_KEY=<some secret key> \
    -p 80:80 \
    -v /srv/revisor:/database \
    dbeliakov/revisor
```

Либо используя docker-compose:

```
version: '2'

services:
    revisor:
        image: dbeliakov/revisor:latest
        restart: always
        container_name: revisor
        environment:
            - SECRET_KEY: <some secret key>
        ports:
            - "80:80"
        volumes:
        - /srv/revisor:/database
```
