# Extensible Web3 Indexer

### ✅ This is a web3 Transaction Indexer (Block Parser) POC, it's designed to be completely extendable to:
- Add other EVM chains
- Use any kind of Databases
- Provide any kind of Public Interfaces (http, grpc, cli, etc.)
- ...


## 📜 Description

This Project Implemented Based on Clean Architecture in Golang.

🔰 Rule of Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

📚 More at [Uncle Bob clean-architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

📚 More at [Martin Fowler PresentationDomainDataLayering](https://martinfowler.com/bliki/PresentationDomainDataLayering.html)

### 🗺 The Diagram: 
![clean architecture](https://github.com/apm-dev/vending-machine/blob/main/clean-arch.png)


### 🏃🏽‍♂️ How To Run This Project
⚠️ Since the project already use Go Module, I recommend to put the source code in any folder except GOPATH.

#### 🧪 Run the Testing

```bash
$ make test
```

#### 🐳 Run the Applications
Here is the steps to run it with `docker-compose`

```bash
# move to directory
$ cd workspace
# Clone it
$ git clone https://github.com/apm-dev/evm-tx-parser.git
# move to project
$ cd evm-tx-parser
# (optional) Build the docker image first
$ make docker
# Run the application
$ make run
# check if the containers are running
$ docker ps
# See the logs
docker logs --follow evm_tx_parser
# subscribe an address
$ curl -X POST http://localhost:8000/address -d '{"address":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"}'
# get txs of an address
$ curl -X GET http://localhost:8000/address/0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD/txs
# get last parsed block
$ curl -X GET http://localhost:8000/block/last
# Stop
$ make stop
```

### 🛠 Tools Used:
In this project, I use some tools listed below. But you can use any similar library that have the same purposes. But, well, different library will have different implementation type. Just be creative and use anything that you really need. 

- All libraries listed in [`go.mod`](https://github.com/apm-dev/web3-api-proxy/blob/develop/go.mod)
- ["github.com/vektra/mockery".](https://github.com/vektra/mockery) To Generate Mocks for testing needs.
