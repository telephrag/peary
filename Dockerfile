FROM golang:1.18.4 as build-env

WORKDIR /go/src/peary
ADD . /go/src/peary

RUN go get -d -v ./...
RUN go build -o /go/bin/peary .

USER 1000

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/peary /

CMD ["./peary"]

# sudo docker run --name=peary_container --env-file=.env -v /home/$(id -nu 1000)/volumes/peary_data:/data -p 8080:8080 peary