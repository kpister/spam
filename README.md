# spam-core

install go: https://golang.org/dl

`mkdir -p ~/go/src && cd $_` OR `mkdir -p some/dir/here/src && cd $_ && export GOPATH=$(some/dir/here)`

`go get -d github.com/kpister/spam`

`cd github.com/kpister/spam`

`go build` will create an executable spam 

OR `export $PATH=$PATH:$GOPATH/bin` and then `go install`. You can now run spam-core from anywhere

edit spam_core.cfg with the IPs of your peers. These will be the ones spam-core connects to by default

`spam`

(new tab) `spam`

(new tab) `spam`


for each instance:

enter a unique port (eg: 8080, 8081, 8082)

use the command `connect` to connect to a single peer (you will be prompted for address)

enter the address of the other clientservers (`127.0.0.1:8080`, `127.0.0.1:8081`, `127.0.0.1:8082`)

use the command `broadcast` to send a message to all of your peers (you will be prompted for a message)

enter your message and see it in the other nodes!
