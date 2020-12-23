# BitHyve Wrapper

BitHyve wrapper is a server instance that wraps around [electrs](https://github.com/Blockstream/electrs) to provide some additional functionality on top of electrs as required by [hexa](https://github.com/bithyve/hexa)

## Getting Started

1. Install golang (replace 1.15.4 with your favorite version)

```
wget https://dl.google.com/go/go1.15.4.linux-amd64.tar.gz
sudo tar -xvf go1.15.4.linux-amd64.tar.gz
sudo mv go /usr/bin/
```

2. Update path

```
export GOROOT=/usr/bin/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

3. Download and install wrapper

```
git clone https://github.com/bithyve/bithyve-wrapper.git && cd bithyve-wrapper
go install github.com/bithyve/bithyve-wrapper && go build
sudo ./BitHyve-wrapper -t
```
