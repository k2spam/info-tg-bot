# INFO_TG_BOT

## Телеграм бот для оповещений

### Скрытые переменные

```env
TOKEN= // токен телеграм бота
USERS= // список сохраненых пользователей для рассылки
```

### Пример запроса к боту

```js
fetch("http://<remote.server>:8081/order", {
    mode: "no-cors",
    method: "POST",
    body: JSON.Stringify({
        name: "User name",
        phone: "+7 (123) 456-78-90"
    })
})
```

### Билд

Билд бота в папку

```bash
go build -o /usr/local/bin/telegram-bot
```

!!! ПОСЛЕ БИЛДА НЕ ЗАБЫВАЕМ СОЗДАТЬ .ENV !!!

### Unit-файл systemd

Добавление сервиса, чтоб перезапускать бота при рестарте и падениях

```bash
sudo nano /etc/systemd/system/telegram-bot.service
```

```nano
[Unit]
Description=Telegram Bot Service
After=network.target

[Service]
ExecStart=/usr/local/bin/telegram-bot
WorkingDirectory=/usr/local/bin/
Restart=always
User=root
EnvironmentFile=/usr/local/bin/.env

[Install]
WantedBy=multi-user.target
```

### Запуск службы

```bash
sudo systemctl daemon-reload
sudo systemctl enable telegram-bot
sudo systemctl start telegram-bot
```

### Проверка статуса работы

```bash
sudo systemctl status telegram-bot
```

### Остановка и перезапуск сервиса

```bash
sudo systemctl stop telegram-bot
sudo systemctl restart telegram-bot
```

### Просмотр логов сервиса

```bash
sudo journalctl -u telegram-bot -f
```

### Если потребоваось обновить конфиг сервиса

```bash
sudo systemctl daemon-reexec
```
