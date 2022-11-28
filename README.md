# Pinger WIP

Small pinging selfcontained service that both allows for pinging arbitrary domains and allows to respond to them.
In case one of target fails - service will send notification to predefined user in telegram.
For telegram notifications, you'll need to set up bot accordingly

This app consists of 3 distinct parts - pinger, telegram bot handler and http server.
HTTP server is enabled always. Pinger can be enabled via `PINGER_ENABLED` env var. Bot handler can be enabled using `BOT_LISTENER_ENABLED`