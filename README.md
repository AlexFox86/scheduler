# Информация
**Scheduler** - это проект веб-сервера, который реализует функциональность планировщика задач — аналог TODO-листа.

В планировщике реализован следующий функционал:

* Простейшая авторизация пользователя
* Добавление задачи
* Получение списка задач
* Удаление задачи
* Получение параметров задачи
* Изменение параметров задачи
* Отметка о выполнении задачи


# Задания со звёздочкой
* Все задания со звёздочкой выполнены полностью


# Запуск
## Локальный запуск
1. Создайте переменные окружения, например:

    * TODO_PORT="7540"
    * TODO_DBFILE="/home/user/Dev/scheduler/db/scheduler.db"
    * TODO_PASSWORD="12345"    
    * SECRET="secret12345" 
    
3. Склонируйте этот репозиторий
4. Соберите проект:
```
go build -o ./app/scheduler
```
5. Запустите приложение:
```
./app/scheduler
````
6. Отредактируйте файл tests/settings.go:
```
var Port = 7540
var DBFile = "../db/scheduler.db"
var FullNextDate = false
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Zc-4uvudFVfEeaqnkDY2c7CeCXnBDa3Fw-i1uh8CgxY`
```
7. Запустите тесты:
```
go test ./tests --count=1`
```

## Сборка и запуск через Docker
1. Склонируйте этот репозиторий
2. Соберите Docker образ:
```
docker build -t scheduler .
```
3. Запустите с параметрами:
```
docker run -v $PWD/db:/data -p 7540:7540 scheduler
```