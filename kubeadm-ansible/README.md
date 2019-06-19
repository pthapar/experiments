# Description
This entails the steps to install a k8s cluster with floating IP. 
Note: The set up assumes that you already have a VM that you can access over ssh as root user.

# Pre-requisites
## Install ansible:
Centos: yum install ansible
Mac: brew install ansible

## Install etcd on one or more nodes:
yum install etcd -y

## Finding a floating IP that you can use(for AHV):
You can use a tool arp-scan to find out all used IP's in the local subnet. It prints IP's in an increasing order. You can find any missing 
IP's in the continous range and pick that as the floating IP. For example:
```
$ sudo arp-scan  -localnet
Interface: eth0, datalink type: EN10MB (Ethernet)
Starting arp-scan 1.9.2 with 4096 hosts (http://www.nta-monitor.com/tools-resources/security-tools/arp-scan/)
10.40.224.1	5c:16:c7:ff:ff:02 (5c:16:c7:08:1a:1c)	Big Switch Networks (DUP: 3)
10.40.224.2	5c:16:c7:0b:2b:f5	Big Switch Networks
10.40.224.3	5c:16:c7:0a:79:79	Big Switch Networks
10.40.224.4	5c:16:c7:08:1a:1c	Big Switch Networks
10.40.224.5	50:6b:8d:27:63:cc	(Unknown)
10.40.224.11	0c:c4:7a:c7:6e:12	Super Micro Computer, Inc.
```
As you can see from above output, IPs from 10.40.224.5 to 10.40.224.11 arenot used by any node on the localnet. So, pick any one as
the floating IP.
Note: Keep this floating IP handy. You would need it for changing the leader election service file below.


## Install an etcd cluster:
Leader election relies on an existing etcd cluster. Hence, you would have to deploy an etcd clsuter. We have an etcd service file in the repo, 
just change that with the node IP that you would like it to run on and and copy it or create a playbook.
Note: Remember the node IP that is advertized by etcd.


# Installation
We leverage ansible and kubeadm to install the k8s cluster

1. Update inventory file: Please update hosts.yml for your intended cluster
2. Remove previous installations if they exist: ansible-playbook -i hosts cleanup.yml
3. Install kube dependencies: ansible-playbook -i hosts kube-dependencies.yml
4. Bring up master node: ansible-playbook -i hosts master.yml
5. Join other nodes: ansible-playbook -i hosts cluster.yml
6. Replace the etcd node IP and floating(a.k.a cluster IP) in leader-election.service 
6. Start leader election: ansible-playbook -i hosts leader-election


# Testing your set up
At this point, you have a running k8s cluster with floating IP. You can try it out by deploying an app which uses node port:
ansible-playbook -i hosts apps.yml
