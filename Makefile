build:
	@go build -o bin/meal-route

run: build
	@./bin/meal-route

seed:
	@go run scripts/seed.go


test:
	@go test -v ./tests

demo:
	 @go run scripts/demo/main.go
	 

docker_build:
	docker build -t m-route .	

docker_run:
	@docker run --name m-route -e PROD=true -e JWT_SECRET=____ -p 5000:5000 m-route


