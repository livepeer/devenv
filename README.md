# Dev Environment

This repository contains everything you need to set up a Livepeer local
dev environment.

There is a virtual machine that is pre-built with key dependencies and
includes a wizard (lpdev) that will let you:

- Run a local Geth node and a local Ethereum network
- Deploy the Livepeer protocol to the local Ethereum network
- Configure and run a Livepeer broadcaster node
- Configure and run a Livepeer transcoder node

## Installation

Install Vagrant 2.0.x: https://www.vagrantup.com/

Install VirtualBox 5.0.x or 5.1.x: https://www.virtualbox.org/

Clone the repo and run this command:

```
host-machine $ git clone https://github.com/livepeer/testenv.git
host-machine $ cd testenv
host-machine $ vagrant up
```

### Host Machine Repos and `$LPSRC`

Vagrant will mount a host machine directory as the location for this
project's source repos in the virtual machine as `$HOME/src`.

You can provide a specific directory to Vagrant when bringing the
virtual machine "up" via the environment variable `$LPSRC`.

```
host-machine $ LPSRC=$HOME/your/local/repos/dir vagrant up
```

Or during a reload:

```
host-machine $ LPSRC=$HOME/your/local/repos/dir vagrant reload
```

If this is not provided and `~/src` does not exist, the default value of
`..` is used and should work in most cases.

## Usage

**Protip:** Virtual machines are cheap. Don’t feel tied to this
directory. Make new directories and copy the Vagrantfile to start fresh
at anytime. Or `vagrant destroy` to blow away the vm in this directory.

### Enter the virtual machine

```
host-machine $ vagrant ssh -- -A
```

### Run the lpdev wizard

```
virtual-machine $ lpdev

+----------------------------------------------------+
| Welcome to the Livepeer local dev environment tool |
|                                                    |
+----------------------------------------------------+

== Current Status ==

Geth miner is running: true (19943)
Geth accounts:
  ca4a8b268e1fb4d7105d95d0c90aa5b2b3d6b4ac (miner)
  3206f44d47f9059366a80b6073bc0132cb77b192
  f938d2bc074b36d480cdee4817c907ec8dc6ef4b
  7a50d0f6f8c07568233bfe62538d70efec4df261

Protocol is built: true (current branch: develop)
Protocol deployed to: 0x4d6d6762160436facc3571f1751ce2cd66a7f553

Broadcaster node is running: true (21101)
Transcoder node is running: true (21182)

--

What would you like to do?
1) Display status                      5) Start & set up transcoder node
2) Set up & start Geth local network   6) Destroy current environmentonment
3) Deploy/overwrite protocol contracts 7) Exit
4) Start & set up broadcaster node
```

_The Current Status values above are specific to your environment. Yours
will be different._

## Additional Details

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

The following file and directories are part of the `ubuntu` user’s $HOME
in the virtual machine:

### ~/.ethereum
This is the $gethDir for this virtual machine. It is created by Geth.

### ~/.lpdata
This directory is created and used by Livepeer nodes that are run using
the wizard on this virtual machine.  Each node has its own Geth Account
and the first 10 characters of the Account is used to represent the node
as subdirectories within `~/.lpdata`.

### ~/.lpdata
This directory is created and used by Livepeer nodes that are run using

### ~/.lpdev_cmds.sh
This file is where the Livepeer local dev environment wizard lives. It
is synced from the `dot_lpdev_cmds.sh` file in this repo during the
Vagrant provisioning step and sourced `source $HOME/.lpdev_cmds.sh` in
any interactive shell.

You can add new commands or edit existing ones by modifying the virtual
machine file and re-source'ing:

```
virtual-machine $ source $HOME/.lpdev_cmds.sh
```

Or modifying the host machine copy of `dot_lpdev_cmds.sh` and
re-provisioning:

```
host-machine $ vagrant provision
```

### ~/go
This is the $GOROOT for this virtual machine.

### ~/src
This is the synced folder on the virtual machine for working copies of
various project repositories checked out and stored on the host machine.
_Changes in the virtual machine are made to the host machine and
vice versa._

The three Repos listed above are technically part of the base virtual
machine image, but they are almost never accessible unless something
seriously weird happens. The `~/src` directory is replaced with any
repositories in the directory on the host machine as provided by
environment var `$LPSRC` or the defaults (`~/src`, `..`).

### Vagrant Commands
These commands are useful for operating a virtual machine run by
Vagrant. Run them in the directory that contains your Vagrantfile.

- If you want to init or boot it: `host-machine $ vagrant up`
- If you want to stop it: `host-machine $ vagrant halt`
- If you want to pause it: `host-machine $ vagrant suspend`
- If you want to blow it away and start over: `host-machine $ vagrant destroy`
- If you want to reboot the machine: `host-machine $ vagrant reload`
_Reload is useful if you change the Vagrantfile configuration and want
to apply those changes._

### Building a custom virtual machine

Step-by-step instructions for building your own base virtual machine
starting with a blank Ubuntu Vagrant box are in
[BUILDING.md](BUILDING.md).
