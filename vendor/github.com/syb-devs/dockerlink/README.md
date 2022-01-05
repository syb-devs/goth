# Dockerlink 

[![Build Status](https://drone.io/github.com/syb-devs/dockerlink/status.png)](https://drone.io/github.com/syb-devs/dockerlink/latest)

Dockerlink is a small package to get data (IP Address, port, protocol) of a linked Docker container. 

More about [Docker linking](https://docs.docker.com/userguide/dockerlinks/)

## Usage example

Given two containers setup as follows:

```bash
$ docker run -d --name db mongo
$ docker run -d -P --name goweb --link db:db myrepo/goapp
```

Where the first one is running a MongoDB server and the second one is our Golang web app, we could use the package to get the IP address and port to connect to the MongoDB database from within our Go app in the second container:

```go
package main

import (
	"github.com/syb-devs/dockerlink"
	"gopkg.in/mgo.v2"
)

func main() {
	link, err := dockerlink.GetLink("db", 27017, "tcp")
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(fmt.Sprintf("%s:%d", link.Address, link.Port))
	if err != nil {
		panic(err)
	}
	...
}
```
