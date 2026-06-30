DB_DSN=host=localhost port=5432 user=postgres dbname=zip_url_db  password=postgres sslmode=disable

migrate-up:
	@goose -dir migrations postgres "$(DB_DSN)" up

migrate-down:
	@goose -dir migrations postgres "$(DB_DSN)" down

migrate-status:
	@goose -dir migrations postgres "$(DB_DSN)" status

migrate-create:
	@goose -dir migrations create $(name) sql