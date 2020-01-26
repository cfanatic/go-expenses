# go-expensegui

This tool is the graphical front-end for [go-expense](https://github.com/cfanatic/go-expense).

It helps you analyze your monthly and yearly expenses by giving answers to following questions:

1. How much do you spend on average per month and year?
2. How much would you be able to set aside each month for savings?
3. What are the highest cost types, e.g. rent, car, food, etc.?
4. Where are potentials to reduce costs?

The graphical user interface shall simplify the process of importing and labeling account statement transactions.

## Requirements

Developed and tested on the following setup:

- macOS 10.15.2
- Go 1.13.4
- Docker 2.2.0.0

## Installation

Run the particular build process for one of the hosts below:

### macOS

```bash
qtmoc desktop
go build
docker pull mongo:latest
docker run -d -p 27017:27017 --name expenses mongo:latest
```

### Linux

```bash
docker pull therecipe/qt:linux_debian_9
docker build -t expensegui:latest -f Dockerfile .
docker run --name expensegui expensegui:latest
docker cp expensegui:/home/user/work/src/github.com/cfanatic/go-expensegui/deploy/linux/go-expensegui .
docker pull mongo:latest
docker run -d -p 27017:27017 --name expenses mongo:latest
```

### Windows

```bash
n/a
```

## Usage

Run the program by calling:

```bash
./go-expensegui
```

Load the sample account statements which can be found in folder `/res/` and label each transaction concerning cost type.
