build:
	@go build -o bin/meal-route

run: build
	@./bin/meal-route

seed:
	@go run scripts/seed.go


test:
	@go test -v ./tests