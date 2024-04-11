# sql

Занятие 5-2 SQL

Создание контейнера с Postgress
```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -d postgres
docker exec -it postgres psql -U test -W
```

## psql — PostgreSQL interactive terminal

### \q - выход
### \l - список баз данных
### \c <database> - переключить базу данных
### \dn - список схем
