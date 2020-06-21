# Ben & Jerry REST API
Building an API for Ben & Jerry's (ice cream) fans.

## Changelog

- **v1**: checkout to the [v1 branch](https://github.com/iqdf/benjerry-api/tree/master) <br>
  This master branch is v1. Proposed on June 2020 <br>
  
## Overviews <a name="overviews"></a>
The following are the content of this docs.

1. [This README](#overviews) - Overview of this project
   1. [Overviews](#overviews)
   2. [Quick Setup and Run](#setup)
2. [API Schema](docs/api/) - Sample Requests and Response schema
3. [Design Docs](docs/design) - Software architecture documentation

### Description
This is an example of implementation of Clean Architecture in Go (Golang) projects.

Rule of Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

This project has  4 Domain layer. More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html
 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Delivery Layer

![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)

### What features are included in example?
* **HTTP REST API**: Implementation of simple CRUD resource management api that communicates using HTTP.
* **Basic Authentication**: Basic authentication using `username` and `password` passed through Basic Auth Header for login.
* **Session based token**: After authenticated, clients will receive session token that can be used to authenticate. Using Redis cache to save and purge expired token.
* **Common Middlewares** [WIP] : Example implementation of using middleware. The middleware includes auth and role/permission check, logging, and http header (add content-types, CORS, etc.). 
* **Database Mongo**: Example implementation of database layer using mongo DB.
* **Dockerize Deployment** Simple Dockerfile and Docker-compose to run mongoDB, Redis, and the application.

### Dependencies
* Golang and Go Pkg under `go.mod`
* Mongo DB: NoSQL Database to store products (ice cream) document
* Redis Cache: In-memory cache for storing and managing session token
* Docker and Docker-Compose [optionally]: For containerised deployment

## Quick Setup and Run <a name="setup"></a>
> Make sure you have Go installed in your machine.

Since the project already use Go Module, I recommend to put the source code in any folder but `$GOPATH`.

### Run Tests

The tests cover almost all the service, but some functionalities in `common` and `repository` layers are yet to be tested.
```bash
$ make test
```

### Run the Applications

We provide two ways to deploy the application: using baremetal and using docker-compose.

#### Running on Unix/Linux based machine
Here is the steps to get it running in local machine. The requirement for running in linux are golang installed, mongodb (started and running at `localhost:27017`), and redis (started and running at `localhost:6379`).

1. Prepare your and source the `.env` file. This contains environment variables to configure the application, e.g. define URI for database/redis. See example `.env` file provided along with this project.

```bash
# sample envvars for app configurations
export ENV_MODE=development
export DB_URI=mongodb://localhost:27017/tutorialDB
export REDIS_URI=redis://localhost:6379
```
2. Build the binary file and run
The application will run at `localhost:8080` by default.
```bash
make app-run
```

#### Running from Docker Compose
Here is the steps to run it with `docker-compose`.

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/iqdf/benjerry-api.git benjerry

# move to project
$ cd benjerry

# To build the docker image first
$ make docker-build

# To run the application using docker compose
$ make compose-run

# To check that containers are running
$ docker ps -a

# To stop docker compose that ran the app
$ make compose-stop
```

## TO DO Work and Features
- [ ] Refactor auth and role middleware. Issue #5
- [ ] Tests: Middlewares, Auth Service, and Database layer. Issue #6 #7
- [ ] Documentations: APi Schema, Software architecture design. Issue #1 #2 #3 #4

## Tools Used:
In this project, I use some tools listed below. But you can use any simmilar library that have the same purposes. But, well, different library will have different implementation type. Just be creative and use anything that you really need. 

- All libraries listed in [`go.mod`](https://github.com/iqdf/benjerry-api/blob/master/go.mod) 
- ["github.com/vektra/mockery".](https://github.com/vektra/mockery) To Generate Mocks for testing needs.
