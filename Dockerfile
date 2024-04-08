# Dockerfile References: https://docs.docker.com/engine/reference/builder/

#Create image to build
FROM golang:latest as build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy all content
COPY . .

# Build app
RUN go build -o main .


# Create final image
FROM golang:latest


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy only the necessary
COPY --from=build /app/main ./main
COPY --from=build /app/.resources ./.resources

# Expose by port
EXPOSE 80

# Command to run the executable
CMD ["./main"]