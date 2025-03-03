FROM golang:1.23-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o esco-search .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates
RUN apk add --no-cache libstdc++

RUN apk --no-cache add curl

# Install ollama to manage ai models
RUN curl -fsSL https://ollama.com/install.sh -o /usr/bin/ollama | sh
RUN chmod +x /usr/bin/ollama

# download granite model from ollama (can take up more than 1 minute)
RUN ["ollama", "pull", "granite-embedding:278m"]

COPY --from=build /app/esco-search .
# Also copy the .env file (for now) and data directory
COPY --from=build /app/.env .
COPY --from=build /app/data .

EXPOSE 8081

#STOPSIGNAL SIGTERM
ENTRYPOINT ["./esco-search"]
