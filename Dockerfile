# Use the official Go image as the base image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go application
RUN go build -o main

# Expose port 8080 for the HTTP server
EXPOSE 8000

# Run the Go application
CMD ["./main"]
