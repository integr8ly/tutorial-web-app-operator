# tutorial-web-app-operator

Openshift operator that handles integreatly tutorial-web-app deployments.

[![Build Status](https://travis-ci.org/integr8ly/tutorial-web-app-operator.svg?branch=master)](https://travis-ci.org/integr8ly/tutorial-web-app-operator)


|                 | Project Info  |
| --------------- | ------------- |
| License:        | Apache License, Version 2.0                      |
| IRC             | [#integreatly](https://webchat.freenode.net/?channels=integreatly) channel in the [freenode](http://freenode.net/) network. |


## Deploying

```sh
#create required resources
make cluster/prepare
#deploys the operator itself
make -B cluster/deploy
```

## Building

```sh
#builds image: quay.io/integreatly/tutorial-web-app-operator:latest
make image/build

#custom image params: registry.io/myusername/image-name:dev
make image/build REG=custom-registry.io ORG=myusername IMAGE=image-name TAG=dev
```

## Running tests

```
make test/unit
```

## Updating Web App Image Version
Update web app image version in the following files:
* Update [WebAppImage](pkg/handlers/webhandler.go) `WebAppImage = "quay.io/integreatly/tutorial-web-app:<version>"`
  * Image version that the deployment config gets reconciled to (**Must be updated as would override the version used by the template**)
* Update [tutorial-web-app.yml template](deploy/template/tutorial-web-app.yml) `image: quay.io/integreatly/tutorial-web-app:<version>`
  * Image version that gets deployed on initial processing of the template file

## Release

Update operator version files:

* Bump [operator version](version/version.go) 
```Version = "<version>"```
* Bump [makefile TAG](Makefile)
```TAG=<version>```
* Bump [operator image version](deploy/operator.yaml)
```image: quay.io/integreatly/tutorial-web-app-operator:v<version>```

Commit changes and open a pull request to master.

When the PR is merged, create a new Github release from master for v<version> e.g. v0.0.10.
This will trigger a ci build of the tag, and a new image pushed to quay.io with that tag.
