# GoBack

A sandbox project to test cool backend features.

## Requirements & Installation
- [Docker Desktop](https://www.docker.com/)
- [Golang Migrate Cli](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
```bash
make install
```

## Database
https://dbdiagram.io/d/6491f9e902bd1c4a5ecb85dd

```bash
# Create a migration.
make db-migration-create MIGRATION_NAME={MIGRATION_NAME)

# Run migration(s).
make db-migration-run
```