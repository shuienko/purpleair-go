# purpleair-go
## about
Simple Telegram Bot which provides PM2.5 AQI index from PurpleAir Sensors.

- Data source: https://www.purpleair.com/json
- Data source documentation: https://www2.purpleair.com/community/faq#hc-json-object-fields
- Sensor ID is hardcoded as a constant: `SensorID`. Feed free to modify it.

## build & run
You should create Telegram bot and obtain **_Telegram Bot API Token_** first. See instructions here: https://core.telegram.org/bots#6-botfather
```bash
docker build . -t purpleair-go
docker run -it --env PURPLEAIR_BOT_TOKEN="<Telegram Bot API Token>" purpleair-go
```

## docker hub (alternative)
```bash
docker pull maffei/purpleair-go
docker run -it --env PURPLEAIR_BOT_TOKEN="<Telegram Bot API Token>" purpleair-go
```