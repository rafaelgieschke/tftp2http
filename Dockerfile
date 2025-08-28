FROM golang

WORKDIR /app
COPY . .
RUN go build

CMD ["./tftp2http"]
EXPOSE 69/udp
