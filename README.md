# Zip-url
- an url Shortener api using golang

## Run the project
```bash
go run main.go
```

## Migrations Helper
- Initiate migrations
```bash
goose -dir migrations -s create create_table_name sql
```

- Run `migration`
```bash
make migrate-up
```
- Check `migration status`
```bash
make migrate-status
```
- Down 'migration'
```bash
make migrate-down
```