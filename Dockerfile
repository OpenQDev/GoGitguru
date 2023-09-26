# Start from the official Go image
FROM golang:1.21.1

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the application
RUN go build -o main .

# Command to run the application
ENTRYPOINT ["./main"]
