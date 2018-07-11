# -*- mode: ruby -*-
# vi:set ft=ruby sw=2 ts=2 sts=2:

NODE_COUNT = 2
POD_NETWORK_CIDR = "10.244.0.0/16"
KUBEADM_TOKEN = "d4daf2.baee52213f63b50b"

MASTER_IP = "192.168.1.11"
NODE_IP_PREFIX = "192.168.1."

Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"

  config.vm.provider "virtualbox" do |v|
    v.memory = 1536
    v.cpus = 2
  end

  config.vm.provision "install-docker", type: "shell", :path => "scripts/docker.sh"
  config.vm.provision "install-kubeadm", type: "shell", :path => "scripts/kubeadm.sh"
  config.vm.provision "install-tools", type: "shell", :path => "scripts/tools.sh"
  config.vm.provision "set-bridge-nf-traffic", type: "shell", :path => "scripts/set-bridge-nf-traffic.sh"
  config.vm.provision "setup-hosts", type: "shell", :path => "scripts/hosts.sh" do |s|
    s.args = ["eth1"]
  end
  config.vm.provision "shell", inline: <<-SHELL
     echo "swapoff -a" >> /root/.bash_profile
     source /root/.bash_profile
  SHELL

  config.vm.define "kube-master" do |node|
    node.vm.hostname = "kube-master"
    node.vm.network :private_network, ip: MASTER_IP

    node.vm.provision "shell", inline: <<-SHELL
        kubeadm init \
            --apiserver-advertise-address=#{MASTER_IP} \
            --pod-network-cidr=#{POD_NETWORK_CIDR} \
            --token #{KUBEADM_TOKEN} \
            --kubernetes-version="1.9.3"
    SHELL

    node.vm.provision "configure-kubeconfig", type: "shell", :path => "scripts/kubeconfig.sh"
    node.vm.provision "install-flannel", type: "shell", :path => "scripts/flannel.sh" do |s|
      s.args = ["eth1"]
    end
  end

  (1..NODE_COUNT).each do |i|
     config.vm.define "kube-node0#{i}" do |node|
        node.vm.hostname = "kube-node0#{i}"
        node.vm.network :private_network, ip: NODE_IP_PREFIX + "#{11 + i}"

        node.vm.provision "shell", inline: <<-SHELL
            kubeadm join --token #{KUBEADM_TOKEN} #{MASTER_IP}:6443 --discovery-token-unsafe-skip-ca-verification
        SHELL
     end
  end
end
