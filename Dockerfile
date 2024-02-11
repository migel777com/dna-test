FROM golang:1.21.1

WORKDIR /app

ADD . /app

RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /dna-test

EXPOSE 8080

# Run
CMD ["/dna-test"]
	
