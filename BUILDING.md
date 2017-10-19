# Building the Dev Environment

## Create an empty VirtualBox vm

```
vagrant init ubuntu/xenial64
vagrant box update
vagrant up
```

## Edit Vagrantfile as needed

Things to consider:
* resources (cpu/memory)
* networking configuration

The latter is unique for Livepeer since testing will often involve
interacting with the virtual machine as if it is both localhost and a
separate remote machine.

See this option to switch between them:
```
# config.vm.network "public_network"
```

## Configure the virtual machine

```
vagrant ssh
```

### Update base packages
```
sudo apt-get update
sudo apt-get install build-essential pkg-config cmake git checkinstall ack-grep curl vim
```

*Note:* curl, git, and vim are needed if this is a Chromium chroot.

### Install go 1.9.1 :
```
wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz
```

### Install node 8.x :
```
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### Install ffmpeg-static:
```
wget https://raw.githubusercontent.com/livepeer/ffmpeg-static/master/bin/linux/x64/ffmpeg
chmod 755 ffmpeg
sudo mv ffmpeg /usr/local/bin/
```

### Update .profile:
```
echo 'PATH="$PATH:/usr/local/go/bin"' >> ~/.profile
echo 'PATH="$PATH:$HOME/go/bin"' >> ~/.profile
source ~/.profile
```

### Install IPFS:
```
wget https://ipfs.io/ipns/dist.ipfs.io/go-ipfs/v0.4.11/go-ipfs_v0.4.11_linux-amd64.tar.gz
tar -zxvf go-ipfs_v0.4.11_linux-amd64.tar.gz
cd go-ipfs
sudo ./install.sh
cd ~
```

### Clean up
```
cd ~
rm *.tar.gz
rm -rf go-ipfs/
```

### Install geth / go-ethereum:
```
sudo apt-get install software-properties-common
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install ethereum
```

### Install truffle:
```
sudo npm install -g truffle
```

### Clone git repos:
```
mkdir src && cd src
git clone https://github.com/livepeer/go-livepeer.git
git clone https://github.com/livepeer/testenv.git
git clone https://github.com/livepeer/protocol.git
cd ~
```

### Compile and install the latest `livepeer` and `livepeer_cli`
```
go get github.com/livepeer/go-livepeer/cmd/livepeer
go get github.com/livepeer/go-livepeer/cmd/livepeer_cli
```
