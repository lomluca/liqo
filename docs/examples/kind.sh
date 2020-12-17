#!/bin/bash

function setup_arch_and_os(){
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
    *) echo "Error architecture '${ARCH}' unknown"; exit 1 ;;
  esac

  OS=$(uname |tr '[:upper:]' '[:lower:]')
  case "$OS" in
    # Minimalist GNU for Windows
    "mingw"*) OS='windows'; return ;;
  esac

  # list is available for kind at https://github.com/kubernetes-sigs/kind/releases
  # kubectl supported architecture list is a superset of the Kind one. No need to further compatibility check.
  local supported="darwin-amd64\n\nlinux-amd64\nlinux-arm64\nlinux-ppc64le\nwindows-amd64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "Error: No version of kind for '${OS}-${ARCH}'"
    return 1
  fi

}

setup_arch_and_os

CLUSTER_NAME=cluster
CLUSTER_NAME_1=${CLUSTER_NAME}1
CLUSTER_NAME_2=${CLUSTER_NAME}2
KIND_VERSION="v0.9.0"
KUBECTL_DOWNLOAD=false

echo "Downloading Kind ${KIND_VERSION}"
TMPDIR=$(mktemp -d -t liqo-install.XXXXXXXXXX)
BINDIR="${TMPDIR}/bin"
mkdir -p "${BINDIR}"

if ! command -v docker &> /dev/null;
then
	echo "MISSING REQUIREMENT: docker engine could not be found on your system. Please install docker engine to continue: https://docs.docker.com/get-docker/"
	return 1
fi

if ! command -v kubectl &> /dev/null
then
    echo "WARNING: kubectl could not be found. Downloading and installing it locally..."
    if ! curl --fail -Lo "${BINDIR}"/kubectl "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/${OS}/${ARCH}/kubectl"; then
        echo "Error: Unable to download kubectl for '${OS}-${ARCH}'"
    	return 1
    fi
    chmod +x "${BINDIR}"/kubectl
    export PATH=${PATH}:${BINDIR}
	KUBECTL_DOWNLOAD=true
fi

curl -Lo "${BINDIR}"/kind https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-${OS}-${ARCH}
chmod +x "${BINDIR}"/kind
KIND="${BINDIR}/kind"

echo "Cleaning: Deleting old clusters"
${KIND} delete cluster --name $CLUSTER_NAME_1
${KIND} delete cluster --name $CLUSTER_NAME_2

cat << EOF > liqo-cluster-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  serviceSubnet: "10.90.0.0/12"
  podSubnet: "10.200.0.0/16"
nodes:
  - role: control-plane
    image: kindest/node:v1.19.1
  - role: worker
    image: kindest/node:v1.19.1
EOF

echo "Creating cluster $CLUSTER_NAME_1..."
${KIND} -v 3 create cluster --name $CLUSTER_NAME_1 --kubeconfig liqo_kubeconf_1 --config liqo-cluster-config.yaml --wait 2m

echo "Create cluster $CLUSTER_NAME_2..."
${KIND} -v 3 create cluster --name $CLUSTER_NAME_2 --kubeconfig liqo_kubeconf_2 --config liqo-cluster-config.yaml --wait 2m

echo ----------------
#Environment variables for E2E testing
CURRENT_DIRECTORY=$(pwd)
#Environment variables for E2E testing
export KUBECONFIG_1=${CURRENT_DIRECTORY}/liqo_kubeconf_1
export KUBECONFIG_2=${CURRENT_DIRECTORY}/liqo_kubeconf_2
export NAMESPACE=liqo
echo "Exported Variables:"
echo "- KUBECONFIG_1=${KUBECONFIG_1}"
echo "- KUBECONFIG_2=${KUBECONFIG_2}"
echo "- NAMESPACE=liqo"
# shellcheck disable=SC2016
echo "If you want to select $CLUSTER_NAME_1, you should simply type:" 'export KUBECONFIG=$KUBECONFIG_1'
# shellcheck disable=SC2016
echo "If you want to select $CLUSTER_NAME_2, you should simply type:" 'export KUBECONFIG=$KUBECONFIG_2'
if [ "$KUBECTL_DOWNLOAD" == "true" ]; then
	echo "kubectl is now installed in ${BINDIR}/kubectl and has been added to your PATH. To make it available without explicitly setting the PATH variable"
	echo "You can copy it to a system-wide location such as /usr/local/bin by typing:"
	echo "sudo cp ${BINDIR}/kubectl /usr/local/bin"
fi
