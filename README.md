# masker

Сервис читает строки из файла (Producer), маскирует `http://...` ссылки, и записывает результат (Presenter).
Бизнес-логика изолирована в `internal/service` и покрыта unit-тестами через `testify/mock`.

## Сборка и запуск

```bash
go mod tidy
go build -o masker ./cmd/app

# usage: ./masker <input_path> [output_path]
printf "see http://abc page\nno links\nhttp://a http://bb\n" > input.txt
./masker input.txt result.txt
cat result.txt
```
