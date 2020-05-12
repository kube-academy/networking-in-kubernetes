FROM golang:latest AS build
COPY demoapp.go /demoapp.go
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/go-errors/errors
RUN go build /demoapp.go

FROM nicolaka/netshoot
COPY --from=build /go/demoapp /demoapp
EXPOSE 9000
CMD "/demoapp"