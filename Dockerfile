FROM golang:latest
WORKDIR /app/auth-backend
COPY ./ ./
ENTRYPOINT [ "go", "run", "main.go" ]