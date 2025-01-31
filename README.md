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
go build /usr/local/bin/telegram-bot -o main.go
```

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
Environment="BOT_TOKEN=your_token_here"

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
