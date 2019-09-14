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
