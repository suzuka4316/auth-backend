FROM golang:alpine

# Install git

RUN apk update && apk add --no-cache git

WORKDIR /app/auth-backend

COPY go.mod go.sum ./

RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

ENTRYPOINT [ "go", "run", "main.go" ]