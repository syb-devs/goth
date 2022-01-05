FROM golang:1.5

RUN apt-get update

RUN go get \
	github.com/skelterjohn/rerun

RUN mkdir -p /go/src/github.com/syb-devs/goth
COPY . /go/src/github.com/syb-devs/goth
RUN cd /go/src/github.com/syb-devs/goth/examples/todo && go get

RUN ln -s /go/src/github.com/syb-devs/goth /code
VOLUME ["/code"]

ENV PORT 80
EXPOSE 80

WORKDIR /go/src/github.com/syb-devs/goth/examples/todo

ENTRYPOINT ["rerun", "--build", "github.com/syb-devs/goth/examples/todo"]
