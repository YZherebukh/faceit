FROM ubuntu:20.04

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /usr/src/app

# Copy the code into the container
COPY . .

# Build the application
# RUN go build 

RUN apt-get update -y
RUN apt-get install -y make
RUN apt-get install curl -y
RUN curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
RUN chmod +x /usr/local/bin/dbmate
RUN chmod +x test

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["./test"]
