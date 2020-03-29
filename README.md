# go-expenses

This program helps you analyze your monthly and yearly expenses by giving answers to following questions:

1. How much do you spend on average per month and year?
2. How much would you be able to set aside each month for savings?
3. What are the highest cost types, e.g. rent, car, food, etc.?
4. Where are potentials to reduce costs?

Labeling is done in automatic fashion whenever a similar transaction is found in the database.

Three sample account statements can be found in folder `cmd/expenses/resources/`.

## Requirements

Developed and tested on the following setup:

- macOS 10.15.2
- Go 1.13.4
- Docker 2.2.0.0

Consider that this program is not deployed as a Go module, because of building issues related to the Qt bindings.

So make sure to check out following library versions:

| Package                                   | Version                              |
| ----------------------------------------- | ------------------------------------ |
| github.com/360EntSecGroup-Skylar/excelize | v2.0.2-0.20190924135319-a34d3b8c86d6 |
| github.com/gonum/stat                     | v0.0.0-20181125101827-41a0da705a5b   |
| github.com/ryanuber/columnize             | v1.1.2-0.20190319233515-9e6335e58db3 |
| github.com/wcharczuk/go-chart             | v2.0.2-0.20190910040548-3a7bc5543113 |

## Installation

Run the particular build process for one of the hosts below:

### macOS

```bash
make generate
make build
```

### Linux

```bash
docker pull therecipe/qt:linux_debian_9
docker build -t expenses:latest -f Dockerfile.linux .
docker run --name expenses expenses:latest
docker cp expenses:/home/user/work/src/github.com/cfanatic/go-expenses/deploy/linux/go-expenses .
```

### Windows

```bash
docker pull therecipe/qt:windows_64_static
docker build -t expenses:latest -f Dockerfile.win .
docker run --name expenses expenses:latest
docker cp expenses:/home/user/work/src/github.com/cfanatic/go-expenses/deploy/windows/go-expenses.exe .
```

## Usage

Start the MongoDB database:

```bash
docker pull mongo:latest
docker run -d -p 27017:27017 --name transactions mongo:latest
```

Run the program.

Import each account statement and specify cost type labels for all transactions:

![Data_Label](https://raw.githubusercontent.com/cfanatic/go-expenses/master/cmd/expenses/resources/go-expenses-1.png)

Show transaction details by hovering over Payee fields:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expenses/master/cmd/expenses/resources/go-expenses-3.png)

Exclude cost groups from statistical analysis by right-clicking on Label fields:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expenses/master/cmd/expenses/resources/go-expenses-2.png)

## Analysis

Analyze the monthly expenses for a particular account statement:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expenses/master/cmd/expenses/resources/go-expenses-4.png)

Show the complete analysis for all account statements which are available in the database:

![Data_Exclude](https://raw.githubusercontent.com/cfanatic/go-expenses/master/cmd/expenses/resources/go-expenses-5.png)
