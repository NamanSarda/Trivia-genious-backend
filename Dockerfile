# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the contents of the entire project into the container
COPY . .

# Install any needed dependencies
RUN go mod download

# Build the Go app
RUN go build -o main ./cmd/quiz-api

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
