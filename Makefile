.PHONY: clean install unittest build app app-run docker-build compose-run compose-stop vendor lint-prepare lint

BINARY=engine
ENVFILE=.env
test: 
	go test -v -cover -covermode=atomic ./...

app:
	go build -o ${BINARY} app/*.go

app-run:
	go build -o ${BINARY} app/*.go
	. ./${ENVFILE}  # provide your environ file as `.env`
	exec ./${BINARY} run --port=8080 --host=localhost 

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker-build:
	docker build --rm -t benjerrysample:latest .

compose-run:
	docker-compose up --build -d

compose-stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...
