<p align="center">
  <img src=https://github.com/liqotech/liqo/workflows/Go/badge.svg>
  <img src=https://goreportcard.com/badge/github.com/liqotech/liqo>
  <img src=https://coveralls.io/repos/github/liqotech/liqo/badge.svg?branch=master>
  <img src=https://app.fossa.com/api/projects/git%2Bgithub.com%2Fliqotech%2Fliqo.svg?type=shield>
  <img src=https://img.shields.io/badge/slack-liqo.io-yellow>
  <img src=https://img.shields.io/twitter/url/https/twitter.com/liqo_io.svg?style=social&label=Follow%20%40liqo_io>

  <br />
  <br />

  <a href="https://github.com/liqotech/liqo">
    <img alt="Liqo Logo" src="https://doc.liqo.io/images/logo-liqo-blue.svg" height="80">
  </a>

  <h3 align="center">A <a href=https://github.com/kubernetes/kubernetes>Kubernetes</a> platform about dynamic and decentralized multi-cluster</h3>

  <br />

  <p align="center">
    <a href="https://doc.liqo.io/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://www.youtube.com/watch?v=tHCKGwnmuAA&t=1s&ab_channel=Liqo">View Demo</a>
    ·
    <a href="https://github.com/liqotech/liqo/issues/new?assignees=&labels=&template=bug_report.md&title=">Report Bug</a>
    ·
    <a href="https://github.com/liqotech/liqo/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=%5BFeature%5D">Request Feature</a>
  </p>

  
</p>

## About the project
Liqo is a platform to enable dynamic and decentralized resource sharing across Kubernetes clusters, either on-prem or managed. Liqo allows to run pods on a remote cluster seamlessly and without any modification of Kubernetes and the applications. With Liqo it is possible to extend the control plane of a Kubernetes cluster across the cluster's boundaries, making multi-cluster native and transparent: collapse an entire remote cluster to a virtual local node, by allowing workloads offloading and resource management compliant with the standard Kubernetes approach.

<br />
<p align="center">
  <img src=https://img.shields.io/badge/Google%20GKE-supported-green>
  <img src=https://img.shields.io/badge/Azure%20AKS-supported-green>
  <img src=https://img.shields.io/badge/Amazon%20EKS-in%20progress-orange>
</p>
<br />

<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#mainfeatures">Main Features</li>
    <li><a href="#quickstart">Quickstart</a></li>
    <li><a href="#installation">Installation</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#community">Community</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

## Main features

* **Decentralized governance**: peer-to-peer paradigm, without any centralized management entity.
* **Clusters discovery**: leverage on three different ways to discover (and peer to) other clusters:
  * Manual configuration: through a custom API representing other clusters
  * DNS: automatic discovery through DNS records
  * LAN: automatic discovery of neighboring clusters available in the same LAN.
* **Transparent offloading**: pods scheduled on the virtual node are offloaded to the remote cluster; they can be controlled by merely accessing the pod objects in the local one; the resources needed by the pods (services, endpoints, configmaps, etc.) are translated and replicated remotely. It allows inter-cluster pod-to-pod and pod-to-service communication.
* **Pod resilience**: the offloaded pods' lifecycle is controlled by a remote replicaset.
* **Inter-cluster networking**: the clusters inter-connection is implemented by a Wireguard tunnel, which ensure encryption and reliability.
* **CNI independence**: compliance with many CNIs (Calico, Cilium, Flannel, etc.) even in heterogeneous scenarios (the two clusters can have different CNIs).

## Quickstart

This quickstart lets you try Liqo in a playground environment built by two clusters in [KinD](https://kind.sigs.k8s.io/).

#### __Provision__ two KinD clusters.

```bash
source <(curl -L https://get.liqo.io/clusters.sh)
```

#### __Install__ Liqo on both clusters:

```bash
export KUBECONFIG=$KUBECONFIG_1
curl -L https://get.liqo.io | bash -s
export KUBECONFIG=$KUBECONFIG_2
curl -L https://get.liqo.io | bash -s
```

Wait that all containers are up and running. When a new virtual-kubelet pops out, a new node modeling the remote cluster is present and ready to receive pods. Check it out with:

```bash
kubectl get nodes
```

#### __Use__ the resources

Create a new namespace and label it to tell Liqo that the pods created in that namespace are suitable for offloading in the remote cluster.

```bash
kubectl create namespace liqo-demo
kubectl label namespace liqo-demo liqo.io/enabled=true
```

Deploy the [Google microservice Shop](https://github.com/GoogleCloudPlatform/microservices-demo) application. 

```bash
kubectl apply -f https://get.liqo.io/app.yaml -n liqo-demo
```

You can observe that:

* Your application is correctly working by exposing the application frontend port and later connecting with a browser to [localhost:8000](localhost:8000). To expose the pod port:
```bash
  kubectl port-forward -n liqo-demo service/frontend 8080:80
```
* Your application is transparently deployed across two different clusters:
```bash
  kubectl get pods -n liqo-demo -o wide  
``` 

## Installation

Liqo can be installed by using Helm. The configuration depends on the cluster type (K8s, K3s, managed) and the provider. Possible scenarios can be found [here](https://doc.liqo.io/user/scenarios/).

Once you identified your scenario, follow the instructions for the proper installation [environment]("https://doc.liqo.io/user/install/").

## Roadmap

Planned features for the next release (v0.3, expected mid-July, 2021) are the following:

* Support for deployments spanning across more than two clusters
* Support for a more balanced scheduling mechanism to distribute jobs across clusters
* Support for Amazon Elastic Kubernetes Service
* Support for more-granular permission control over remote cluster resources

## Contributing

All contributors are warmly welcome. If you want to become a new contributor, we are so happy!. Just, before doing it, read the repo's guidelines presented on our [documentation website](https://doc.liqo.io/contributing/).

## Community

To get involved with the Liqo community, join the slack [channel](https://join.slack.com/t/liqo-io/shared_invite/zt-h20212gg-g24YvN6MKiD9bacFeqZttQ).

|![notification](docs/images/readme/bell-outline-badged.svg) Community Meeting|
|------------------|
|Liqo holds weekly community meeting on Monday, 5.30pm UTC (6.30 CET, 9.30am PST). To join the community meeting, follow this [link](https://meet.google.com/dyr-ieso-smu). Convert to your timezone [here](https://www.thetimezoneconverter.com/?t=17:30&tz=UTC%20%28Universal%20Time%20Coordinated%29).|

## License

Distributed under the Apache-2.0 License. See [License](LICENSE) for more information.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fliqotech%2Fliqo.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fliqotech%2Fliqo?ref=badge_large)

<p align="center">
Liqo is a project kicked off at Polytechnic of Turin (Italy) and actively maintained with :heart: by this open-source community.
</p>
