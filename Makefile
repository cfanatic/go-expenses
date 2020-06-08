generate:
	qtmoc desktop cmd/expenses/

build:
	go build -o bin/go-expenses cmd/expenses/main.go

run:
	@go build -o bin/go-expenses cmd/expenses/main.go
	@./bin/go-expenses

deploy:
	qtdeploy build darwin cmd/expenses

clean:
	rm -f gui/moc*
	rm -f -r bin/
	rm -f -r cmd/expenses/darwin/
	rm -f -r cmd/expenses/deploy/

all: generate run
