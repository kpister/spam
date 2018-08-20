# spam
_Always check the corresponding branch for details on actual updates_

Many of the following tasks are not yet implemented on Master branch and can only be used on their sub-branch

# TODO
* - [x] Console
  * - [x] Implement simple file io console
  * - [ ] Brainstorm other ideas which might be better
  
* - [ ] Authentication:
  * - [x] Read and understand RSA aglorithm
  * - [x] Look into Go-lang libraries which will help
  * - [x] Implement RSA on a branch
    * - [x] Easy way to generate private and public keys
    * - [x] Encrypt and decrypt message functions
    * - [x] Link encryption and decryption to identify for authentication
      * - [x] Send encrypted messages
      * - [x] Receive and decrypt the messages - creating the association
  * - [ ] Improve on the algorithms
  * - [ ] Add digital signing to all messages
  * - [ ] Test more
  * - [ ] We will want to consider some policies for openness
    * Accept all
    * Accept people I have in cfg
    * Accept something else
  * - [ ] It would be nice to shrink the public private key values. There is a lot of documentation on this

* - [ ] Implement SCP
  * - [x] Reread SCP paper
  * - [ ] Read Stellar's SCP implementation
  * - [ ] Draw out how we will do it in spam (including classes, file directory structure, etc)
  * - [ ] Start writing code

# To install and use

install go: https://golang.org/dl

`mkdir -p ~/go/src && cd ~/go/src && export GOPATH=~/go`

`go get -d github.com/kpister/spam`

`cd github.com/kpister/spam`

`export PATH=$PATH:$GOPATH/bin` and then `go install`. You can now run spam from anywhere


edit spam_core.cfg with the IPs and public keys of your peers. These will be the ones spam connects to by default

`spam` will default to use spam_core.cfg and .log

For a test network to show it is working run the following nodes and cfgs

`spam -i cfgs/node1.cfg`

(new tab) `spam -i cfgs/node2.cfg`

(new tab) `spam -i cfgs/node3.cfg`

Feel free to close and reconnect each of the nodes.

Expected behavior is that the nodes will connect within a couple of seconds and start broadcasting their current clock time every five seconds.

# Using the console

There is a simple command line console with which you can control your node. Run it with `spam -c`

Some example commands include: 

* `peers` - receive a list of all your peers and their status
* `broadcast` - send out a message to all your connected peers
* `drop peer by name` - remove a peer from your list (you will not reconnect to them)
* `add peer` - add a peer (using [address name publickey] format


If you change the name of the cfg file when you run your node, the console will only work if you do the following:

`spam -c .log_newCfgFileName.cfg`

#### Tags
go, plsyssec
