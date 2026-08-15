FROM golang:1.26-alpine

# Define work directory
WORKDIR /var/www

COPY src/go.mod ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

# Copy the app to container
COPY src/*.go ./

EXPOSE 8080

CMD ["air"]