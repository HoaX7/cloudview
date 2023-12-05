## App service
- This service is written in Go lang. It is responsible for serving important business logic of the app.

# Database and cache being used
- Mainly using postgres and redis cache

## Generating Models and Tables
- We are using `go-jet` to generate models after running migrations.
- Use cli command `jet -source=postgres -host=localhost -password=postgres -user=postgres -port=5432 -dbname=cloudview -schema=public -path=./dbschema` to scan through all tables in the DB to generate models.
- For more info see https://github.com/go-jet/jet

## Logging
- checkout https://pkg.go.dev/log/slog
- also see Logrus and Rollbar