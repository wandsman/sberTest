markdown
# Описание
Этот проект представляет собой пример многопоточного приложения на Go, которое выполняет HTTP-запросы к внешнему серверу.

## Компиляция
Для компиляции проекта вам потребуется Go версии 1.16 или выше. Вы можете скачать его с [официального сайта](https://golang.org/dl/).

Склонируйте репозиторий в вашу рабочую директорию:
```bash
git clone https://github.com/yourusername/yourrepository.git
```

Перейдите в директорию проекта:

```bash
cd yourrepository
```
Скомпилируйте проект:

```bash
go build -o app cmd/app/main.go
```

## Конфигурация
Количество горутин, которые будут запущены приложением, можно настроить с помощью флага `-threads`. По умолчанию, это значение равно 3.

## Запуск
Запустите скомпилированный бинарный файл:

```bash
./app -threads=5
```

В этом примере приложение запустит 5 горутин.

## Примечание
Если вы получаете ошибки при запуске более трех горутин