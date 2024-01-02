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