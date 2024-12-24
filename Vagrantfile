# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV['VAGRANT_SERVER_URL'] = 'https://vagrant.elab.pro'
Vagrant.configure("2") do |config|
  config.vm.box = "debian/bullseye64"
  config.vm.box_check_update = false
  config.vm.box_version = "1.0.0"

  config.vm.synced_folder ".", "/home/vagrant/auth-service"

  config.vm.define "manager" do |manager|
    manager.vm.provider :virtualbox do |v|
      v.cpus = 2
      v.memory = 8192
    end

    manager.vm.hostname = "manager"
    manager.vm.network "private_network", type: "dhcp"

    manager.vm.network :forwarded_port, guest: 8080, host: 8080
    manager.vm.network :forwarded_port, guest: 6060, host: 6060
    manager.vm.network :forwarded_port, guest: 5432, host: 5432
    manager.vm.network :forwarded_port, guest: 6379, host: 6379
    manager.vm.network :forwarded_port, guest: 16686, host: 16686
    manager.vm.network :forwarded_port, guest: 9091, host: 9091
    manager.vm.network :forwarded_port, guest: 9092, host: 9092
    manager.vm.network :forwarded_port, guest: 9093, host: 9093
    manager.vm.network :forwarded_port, guest: 9200, host: 9200
    manager.vm.network :forwarded_port, guest: 5601, host: 5601
    manager.vm.network :forwarded_port, guest: 9090, host: 9090
    manager.vm.network :forwarded_port, guest: 9100, host: 9100
    manager.vm.network :forwarded_port, guest: 9187, host: 9187
    manager.vm.network :forwarded_port, guest: 3000, host: 3000

    manager.vm.provision "shell", inline: <<-SHELL
      sudo bash /home/vagrant/auth-service/scripts/install_docker.sh
    SHELL
  end
end
