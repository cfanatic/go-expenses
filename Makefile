generate:
	qtmoc desktop

build:
	go build

deploy:
	go build
	./go-expensegui

clean:
	rm gui/moc*
	rm go-expensegui

all: generate deploy
