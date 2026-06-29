# Zip-url
- an url shortener api using golang

## Migrations Helper
```bash
goose -dir migrations -s create create_table_name sql
```

- Run `migration`
```bash
make migrate-up
```
- Check `migration status`
```bash
make migrate-down
```