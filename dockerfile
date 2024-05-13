FROM golang:1.22.2

WORKDIR /app

COPY ["go.mod", "main.go", "./"]
COPY ["solution/", "./solution/"] 

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]