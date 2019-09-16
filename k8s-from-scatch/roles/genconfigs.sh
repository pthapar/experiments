# Execute this script on one of the nodes that has all certs to generate kubeconfigs for individual components
kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=${CERTS_DIR}/ca-k8s-apiserver.pem \
    --embed-certs=true \
    --server=https://127.0.0.1:6443 \
    --kubeconfig=${KUBECONFIG_DIR}/admin.kubeconfig

  kubectl config set-credentials admin \
    --client-certificate=${CERTS_DIR}/cert-admin.pem \
    --client-key=${CERTS_DIR}/cert-admin-key.pem \
    --embed-certs=true \
    --kubeconfig=${KUBECONFIG_DIR}/admin.kubeconfig

  kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=admin \
    --kubeconfig=${KUBECONFIG_DIR}/admin.kubeconfig


kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=${CERTS_DIR}/ca-k8s-apiserver.pem \
    --embed-certs=true \
    --server=https://127.0.0.1:6443 \
    --kubeconfig=${KUBECONFIG_DIR}/kube-controller-manager.kubeconfig

  kubectl config set-credentials system:kube-controller-manager \
    --client-certificate=${CERTS_DIR}/cert-k8s-controller-manager.pem \
    --client-key=${CERTS_DIR}/cert-k8s-controller-manager-key.pem \
    --embed-certs=true \
    --kubeconfig=${KUBECONFIG_DIR}/kube-controller-manager.kubeconfig

  kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=system:kube-controller-manager \
    --kubeconfig=${KUBECONFIG_DIR}/kube-controller-manager.kubeconfig

  kubectl config use-context default --kubeconfig=${KUBECONFIG_DIR}/kube-controller-manager.kubeconfig


  kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=${CERTS_DIR}/ca-k8s-apiserver.pem \
    --embed-certs=true \
    --server=https://127.0.0.1:6443 \
    --kubeconfig=${KUBECONFIG_DIR}/kube-scheduler.kubeconfig

  kubectl config set-credentials system:kube-scheduler \
    --client-certificate=${CERTS_DIR}/cert-k8s-scheduler.pem \
    --client-key=${CERTS_DIR}/cert-k8s-scheduler-key.pem \
    --embed-certs=true \
    --kubeconfig=${KUBECONFIG_DIR}/kube-scheduler.kubeconfig

  kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=system:kube-scheduler \
    --kubeconfig=${KUBECONFIG_DIR}/kube-scheduler.kubeconfig

  kubectl config use-context default --kubeconfig=${KUBECONFIG_DIR}/kube-scheduler.kubeconfig

## Below commands on all worker nodes:
# CERTS_DIR=~/k8s/certs/
# CONFIGS_DIR=~/k8s/configs/
# NODE_NAME=node[1,2,3]
kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=${CERTS_DIR}/ca-k8s-apiserver.pem \
    --embed-certs=true \
    --server=https://127.0.0.1:443 \
    --kubeconfig=${CONFIGS_DIR}/${NODE_NAME}.kubeconfig
kubectl config set-credentials system:node:${NODE_NAME} \
    --client-certificate=${CERTS_DIR}/cert-${NODE_NAME}.pem \
    --client-key=${CERTS_DIR}/cert-${NODE_NAME}-key.pem \
    --embed-certs=true \
    --kubeconfig=${CONFIGS_DIR}/${NODE_NAME}.kubeconfig

kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=system:node:${NODE_NAME} \
    --kubeconfig=${CONFIGS_DIR}/${NODE_NAME}.kubeconfig

kubectl config use-context default --kubeconfig=${CONFIGS_DIR}/${NODE_NAME}.kubeconfig

# Each node
kubectl config set-cluster kubernetes-the-hard-way \
    --certificate-authority=${CERTS_DIR}/ca-k8s-apiserver.pem \
    --embed-certs=true \
    --server=https://127.0.0.1:443 \
    --kubeconfig=${CONFIGS_DIR}/kube-proxy.kubeconfig

  kubectl config set-credentials system:kube-proxy \
    --client-certificate=${CERTS_DIR}/cert-k8s-proxy.pem \
    --client-key=${CERTS_DIR}/cert-k8s-proxy-key.pem \
    --embed-certs=true \
    --kubeconfig=${CONFIGS_DIR}/kube-proxy.kubeconfig

  kubectl config set-context default \
    --cluster=kubernetes-the-hard-way \
    --user=system:kube-proxy \
    --kubeconfig=${CONFIGS_DIR}/kube-proxy.kubeconfig

  kubectl config use-context default --kubeconfig=${CONFIGS_DIR}/kube-proxy.kubeconfig
