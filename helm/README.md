# FogFlow

![Version: 3.2.6](https://img.shields.io/badge/Version-0.0.15-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

A Helm chart for running the fiware Fog Flow on kubernetes.
Repository for providing [HELM Charts](https://helm.sh/) of Generic Enablers from the [FIWARE Catalogue](https://github.com/FIWARE/catalogue). The 
charts can be install into  [Kubernetes](https://kubernetes.io/) with [helm3](https://helm.sh/docs/).
FIWARE is a curated framework of open source platform components which can be assembled together and with other third-party platform components to
accelerate the development of Smart Solutions.

**Homepage:** <https://fogflow.readthedocs.io/en/latest/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Bin Cheng | Bin.Cheng@neclab.eu |
| Naveen Singh Bisht | Naveen.b@india.nec.com |
| Neeraj Srivastava  | Neeraj.Srivastava@india.nec.com |
| Vinod Rawat  | Vinod.Rawat@india.nec.com |
| Aniket babuta  | aniket.babuta@india.nec.com |


## Source Code

* <https://github.com/smartfog/fogflow>


## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` |initial number of target replications, can be different if autoscaling is enabled  |
| serviceAccount | object | `{"annotations":{},"create":true,"name":"fogflow-dns"}` |  |
| serviceAccount.name | string | `"fogflow-dns"` |  |

