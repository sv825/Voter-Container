# Use an official Go runtime as a parent image
FROM golang:1.16

# Set the working directory
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Voter API command inside the container
RUN go build -o voterapiredis

# Run the voterapi command when the container starts
CMD ["./voterapiredis"]