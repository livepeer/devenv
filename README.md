# testenv
Setting up a local test environment.  
- Local Geth node and a test network
- Deploy protocol to the test network
- Start broadcaster
- Start transcoder

To use this tool, clone it and run `go run *.go`

# Dev Environment

**Note:** This is currently a work in progress. The previous `testenv` will move to support this model and deprecate the instructions above once the proper tooling is in place.

## Installation

Install Vagrant 2.0.x: https://www.vagrantup.com/

Install VirtualBox 5.0.x or 5.1.x: https://www.virtualbox.org/

Clone the repo and run this command:

```
git clone https://github.com/livepeer/testenv.git
cd testenv
vagrant up
```

Then connect to the virtual machine:

```
vagrant ssh -- -A
```

## Usage

**Protip:** Virtual machines are cheap. Don’t feel tied to this directory. Make new directories and copy the Vagrantfile to start fresh at anytime. Or `vagrant destroy` to blow away the vm in this directory.

The following software is included in this virtual machine:

```
Welcome to the Livepeer Dev Environment
Based on: Ubuntu 16.04.3 LTS (GNU/Linux 4.4.0-96-generic x86_64)

Contains:
	* Go: 1.9.1 linux/amd64
	* Node: 8.7.0
	* npm: 5.4.2
	* Geth: 1.7.2-stable-1db4ecdc
	* Truffle: v3.4.11
	* TestRPC: n/a
	* ffmpeg: 3.1-static (livepeer/ffmpeg-static)
	* livepeer: 0.1.2
	* livepeer_cli: 0.0.0

Repos:
	* go-livepeer: master (2e545787a)
	* protocol: master (8fdec47a0)
	* testenv: master (7c2aa5896)

Documentation: https://github.com/livepeer/testenv
```

The following directories are part of the `ubuntu` user’s $HOME:

### ~/go
This is the $GOROOT for this virtual machine.

### ~/src
This is the local copies of various source repositories. Currently it includes a handful of repos. Feel free to add more as needed.

### Vagrant Commands
These commands are useful for operating a virtual machine run by Vagrant. Run them in the directory that contains your Vagrantfile.

- If you want to init or boot it: `vagrant up`
- If you want to stop it: `vagrant halt`
- If you want to pause it: `vagrant suspend`
- If you want to blow it away and start over: `vagrant destroy`
- If you want to reboot the machine: `vagrant reload`
_Reload is useful if you change the Vagrantfile configuration and want to apply those changes._

## Building

Step-by-step instructions for building live in [BUILDING.md](BUILDING.md).
