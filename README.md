# urlshortener

Путь до конфиг-файла получаем из переменной окружения CONFIG_PATH, дефолтный путь не предусмотрен. Чтобы передать значение такой переменной, можно запустить приложение следующей командой:

```
export CONFIG_PATH=./config/local.yaml
```

```
CONFIG_PATH=./config/local.yaml ./urlshortener
```

```
CONFIG_PATH=./config/local.yaml go run ./cmd/url-shortener/
main.go
```

Если возникла оршибка: 
`level=ERROR msg="failed to initialize storage" env=local error="storage.sqlite.NewStorage: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub"`
Для работы sqlite3 необходим GCC
https://pkg.go.dev/github.com/mattn/go-sqlite3#section-readme
Установка GCC в Ubuntu
Если вас устраивает текущая версия GCC, которая есть в официальных репозиториях дистрибутива, то вам достаточно установить пакет build-essential. Для этого выполните команду:

```
sudo apt -y install build-essential
```