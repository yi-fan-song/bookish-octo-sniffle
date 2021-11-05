# octo


## Setup
Build and run the docker container
```shell
~$ docker build -t octo:latest .
~$ docker run --rm -it -p 8080:8080 octo:latest
```
The server will be listening on port 8080

## Endpoints

### GET `/timer/:id`

- `id` is a positive integer

A succesful request return 200 and the response looks like
```json
{
    "id": "1",
    "duration": 250,
    "timeRemaining": 244,
    "expiresOn": "05 Nov 21 16:12 UTC",
    "isPaused": false,
    "webhookUrl": "http://localhost:8080/status"
}
```

### POST `/timer`
Request body should be
```json
{
    "duration": 250,
    "webhookUrl": "http://localhost:8080/status"
}
```
- `duration` is in seconds
- when the timer is up, a POST request will be made to `webhookUrl` with body
```json
{"message": "Time's up!"}
```

A succesful request return 200 and the response looks like
```json
{
    "id": "1",
    "duration": 250,
    "timeRemaining": 250,
    "expiresOn": "05 Nov 21 16:12 UTC",
    "isPaused": false,
    "webhookUrl": "http://localhost:8080/status"
}
```

### PUT `/timer/:id`
This endpoint is used to pause and unpause a timer

The request body should be
```json
{
    "isPaused": false
}
```

A succesful request return 200 and the response looks like
```json
{
    "id": "1",
    "duration": 250,
    "timeRemaining": 151,
    "expiresOn": "05 Nov 21 16:12 UTC",
    "isPaused": true,
    "webhookUrl": "http://localhost:8080/status"
}
```

## Configuration
You need to create a .env file, default like so should work
```
DB_HOST=localhost
DB_USER=dbuser
DB_PASSWORD=Vj23urju
DB_NAME=pomodoro
DB_PORT=5432
```

If `DB_USER` and `DB_PASSWORD` are different, make sure to update them in the `run.sh` script as well
