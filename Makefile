BINARY=engine

app:
	go build -o ${BINARY} engine/*.go

