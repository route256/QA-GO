Требования к сервису act-device-api

# Назначение

Сервис предназначен для регистрации мобильных устройств и отправки уведомлений по ним.
Сервис должен поддерживать работу по протоколам: gRpc и Rest api.

# 1. Описание протокола act_device_api

## 1.1 Создание устройства

Сервис должен обеспечивать возможность создания устройства.

### gRPC

Запрос на создание устройства (CreateDeviceV1Request) должен включать следующие поля:

| Поле     | Тип данных | Описание              | Валидация                                 |
|----------|------------|-----------------------|-------------------------------------------|
| platform | String     | Платформа ios/android | Мин 1 символ.<br/>Может быть любая строка |
| user_id  | Uint64     | Id пользователя.      | Не нулевое                                |

Запрос должен возвращать CreateDeviceV1Response со следующими полями:

| Поле      | Тип данных | Описание                 |
|-----------|------------|--------------------------|
| device_id | Uint64     | Уникальное ID устройства |

### Rest

Rest протокол должен соответствовать gRPC
и обрабатываться посредством POST запросы по пути /api/v1/devices:

| Поле     | Тип данных     | Описание              | Валидация                                 |
|----------|----------------|-----------------------|-------------------------------------------|
| platform | String         | Платформа ios/android | Мин 1 символ.<br/>Может быть любая строка |
| user_id  | String(Uint64) | Id пользователя.      | Не нулевое                                |

### БД

При создании устройство должно сохраняться в БД в табличном виде:

| Поле       | Тип данных  | Описание                                   |
|------------|-------------|--------------------------------------------|
| Id         | integer     | Уникальное ID устройства                   |
| platform   | Varchar(32) | Имя платформы                              |
| removed    | boolean     | Признак удаления устройства (см пункт 1.4) |
| user_id    | Bigint      | Id пользователя                            |
| entered_at | Timestamp   | Timestamp создания устройства в приложении |
| created_at | Timestamp   | Timestamp создания устройства в БД         |
| updated_at | Timestamp   | Timestamp - последнего обновления          |

### Логика

При создании устройства - должна произойти запись о созданном устройстве в таблицу Devices.

## 1.2 Получение информации об устройстве

### gRPC

Запрос информации об устройстве DescribeDeviceV1Request должен содержать:

| Поле      | Тип данных | Описание                 |
|-----------|------------|--------------------------|
| device_id | Uint64     | Уникальное ID устройства |

На запрос должен возвращаться ответ DescribeDeviceV1Response со следующими полями:

| Поле       | Тип данных | Описание                                   |
|------------|------------|--------------------------------------------|
| Id         | uint64     | Уникальное ID устройства                   |
| platform   | String     | Имя платформы                              |
| user_id    | uint64     | Id пользователя                            |
| entered_at | Timestamp  | Timestamp создания устройства в приложении |

### REST

Соответствующий запрос в HTTP протоколе должен происходить с помощью GET /api/v1/devices/{device_id},
где device_id - ID - существующего устройства.

Сущность ошибки представлена следующими полями:

| Поле    | Тип данных | Описание             |
|---------|------------|----------------------|
| code    | int4       | Код ошибки           |
| message | String     | Текст ошибки         |
| details | Array      | Доп детали по ошибке |

### БД

Отсутствует

### Логика

При поступлении запроса на получение устройства - соответствующее актуальное устройство - должно быть возвращено в
ответе.

Если ID не существует или устройство было удалено - должна вернуться 404 ошибка
со следующим содержимым:

| Поле    | Значение         |
|---------|------------------|
| code    | 5                |
| message | device not found |

## 1.3 Редактирование устройства

Устройство может быть отредактировано с помощью запроса: UpdateDeviceV1Request, на который должен прийти ответ
UpdateDeviceV1Response.

### gRPC

UpdateDeviceV1Request

