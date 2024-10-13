GoNews Scraper Service
======================
Агрегатор новостей по RSS фидам.

# Конфигурационные параметры

| Параметр         | Описание                                                       | Значение по-умолчанию                                                                                                                               |
|------------------|----------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| `http_port`      | порт Scraper сервиса                                           | `8081`                                                                                                                                              |
| `request_period` | интервал опроса RSS фидов (т.е. как часто ходить за новостями) | `10m0s`                                                                                                                                             |
| `db_conn_string` | строка подключения к СУБД Postgres                             | `postgres://postgres@localhost:5432/news?sslmode=disable`                                                                                           |
| `rss_feeds`      | список URL с RSS фидами                                        | `["https://habr.com/ru/rss/hub/go/all/?fl=ru", "https://habr.com/ru/rss/best/daily/?fl=ru", "https://cprss.s3.amazonaws.com/golangweekly.com.xml"]` |

# Сборка и запуск

## Требования

-   docker >=23.0.0
-   golang 1.22


Приложение поддерживает конфиг-файлы в формате yaml:

    $ make build
    $ ./bin/go-news-scraper -print-config > config.yaml
    $ ./bin/news-server -config ./config.yaml

---

Для быстрого запуска с дефолтным конфигом:

    $ make run

Логи будут писаться сюда:

    $ tail -f log/go-news-scraper.log

Остановить приложение и удалить контейнер с базой:

    $ make clean

Показать версию сборки:

    $ ./bin/go-news-scraper -version

# Тесты
 
Полный прогон имеющихся тестов:

    $ make test

Для удаления тестовой базы:

    $ make clean
