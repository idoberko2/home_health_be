# Home HealthCheck Application BackEnd

This application receives `/ping` requests from any device that is connected to a power supply and internet.

If these requests aren't received for a pre-defined amount of time (`GracePeriod`), a notification is sent to a Telegram chat, alerting that the state has changed to "Unhealthy".

Once a request is received, another notification is sent to the Telegram chat, notifying the state is restored to "Healthy".

## Environment variables

| Variable      | Description | Default |
| ----------- | ----------- | -------- |
| HC_HISTORY_LENGTH      | The length of the history to track       | `1000` |
| HC_PASSPHRASE   | A string that should be included in every ping request        | - |
| HC_GRACE_PERIOD   | How long to wait before notifying state has changed to "Unhealthy"        | `5m` |
| HC_SAMPLE_RATE   | How often to check for an "Unhealthy" state        | `10s` |
| HC_TOKEN   | Telegram Bot token to send the notifications        | - |
| HC_CHAT_ID   | Telegram chat id to send the notifications to        | - |
| HC_IS_DEBUG   | Whether or not to turn on Telegram Bot debug mode        | `false` |
| HC_IS_LISTEN   | Setting this to "true" is used only to find out a chat's id        | `false` |
| HC_AVOID_DOT_ENV   | Whether or not to load a .env configuration file        | `false` |
| HC_HOST   | The server host        | `localhost` |
| HC_PORT   | The server port        | - |
| HC_WRITE_TIMEOUT   | The server write timeout        | `10s` |
| HC_READ_TIMEOUT   | The server read timeout        | `10s` |
| HC_IDLE_TIMEOUT   | The server idle timeout        | `60s` |
| HC_SERVER_SHUTDOWN_TIMEOUT   | The timeout for gracefully shutting down the server        | `5s` |

## Test
### Unit tests
```
$ go test ./...
```

### Unit and Integration tests
The integration test require a user to manually assert the incoming Telegram notifications. Once starting the test you should receive:
1. An immediate message notifying state changing to "Healthy"
1. Another message, 5 seconds later notifying state change to "Unhealthy"
```
$ go test -tags=integration ./...
```
