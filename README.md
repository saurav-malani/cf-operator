# cf-operator
[![master](https://ci.flintstone.cf.cloud.ibm.com/api/v1/teams/containerization/pipelines/cf-operator/badge)](https://ci.flintstone.cf.cloud.ibm.com/teams/containerization/pipelines/cf-operator)
[![go report card](https://goreportcard.com/badge/code.cloudfoundry.org/cf-operator)](https://goreportcard.com/badge/code.cloudfoundry.org/cf-operator)


|Nightly build|[![nightly](https://ci.flintstone.cf.cloud.ibm.com/api/v1/teams/containerization/pipelines/cf-operator-nightly/badge)](https://ci.flintstone.cf.cloud.ibm.com/teams/containerization/pipelines/cf-operator-nightly)|
|-|-|

cf-operator will enable the deployment of BOSH Releases, especially Cloud Foundry, to Kubernetes.

It's implemented as a k8s operator, an active controller component which acts upon custom k8s resources.

* Incubation Proposal: [Containerizing Cloud Foundry](https://docs.google.com/document/d/1_IvFf-cCR4_Hxg-L7Z_R51EKhZfBqlprrs5NgC2iO2w/edit#heading=h.lybtsdyh8res)
* Slack: #cf-containers on <https://slack.cloudfoundry.org>
* Backlog: [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/2192232)

## Install

cf-operator is still missing core functionality.

## Development

### Start Operator Locally

    make up
    kubectl apply -f docs/examples/fissile_v1alpha1_boshdeployment_cr.yaml
    kubectl get boshdeployments.fissile.cloudfoundry.org
    kubectl get pods --watch

    # clean up
    kubectl delete configmap bosh-manifest
    kubectl delete configmap bosh-ops
    kubectl delete secret bosh-ops-secret
    kubectl delete boshdeployments.fissile.cloudfoundry.org example-boshdeployment
