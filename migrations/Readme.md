## Migrations
- This service is responsible for managing PostgreSQL DB migrations.
- We are using pg-migrate to run `up` and `down` migrations.

# Creating migration files
- To create .sql migration file use npm script `npm run migrate:create <migration name>`. By default `node-pg-migrate` uses .js ext.
- To run the migration use `npm run migrate:up`.
- For more info about `node-pg-migrate` see https://salsita.github.io/node-pg-migrate/#/cli