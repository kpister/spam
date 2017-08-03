# spam-core

install go: https://golang.org/dl

`mkdir -p ~/go/src && cd $_` OR `mkdir -p some/dir/here/src && cd $_ && export GOPATH=$(some/dir/here)`

`go get -d github.com/kpister/spam-core`

`cd github.com/kpister/spam-core`

`go build` will create an executable spam-core 

OR `export $PATH=$PATH:$GOPATH/bin` and then `go install`. You can now run spam-core from anywhere

`spam-core`

(new tab) `spam-core`

(new tab) `spam-core`


for each instance:

enter a unique port (eg: 8080, 8081, 8082)

use the command `connect` to connect to a single peer (you will be prompted for address)

enter the address of the other clientservers (`127.0.0.1:8080`, `127.0.0.1:8081`, `127.0.0.1:8082`)

use the command `broadcast` to send a message to all of your peers (you will be prompted for a message)

enter your message and see it in the other nodes!
