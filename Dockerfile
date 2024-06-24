FROM golang:1.22

RUN mkdir /AllChat

WORKDIR /AllChat

COPY go.mod go.sum ./

COPY . .

RUN go build cmd/main.go

EXPOSE 8080

CMD ["./main"]