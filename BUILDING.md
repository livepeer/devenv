# Building the Dev Environment

## Create an empty VirtualBox vm

```
vagrant init bento/ubuntu-16.04
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

### Install IPFS:
```
wget https://ipfs.io/ipns/dist.ipfs.io/go-ipfs/v0.4.11/go-ipfs_v0.4.11_linux-amd64.tar.gz
tar -zxvf go-ipfs_v0.4.11_linux-amd64.tar.gz
cd go-ipfs
sudo ./install.sh
cd ~
```

### Install Livepeer Release
```
wget https://github.com/livepeer/go-livepeer/releases/download/0.1.6/livepeer_linux.tar
tar -xf livepeer_linux.tar
```

### Update .profile:
```
echo 'PATH="$PATH:/usr/local/go/bin"' >> ~/.profile
echo 'PATH="$PATH:$HOME/go/bin"' >> ~/.profile
echo 'PATH="$PATH:$HOME/livepeer_linux"' >> ~/.profile
source ~/.profile
```

### Clean up
```
cd ~
rm *.tar*
rm -rf go-ipfs/
```

#### Remove apt cache
```
sudo apt-get clean -y
sudo apt-get autoclean -y
```

#### Remove bash history
```
unset HISTFILE
sudo rm -f /root/.bash_history
rm -f /home/vagrant/.bash_history
```

#### Remove log files
```
sudo find /var/log -type f -exec truncate -s 0 {} \;
```

#### Packaging the box

```
vagrant package --base livepeer-ubuntu1604
```

#### Run the packaged box
```
vagrant box add ...
vagrant init ...
vagrant up
```
