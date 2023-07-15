# Use the official Golang base image
FROM golang:1.20.6

# Set the working directory inside the container
WORKDIR /app

# Copy the built application into the container
COPY ./happyapp .

# Expose the port on which the application listens
EXPOSE 3000

# Set the command to run when the container starts
CMD ["./happyapp"]