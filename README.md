[![Travis Build Status](https://travis-ci.org/radoondas/safecastbeat.svg?branch=master)](https://travis-ci.org/radoondas/safecastbeat)

# Safecastbeat

Welcome to Safecastbeat.

Safecastbeat is a beat which will periodically pull data from [Safecast](https://safecast.org/) [API](https://api.safecast.org/en-US/home). 

This will specifically pull data every `Period` defined and will use api call with `since` parameter.
Example call may look following:

```
https://api.safecast.org/measurements.json?per_page=1000&since=2019-03-22%2B20%3A51%3A46
```

Above command will request any new data added to Safecast DB since date and time specified in URI. It can be mix of current and older data as measurements can be uploaded with delays.

## Installation
Download and install appropriate package for your system. Check release [page](https://github.com/radoondas/safecast/releases) for latest packages.

You also can use Docker image `docker pull radoondas/safecastbeat`


## Configuration

To run Safecastbeat you have to define `Period` for data pull. 1m or 2m should be sufficient.

```yaml
  period: 5m
```

Define the path to CA file which requires TLS call. One CA is provided in the repository. Feel free to use it.

## Run

```
./safecastbeat -c safecastbeat.yml -e 
```

## Visualisations
This is an example of visualisation for measurements.

![Map](docs/img/map_01.png)

## Build
If you want to build Owmbeat from scratch, follow [build](BUILD.md) documentation.
