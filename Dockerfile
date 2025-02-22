FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o /dofusnoobsidentifier

CMD [ "/dofusnoobsidentifier" ]
