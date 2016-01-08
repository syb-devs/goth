FROM golang:1.5

RUN apt-get update

RUN go get \
	github.com/skelterjohn/rerun

RUN mkdir -p /go/src/bitbucket.org/syb-devs/goth
COPY . /go/src/bitbucket.org/syb-devs/goth
RUN cd /go/src/bitbucket.org/syb-devs/goth/examples/todo && go get

RUN ln -s /go/src/bitbucket.org/syb-devs/goth /code
VOLUME ["/code"]

ENV PORT 80
EXPOSE 80

WORKDIR /go/src/bitbucket.org/syb-devs/goth/examples/todo

ENTRYPOINT ["rerun", "--build", "bitbucket.org/syb-devs/goth/examples/todo"]
