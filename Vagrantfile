# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  config.vm.box = "joewest/livepeer-ubuntu1604"
  config.vm.box_version = "0.0.2"
  config.vm.hostname = "livepeer-ubuntu1604"

  # Livepeer RTMP and HTTP ports used by the guest VM.
  config.vm.network "forwarded_port", guest: 1935, host: 1935
  config.vm.network "forwarded_port", guest: 8935, host: 8935

  # NOTE: In some use cases it might be desirable to have the VM act as an
  # independent device on the network. In this case, disable the port
  # forwarding above and enable bridged networking by uncommenting the line below.
  # config.vm.network "public_network"

  config.ssh.username = "ubuntu"

  config.vm.provider "virtualbox" do |vb|
    # Customize the number of CPUs and amount of memory on the VM:
    vb.cpus = "2"
    vb.memory = "2048"
  end

end
