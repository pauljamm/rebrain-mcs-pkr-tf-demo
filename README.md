# Материалы к вебинару по Packer и Terraform в VK Cloud Solutions с Rebrain

## Подготовка

1. Регистрируемся в Mail.ru Cloud Solutions
2. Заходим в личный кабинет MCS, переходим в раздел "Настройки проекта"
и на вкладке "API ключи" нажимаем "Скачать openrc версии 3"
3. После скачивания файла нужно загрузить из него переменные командой

```bash
source <путь к скачанному файлу>
```

При вызове этой команды скрипт поросит ввести пароль от своей учетной записи в MCS

## Подготовка очереди сообщений для приложения

1. Переходим в интерфейс облака на вкладку "Очереди сообщений" -> "Ключи доступа"
1. Нажимаем добавить новый ключ, вводим любое имя
1. Копируем содержимое полей Access Key ID и Secret Key
и вставляем в файл devops/packer/app/templates/application.service.j2
заменяя плейсхолдеры

```bash
Environment=AWS_ACCESS_KEY_ID=<ACCESS KEY ID из интерфейса облака>
Environment=AWS_SECRET_ACCESS_KEY=<SECRET KEY из интерфейса облака>
```

## Сборка бинарных файлов приложения

1. Переходим в директорию server/
1. Выполняем

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../sample_server .
```

1. Повторяем тоже самое для директории worker/

## Сборка образа с помощью Packer

1. Локально должены быть установлены утилиты packer и ansible
1. Переходим в директорию devops/packer/ в этом проекте
1. Запускаем сборку образа командой

```bash
packer build -var image_tag=0.0.1 app.pkr.hcl
```

## Создание инфраструктуры с помощью Terraform

1. Локально должна быть установлена утилита Terraform
1. Переходим в директорию terraform/ в этом проекте
1. Создаем там файл vars.tfvars со следующим содержимым

```hcl
image_tag = "0.0.1"
node_count = 3
```

1. Запускаем создание инфраструктуры командой

```bash
terraform apply -var-file vars.tfvars -auto-approve
```

## Проверка

После того как терраформ отработает, он выдаст IP адрес созданного балансировщика.
По этому адресу можно делать запросы в API тестового приложения

```bash
curl <LB IP>/publish
```
