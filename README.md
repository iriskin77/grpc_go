## Как работать с proto файлами

1) Установить компилятор для прото файлов 

2) Написать прото файлы

3) Использовать команду: 

protoc -I proto proto/sso/sso.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

protoc -I proto proto/sso/*.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

Где:

+ **-I proto** - это папка, где находятся все прото файлы;
+ **proto/sso/sso.proto** - путь до конкретного прото файла, который необходимо скомпилировать;
+ **--go_out=./gen/go** - путь, куда будут сохраняться скомпилированные файлы;
+ **--go-grpc_out=./gen/go/** - путь, куда будет сохраняться go-grpc код;

## Запуск приложения:

+ go run cmd/sso/main.go --config=./config/local.yaml

## Создание миграций

+ go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations