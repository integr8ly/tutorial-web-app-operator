SHELL=/bin/bash
REG=quay.io
ORG=integreatly
IMAGE=tutorial-web-app-operator
TAG=v0.0.58
KUBE_CMD=oc apply -f
DEPLOY_DIR=deploy
OUT_STATIC_DIR=tmp/_output
OUTPUT_BIN_NAME=./tmp/_output/bin/tutorial-web-app-operator
TARGET_BIN=cmd/tutorial-web-app-operator/main.go

.PHONY: setup/dep
setup/dep:
	@echo Installing golang dependencies
	@go get golang.org/x/sys/unix
	@go get golang.org/x/crypto/ssh/terminal
	@echo Installing dep
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	@echo setup complete

.PHONY: setup/travis
setup/travis:
	@echo Installing Operator SDK
	@curl -Lo operator-sdk https://github.com/operator-framework/operator-sdk/releases/download/v0.0.7/operator-sdk-v0.0.7-x86_64-linux-gnu && chmod +x operator-sdk && sudo mv operator-sdk /usr/local/bin/

.PHONY: code/compile
code/compile:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${OUTPUT_BIN_NAME} ${TARGET_BIN}

.PHONY: code/gen
code/gen:
	@operator-sdk generate k8s

.PHONY: code/check
code/check:
	@diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: code/fix
code/fix:
	@gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: image/build
image/build: code/compile
	@mkdir -p ${OUT_STATIC_DIR}/deploy/template
	@cp ${DEPLOY_DIR}/template/tutorial-web-app.yml ${OUT_STATIC_DIR}/deploy/template
	operator-sdk build ${REG}/${ORG}/${IMAGE}:${TAG}

.PHONY: image/build/push
image/build/push: image/build
	@docker push ${REG}/${ORG}/${IMAGE}:${TAG}

.PHONY: test/unit
test/unit:
	go test -v -race -cover ./pkg/...

.PHONY: cluster/prepare
cluster/prepare:
	${KUBE_CMD} ${DEPLOY_DIR}/rbac.yaml
	${KUBE_CMD} ${DEPLOY_DIR}/sa.yaml
	${KUBE_CMD} ${DEPLOY_DIR}/crd.yaml
	${KUBE_CMD} ${DEPLOY_DIR}/cr.yaml

.PHONY: cluster/deploy
cluster/deploy:
	${KUBE_CMD} ${DEPLOY_DIR}/operator.yaml
