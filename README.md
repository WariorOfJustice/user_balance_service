# Микросервис для работы с балансом пользователей

### Запуск
Перед запуском сервиса нужно создать базу (pg) и пользователя c правом входа в базу по данным из .env .\
Его содержание можно поменять, но имя базы и пользователя должны быть аналогичны информации в файле .env .\
Затем нужно накатить скрипт из init_sql.sql . После этого можно запускать сервис. \
Пробовал через docker-compose, но возникает ошибка при подключении к базе, которую не успел пофиксить. Ошибка может быть локальной, т к файлы docker и docker-compose вроде в порядке.


### Выполненные задачи:
* Метод получения баланса пользователя.
* Метод начисления средств на баланс.
* Метод резервирования средств с основного баланса на отдельном счете.
* Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии.\
У меня данный метод принимает id_reservation,
а не ИД пользователя, ИД услуги, ИД заказа, сумму.\
В чем суть и почему я так сделал. Метод резервирования у меня возвращает в json-е id_reservation.
Сервис, который запросил резервирование, его получает, сохраняет к себе. Таким образом, у нас есть четкая взаимосвязь таблицы с заказом из другого сервиса
и строкой резервирования. Затем, если нужно признать выручку/разрезервировать, сервис может отослать по этой этому заказу просто id резервирования.
Сервису не нужно отслеживать связку услуга-заказ-пользователь (сумму сюда не включил, поскольку она не является чем-то уникальным и важным для связки) + т к id_reservation это PK, то скорость поиска строки будет выше.

### Примеры запросов/ответов
Кейсы выполнял в цифровом порядке, зафиксированном в примерах

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/74541d525400d06fd039?action=collection%2Fimport)
