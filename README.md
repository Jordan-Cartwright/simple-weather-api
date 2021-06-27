# Simple Weather API
[![Build Status](https://travis-ci.com/Jordan-Cartwright/simple-weather-api.svg?branch=main)](https://travis-ci.com/Jordan-Cartwright/simple-weather-api) ![GitHub](https://img.shields.io/github/license/Jordan-Cartwright/simple-weather-api?color=blue) ![Docker Pulls](https://img.shields.io/docker/pulls/jordancartwright/simple-weather-api)

This is a simple Golang weather API service that can be run in a container environment. This demo application allows users to look up weather by location, using latitude and longitude coordinates.

## Prerequisites and Considerations
- Obtain an API key from OpenWeatherMap by signing up [here](https://openweathermap.org/appid).

## CI/CD Setup
This demo application makes use of TravisCI build stages to accomplish the following:
- Trigger TravisCI on every commit
- On every commit
  - Run our tests
- On commits to `main`
  - Build our simple-weather-api docker image
  - Tag the images and push thea architecture specific images
  - Create Docker manifests so users are able to run the image on multiple architectures effortlessly

This setup allows us to run our CI/CD pipeline in parallel, when applicable, and publish multi-architecture images. This therefore makes the demo application available to run in a variety of infrastructures.

## Deployment
### Deploy with Docker (locally)
We can build the Dockerfile locally using the following commands:
```
docker build -t simple-weather-api:latest .
```

Additionally, we can pull the image from DockerHub with the following command:
```
docker pull jordancartwright/simple-weather-api:latest
```

When the image has either been pulled or built, we can then run the image detached in the background
with the following:
```
docker run -d -p 8080:8080 -e "APIKEY=myapikey" jordancartwright/simple-weather-api:latest
```

> **Note:** you can optionally use the `TZ` environment variable to set the container timezone for logging
>
> Example: `-e "TZ=America/New_York"`

This will result in the image being run on the machine in the background. This can be verified with `docker ps`.
```
$ docker ps
CONTAINER ID   IMAGE                                        COMMAND   CREATED         STATUS         PORTS                    NAMES
fa2f06278d00   jordancartwright/simple-weather-api:latest   "/api"    5 seconds ago   Up 3 seconds   0.0.0.0:8080->8080/tcp   laughing_moore
```

Now that you have confirmed, you can verify the simple weather api is running by issuing the following command:
```
$ curl localhost:8080/api/v1/ping
{"message": "pong"}
```

### Deploy to Kubernetes

## Example API Output
`/api/v1/weather?latitude&longitude` is the expected input for a location (e.g., `/api/v1/weather?latitude=33.7984&longitude=-84.3883`).

It will return JSON with the current weather and 7 day forecast. The response looks like this:
```
{
  "date": "2018-01-23",
  "type": "partly-cloudy-day",
  "description": "Partly Cloudy",
  "temperature": 61.78,
  "wind": {
    "speed": 4.66,
    "bearing": 147
  },
  "precip_prob": 0,
  "daily": [
    {
      "date": "2018-01-23",
      "type": "partly-cloudy-day",
      "description": "Mostly cloudy throughout the day.",
      "temperature": {
        "low": 46.78,
        "high": 68.66
      }
    },
    ...
  ]
}
```
