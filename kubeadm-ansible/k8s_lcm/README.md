# Introduction
This repo contains k8s life cycle management ansible playbooks. Underneath, we leverage kubeadm.

## Heads up
1. Playbooks in this repo cover a subset of use cases. Workflows/playbooks have not undergone production grade testing yet
2. ATM, these playbooks only support HA k8s clusters w/o reverse proxy which means that that all certs, workers, control plane nodes have a SPOF which is the bootstrap master.
3. HA is only supported in k8s/kubeadm version 1.15.3 onwards.
4. For now, we are creating a new bootstrap token each time a new node join is attempted.
4. If you are using sherlock base image, it is recommended to run cleanup.yml playbook on any participating node before anything else

# Getting started
In order to get started, you would need a centos-7 VM/machine with ansible-2.8.4 installed on it. In order to install ansible-2.8.4 on centos, you may use:
```
$ sudo yum install -y ansible-2.8.4
```

Also, you would need to add public key of the machine where you will be running these playbooks to authorized_keys of all participating nodes. Please use usual tools like `cat/echo` to write the public key to `~/.ssh/authorized_keys` on target hosts

# Ansible variables & hosts
We leverage ansible variables including:
1. k8s_version: Version that you would like to create cluster with or upgrade cluster to.
2. encrypt_key: This is a hex string that is used by kubeadm to encrypt certificates uploaded to k8s secrets. This key is used by new control plane nodes
   to decrypt cert key pairs which is then used by new control plane nodes to generate new cert key/pairs for various control plane components.

Sample hosts.ini file has the following  groups:
1. bootstrapmaster: This is the first control plane node
2. new_control_plane_nodes: This is the set of nodes that you would like to add to the control plane
3. other_control_plane: This is the set of nodes that are control plane nodes but are not a bootstrapmaster control plane node.

# LCM Exercise
This section goes through a set of workflows that would provide a hands-on experience of performing basic LCM operations on k8s cluster.

* Pre-req
  * It is recommended to run a cleanup in case you are trying to use sherlock base image or a VM/machine that used to be a k8s node.
    ```
      ansible-playbook -i hosts.ini cleanup.yml
    ```

* Set up a 1 node k8s 1.14.4 cluster
  * Add a node to hosts.ini like show below:
    ```
    [all]
    nodeA ansible_host=<nodeA_IP_HERE> ansible_user=root

    [bootstrapmaster]
    nodeA

    [other_control_plane]

    [new_control_plane_nodes]
    ```
  * This step creates a 1 node 1.14.4 cluster.
    ```
    ansible-playbook -i hosts.ini create_one_node_cluster.yml --extra-vars "k8s_version=1.14.4" -vv
    ```

* Upgrade 1 node cluster to 1.15.3
  * Stage upgrade: This step downloads binaries and docker images on all target hosts.
    ```
    ansible-playbook -i hosts.ini stage.yml --extra-vars "k8s_version=1.15.3" -vv
    ```
  * Apply upgrade: This step uses the binaries and docker images from previous step and applies the upgrade.
    ```
    ansible-playbook -i hosts.ini apply.yml --extra-vars "k8s_version=1.15.3" -vv
    ```

* Add a control plane node to 1 node cluster
  * Add the new node to `new_control_plane_nodes` & `other_control_plane` like shown below for new node 10.45.27.59
    ```
    [all]
    nodeA ansible_host=<nodeA_IP_HERE> ansible_user=root
    nodeB ansible_host=<nodeB_IP_HERE> ansible_user=root

    [bootstrapmaster]
    nodeA

    [other_control_plane]
    nodeB
    
    [new_control_plane_nodes]
    nodeB
    ```
  * Join the new node as shown below:
    ```
    ansible-playbook -i hosts.ini add_control_plane.yml --extra-vars "k8s_version=1.15.3 encrypt_key=b2397eef1fb873b1ae6cde112795a852703b307909cc081bdb948ceddfb623ad" -vv
    ```
    Note: You can use custom value for encrypt_key as long as it meets kubeadm requirements.
 
 * Upgrade control plane nodes of an HA cluster to 1.15.4
  * Stage upgrade: This step downloads binaries and docker images on all target hosts.
    ```
    ansible-playbook -i hosts.ini stage.yml --extra-vars "k8s_version=1.15.4" -vv
    ```
  * Apply upgrade: This step uses the binaries and docker images from previous step and applies the upgrade.
    ```
    ansible-playbook -i hosts.ini apply.yml --extra-vars "k8s_version=1.15.4" -vv
    ```
