# purpleair-go
## about
Simple Telegram Bot which provides PM2.5 AQI index from PurpleAir Sensors.

- Data source: https://www.purpleair.com/json
- Data source documentation: https://www2.purpleair.com/community/faq#hc-json-object-fields
- Sensor ID is hardcoded as a constant: `SensorID`. Feed free to modify it.

## build & run
You should create Telegram bot and obtain `PURPLEAIR_BOT_TOKEN` first. See instructions here: https://core.telegram.org/bots#6-botfather

`PURPLEAIR_BOT_SENSOR_ID` might be taken from https://www.purpleair.com/map. Find your sensor, click on _Get This Widget_ and you'll see Sensor ID.
```bash
docker build . -t purpleair-go
docker run -it --env PURPLEAIR_BOT_TOKEN="<Telegram Bot API Token>"  --env PURPLEAIR_BOT_SENSOR_ID="<Sensor ID>" purpleair-go
```

## docker hub (alternative)
```bash
docker pull maffei/purpleair-go
docker run -it --env PURPLEAIR_BOT_TOKEN="<Telegram Bot API Token>" --env PURPLEAIR_BOT_SENSOR_ID="<Sensor ID>" purpleair-go
```