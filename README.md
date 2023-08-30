![GitHub repo size](https://img.shields.io/github/repo-size/Morozhkaa/User-segmentation-service)
![Repository Top Language](https://img.shields.io/github/languages/top/Morozhkaa/User-segmentation-service)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/Morozhkaa/User-segmentation-service)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/Morozhkaa/User-segmentation-service)

# User segmentation service

<img align="right" width="50%" src="./images/Customer-Segmentation.png">

Cервис, хранящий пользователя и сегменты, в которых он состоит (создание, удаление сегментов, а также добавление и удаление пользователей в сегмент).

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- gin (веб фреймворк)
- pgx (драйвер для работы с PostgreSQL)
- slog для логгирования

Сервис был написан с Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его. Также был реализован Graceful Shutdown для корректного завершения работы сервиса.


# Usage

Запустить сервис можно с помощью команды `make compose-up`.

Документацию после запуска сервиса можно посмотреть по адресу `http://localhost:3000/swagger/index.html`.

Для запуска линтера необходимо выполнить команду `make linter`.


## Parametrs formats
  * slug  `^[\w-]+$` - короткое имя, содержащее только буквы, цифры, символы подчеркивания, или дефисы.
  * userID  `uuid` - идентификатор пользователя.
  * period  `^\d{4}-\d{2}$` - месяц, за который вы хотите отобразить.


## Examples

Некоторые примеры запросов
- [Создание сегмента](#create)
- [Удаление сегмента](#delete)
- [Обновление информации о сегментах у пользователя](#update)
- [Получение всех сегментов пользователя](#getSegments)
- [История событий за заданный месяц в формате csv файла](#report)
- [История событий за заданный месяц для конкретного пользователя в формате csv файла](#userreport)


### Создание сегмента <a name="create"></a>

```curl
curl -X 'POST' \
  'http://localhost:3000/api/v1/createSegment' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "slug": "AVITO_VOICE_MESSAGES"
}'
```
Пример ответа:
```json
{
  "success": "segment with slug 'AVITO_VOICE_MESSAGES' created"
}
```

### Удаление сегмента <a name="delete"></a>

```curl
curl -X 'DELETE' \
  'http://localhost:3000/api/v1/deleteSegment' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "slug": "AVITO_VOICE_MESSAGES"
}'
```
Пример ответа:
```json
{
  "success": "segment with slug 'AVITO_VOICE_MESSAGES' deleted"
}
```


### Обновление информации о сегментах у пользователя <a name="update"></a>

```curl
curl -X 'POST' \
  'http://localhost:3000/api/v1/updateUserSegments/550e8400-e29b-41d4-a716-446655440000' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segments-to-add": [
    "AVITO_DISCOUNT_50", "AVITO_PERFORMANCE_VAS"
  ],
  "segments-to-remove": [
    "AVITO_DISCOUNT_30"
  ]
}'
```
Пример ответа:
```json
{
  "success": "segment information for user with userID = 550e8400-e29b-41d4-a716-446655440000 updated"
}
```


### Получение всех сегментов пользователя <a name="getSegments"></a>

```curl
curl -X 'GET' \
  'http://localhost:3000/api/v1/getUserSegments/f15a3c54-90bc-4a9c-8070-e9e2394d872e' \
  -H 'accept: application/json'
```
Пример ответа:
```json
{
  "segments": [
    "AVITO_VOICE_MESSAGES",
    "AVITO_DISCOUNT_50"
  ]
}
```


### История событий за заданный месяц в формате csv файла <a name="report"></a>

```curl
curl -X 'GET' \
  'http://localhost:3000/api/v1/getReport/2023-08' \
  -H 'accept: application/json'
```
Пример ответа - csv файл с содержимым: 


```text/csv 
550e8400-e29b-41d4-a716-446655440000,AVITO_VOICE_MESSAGES,add,2023-08-30 14:38:42
550e8400-e29b-41d4-a716-446655440000,AVITO_DISCOUNT_50,remove,2023-08-30 14:44:57
```


### История событий за заданный месяц для конкретного пользователя в формате csv файла <a name="userreport"></a>

```curl
curl -X 'GET' \
  'http://localhost:3000/api/v1/getUserReport/2023-08/da3626c2-4747-11ee-be56-0242ac120002' \
  -H 'accept: text/csv'
```
Пример ответа - csv файл с содержимым: 


```text/csv 
da3626c2-4747-11ee-be56-0242ac120002,AVITO_VOICE_MESSAGES,add,2023-08-30 14:38:42
da3626c2-4747-11ee-be56-0242ac120002,AVITO_DISCOUNT_50,remove,2023-08-30 14:44:57
```

# Decisions <a name="decisions"></a>

1. При создании индентификатора пользователя использовать uuid или обычный auto-increment(indentity)?

> Решила использовать uuid, поскольку это бы избавило от проблемы создания одинакового идентификатора на разных рабочих серверах, если их в дальнейшем будет несколько.

2. Как реализовать хранение данных?

> Сначала планировала создать отдельные таблицы для сегментов и пользователей, а также общую segments_users таблицу для реализации связи многие-ко-многим. Но, поскольку у нас нет дополнительной информации о пользователях, пока не стала выделять на них отдельную таблицу. Вместо этого храню пары (сегмент - пользователь) в segments_users таблице.

3. Как вернуть отчёт о событиях?

> Подумала, что будет удобно автоматически сохранять файл с отчетом. При этом можно получить как общую информацию о событиях за месяц, так и отчет только по интересующему пользователю.


### Сomments

Вообще, сервис получился довольно простым, в основном фокус был на следовании принципам чистой архитектуры и написании автоматической документации. В дальнейшем, если будет большая нагрузка на базу, можно думать в сторону использования индексов, добавления кеширования, репликации или шардирования (кроме кеширования ни с чем не работала). Возникли трудности с написанием тестов. Зато осознала слабые места.