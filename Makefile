run:
	air -c .air.toml

dev:
	go run cmd/main.go

clean:
	rm -rf tmp/
	go clean -cache

gen:
	wire gen ./api
