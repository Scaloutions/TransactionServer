# Transaction Server

This server is responsible for the business logic and user authentication.

# Running the server

1. `go build` insisde the project directory
2. `./TransactionServer`

_alernatively_: 

Run `go run *.go -logtostderr=true`

Server will be running on [http://localhost:9090/](http://localhost:9090/) 

# Running with Docker image

*Build image*
inside the directory run:  `docker build -t transactionserver .`
This will build a new docker image named <transactionserver>.

*Saving Docker image*
`docker save transactionserver > PATH/dockerTS.tar`

*Exporting Docker image for .tag*
`docker load < PATH/dockerTS.tar`

*Running Docker image*
You can run `docker-compose up` inside the directory and that will bring two docker containers MySQL and transcation server.
Alternatively you can build two images separately with `docker build` command and then run
`docker run -p 9090:9090 transactionserver`

*Accessing MYSQL database container*
1. `docker exec -it <CONTAINER ID> bash`
2. `mysql -u root -p` enter password
3. Now you can execute queries

*DB Initialization*
To run initial sql script that will set up the tables you can execute inside mysql container:
`mysql -u user -p < dbSetup.sql`

# Running tests

To run tests for all packages run `go test ./...`
Go will look for *_test.go files

### Run into problems?

Make sure you have correct setup fpr $GOPATH
more about that here: [https://golang.org/doc/code.html#GOPATH](https://golang.org/doc/code.html#GOPATH)
