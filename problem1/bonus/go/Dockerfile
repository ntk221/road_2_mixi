FROM golang:1.18

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app
COPY  . .
EXPOSE 1323

CMD ["air"]
