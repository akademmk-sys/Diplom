    FROM golang:1.27.1

    ENV CGGO_ENABLED=1

    WORKDIR /app

    COPY go.mod go.sum ./

    RUN go mod download

    COPY . .

    RUN go build -o todo_app .

    ENV TODO_PORT=7540

    ENV TODO_DBFILE=/app/scheduler.db

    ENV TODO_PASSWORD=myPass

    EXPOSE 7540

    CMD [ "./todo_app" ]