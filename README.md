# spam

install go: https://golang.org/dl

`mkdir -p ~/go/src && cd $_` OR `mkdir -p some/dir/here/src && cd $_ && export GOPATH=$(some/dir/here)`

`go get -d github.com/kpister/spam`

`cd github.com/kpister/spam`

`go build` will create an executable spam 

OR `export $PATH=$PATH:$GOPATH/bin` and then `go install`. You can now run spam-core from anywhere

edit spam_core.cfg with the IPs of your peers. These will be the ones spam-core connects to by default

`spam` will default to use spam_core.cfg and .log, if you change the name of the cfg file, the console will only work if you do the following:

`spam -c .log_new_cfg_file_name.cfg`

Once you have the console running, you can check your peer status by typing `peer`

For a test network to show it is working?

`spam -i node1.cfg`

(new tab) `spam -i node2.cfg`

(new tab) `spam -i node3.cfg`

Feel free to close and reconnect each of the nodes.

Expected behavior is that the nodes will connect within a couple of seconds and start broadcasting their current clock time every five seconds.
