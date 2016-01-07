FROM golang:1.5

RUN apt-get update

RUN go get \
	github.com/skelterjohn/rerun

RUN mkdir /code
COPY . /code

RUN mkdir -p /go/src/bitbucket.org/syb-devs
RUN ln -s /code /go/src/bitbucket.org/syb-devs/goth

RUN ls -l /go/src/bitbucket.org/syb-devs/goth
RUN cd /go/src/bitbucket.org/syb-devs/goth/examples/todo && go get

VOLUME ["/code"]

ENV PORT 80
EXPOSE 80

WORKDIR /go/src/bitbucket.org/syb-devs/goth/examples/todo

#ENTRYPOINT ["rerun", "--build", "bitbucket.org/syb-devs/goth/examples/todo"]
ENTRYPOINT ["sleep", "3600"]
