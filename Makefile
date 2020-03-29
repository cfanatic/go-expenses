generate:
	qtmoc desktop cmd/expenses/

build:
	go build -o bin/go-expenses cmd/expenses/main.go

deploy:
	@go build -o bin/go-expenses cmd/expenses/main.go
	@./bin/go-expenses

clean:
	rm -f gui/moc*
	rm -f -r bin/

all: generate deploy
