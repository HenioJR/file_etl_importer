# file_etl_importer

This is an open source tool to import files, do any work (data transformation or/and data cleasing) and load data on repository (like relational or non relational DB, queues).

## Features

- Read file (tested using csv and txt extensions =) )
- Transform data
- Load data in database (working to PostgreSQL connection; more databases soon)

## Dependencies

- [golang](https://golang.org/)
- [toml](https://github.com/BurntSushi/toml)
- [pq](https://github.com/lib/pq)
- [logger](https://github.com/NeowayLabs/logger)
- [httpRouter](https://github.com/julienschmidt/httprouter)

## Config

In config file, you need to configurate parameters to read file, to connect database and to processing load.

The config file is in project root path, named **config.toml**.

Bellow, an example:
```
[file]
pathName = "/home/user/file.csv"
separator = ","

[processing]
numberOfThreads = 20
batchSizeCommit = 500
validateLineByLine = true #true to validate parameters length, line by line
removeDoubleQuote = true #true to remove double quote from all lines and all atributes

[database]
[database.postgres]
driver = "postgres"
user = "x"
password = "x"
port = "x"
dbname = "x"
host = "x"
maxOpenConns = 50
maxIdleConns = 10
schemaOutput = "x" #schema to save imported data
tableOutput = "x" #table to save imported data
[database.mongo]
#to do
```

## Install and Execute

Go get it:
```sh
# Make sure GOPATH/bin is in your PATH
go get github.com/BurntSushi/toml
go get github.com/lib/pq
go get github.com/NeowayLabs/logger
go get github.com/julienschmidt/httprouter
```

Set **config.toml** file as you wish.

After that, you can build project:
```sh
# Make sure that you in project root path
go build
```

After build, you can run:
```sh
# Make sure that you in project root path
./file_etl_importer
```
