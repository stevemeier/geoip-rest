# geoip-rest
A simple REST-API for Maxmind GeoIP database

This piece of software (more or less) replicates https://github.com/observabilitystack/geoip-api
but without all the Java bloat. It's written in Go for easy deployment as a standalone binary.

## Compilation

```
go build main.go
```

## Running

```
./main --db <MAXMIND-GEOIP.mmdb> [ --listen 0.0.0.0:8888 ]
```

* The `db` parameter defines the location of your copy of the MaxMind GeoIP database
* The `listen` parameter is optional and allows you define on what IP address and/or port to listen

## Example

```
$ curl http://127.0.0.1:8000/72.44.167.158
{
  "country": "US",
  "latitude": 39.0991,
  "logitude": -75.5966,
  "continent": "NA",
  "timezone": "America/New_York",
  "stateprov": "Delaware",
  "stateprovCode": "DE",
  "city": "Camden"
}
```

Note: The fields `stateprov`, `stateprovCode` and `city` are optional and may not be present in each response.

## IPv6

IPv6 is supported for lookups but not for transport.

## Logging

No logs are written, neither to Stdout nor to any files.
