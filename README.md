# go-jarvis
This is a linux ci-cd deploy app. Just like gocd

## Installation

~~~shell
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/lvyong1985
cd $GOPATH/src/github.com/lvyong1985
git clone https://github.com/lvyong1985/go-jarvis.git
cd go-jarvis
go get ./...
./control.sh build
./control.sh start

# goto http://localhost:10010
~~~

## Configuration

~~~yaml

~~~

