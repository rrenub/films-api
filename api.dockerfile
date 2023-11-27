FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and env variables
COPY . .

# Build
RUN go build -o ./movies-api -v ./src/api

# Expose port
EXPOSE 4000

# Run
CMD ["./movies-api"]