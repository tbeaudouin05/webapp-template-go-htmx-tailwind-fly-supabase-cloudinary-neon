ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
# build the app binary and save the binary at root folder
RUN go build -v -o /run-app .


FROM debian:bookworm

RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

# Copy the built binary to executable path
COPY --from=builder /run-app /usr/local/bin/

# Copy the static directory to the root
COPY --from=builder /usr/src/app/static /static

# Copy the gemini-key.json to the root
COPY --from=builder /usr/src/app/gemini-key.json /gemini-key.json

# run the app from the root folder (the command will follow executable path but static folder etc. should be at root folder)
CMD ["run-app"]
