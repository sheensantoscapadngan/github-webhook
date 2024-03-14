FROM golang:1.22.0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /github-app
EXPOSE 8080
CMD [ "/github-app" ]
