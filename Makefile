include .env

run: 
	go run cmd/main.go

migrationsup:
	migrate -path db/migrations -database "$(DBDRIVER)://$(USER):$(PASSWORD)@$(HOST):$(PORT)/$(DBNAME)?sslmode=disable" -verbose up 

migrationsdown:
	migrate -path db/migrations -database "$(DBDRIVER)://$(USER):$(PASSWORD)@$(HOST):$(PORT)/$(DBNAME)?sslmode=disable" -verbose down 