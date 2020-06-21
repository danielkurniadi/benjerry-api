# Builder
FROM golang:1.14.2-alpine3.11 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /home/benjerry

COPY . .

# Build application binary 
# into file `engine` 
RUN make app

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app
EXPOSE 8080

# Copy environ file from project
# and copy application binary from builder
COPY --from=builder /home/benjerry/engine /home/benjerry/.env /app/

# IMPORTANT: Provide your .env file
# containing env variables for app configs
SHELL ["/bin/bash", "-c", "source /app/.env"]
ENTRYPOINT [ "/app/engine" , "run"]
CMD ["--port=8080", "--host=0.0.0.0"]