| Поле      | Тип данных | Описание                                         |
|-----------|------------|--------------------------------------------------|
| device_id | uint64     | Ненулевое Уникальное ID устройства               |
| platform  | String     | Имя платформы, Не больше 32 символов, Не пустое. |
| user_id   | uint64     | Ненулевое ID пользователя во внешней системе     |

UpdateDeviceV1Response

| Поле    | Тип данных | Описание                                                        |
|---------|------------|-----------------------------------------------------------------|
| success | Boolean    | Результат выполнения операции True - успешно, False - неуспешно |

### Rest

При отправке PUT запроса по пути `/api/v1/devices/{device_id}`, где device_id - id существующего устройства в системе.

Тело запроса JSON с полями:

| Поле     | Тип данных | Описание                           |
|----------|------------|------------------------------------|
| platform | String     | Платформа - Длина 32 символа       |
| userId   | uint64     | ID пользователя во внешней системе |


### Логика
При поступлении запроса на обновление, существующего устройства, запись об устройстве, должна быть обновлена в таблице Devices. Поле Updated_At должно быть обновлено на timestamp времени
обновления.
При поступлении запроса на обновление, не существующего или удаленного устройства - должна быть возвращен ответ UpdateDeviceV1Response c полем success =
False.

## 1.4 Удаление устройства

Система должна поддерживать функционал удаление устройства.

### gRPC

RemoveDeviceV1Request

| Поле      | Тип данных | Описание                           |
|-----------|------------|------------------------------------|
| device_id | uint64     | Ненулевое Уникальное ID устройства |

RemoveDeviceV1Response

| Поле  | Тип данных | Описание                   |
|-------|------------|----------------------------|
| found | bool       | Статус удаления устройства |

### Rest

Для REST протокола, удаление должно происходить при получении запроса:
`DELETE /api/v1/devices/{device_id}`, где device_id - id существующего устройства.

### БД

Обновить в таблице Devices следующие поля:

| Поле    | Тип данных | Описание                    |
|---------|------------|-----------------------------|
| removed | bool       | Признак удаления устройства |

### Логика:

При поступлении запроса с существующим device_id, устройство с таким идентификатором помечается как удаленное.
Если устройство уже помечено, как удаленное - возвращается ответ RemoveDeviceV1Response со статусом found = false.

При получении запроса на удаление устройства - необходимо обновить в таблице Devices:
Поле removed в соответствующем устройстве - проставляется как True.
Удаленные устройства не отображаются в списке устройств.
Удаленное устройство невозможно запросить по ID.
Удаленное не может быть отредактировано.

## 1.5 Список устройств

При запросе списка доступных устройств - должна вернуться соответствующая таблица, отсортированного списка устройства.

### GRPC

ListDevicesV1Request

| Поле     | Тип данных | Описание                               |
|----------|------------|----------------------------------------|
| page     | uint64     | Страница списка устройств              |
| per_page | uint64     | Количество устройств на странице (>=1) |

ListDevicesV1Response

| Поле  | Тип данных | Описание            |
|-------|------------|---------------------|
| items | []Device   | Коллекция устройств |

Device

| Поле       | Тип данных | Описание                             |
|------------|------------|--------------------------------------|
| id         | uint64     | ID устройства                        |
| platform   | string     | Платформа ios/android                |
| user_id    | uint64     | UserID устройства во внешней системе |
| entered_at | Timestamp  | Дата-время создания устройства       |


### Rest

Соответствующй маппинг должен быть проброшен на запрос GET /api/v1/devices
с параметрами:

| Поле     | Тип данных | Описание                                              |
|----------|------------|-------------------------------------------------------|
| page     | string     | (Обязательный) Страница списка устройств              |
| per_page | string     | (Обязательный) Количество устройств на странице (>=1) |

### БД

Отсутствует

### Логика

При поступлении запроса - должен быть вычислена страница, исходя из переданного размера и должна быть возвращена в
Ответ.
