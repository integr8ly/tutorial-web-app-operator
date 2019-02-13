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
test change
