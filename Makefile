SHELL=/bin/bash
REG=quay.io
ORG=integreatly
IMAGE=tutorial-web-app-operator
TAG=v0.0.63
KUBE_CMD=oc apply -f
DEPLOY_DIR=deploy
OUT_STATIC_DIR=tmp/_output
OUTPUT_BIN_NAME=./tmp/_output/bin/tutorial-web-app-operator
TARGET_BIN=cmd/tutorial-web-app-operator/main.go
OPERATOR_IMAGE=$(REG)/$(ORG)/$(IMAGE):$(TAG)

.PHONY: setup/dep
setup/dep:
	@echo Installing golang dependencies
	@go get -u golang.org/x/lint/golint
	@echo setup complete

.PHONY: setup/travis
setup/travis:
	@echo Installing Operator SDK
	@curl -Lo operator-sdk https://github.com/operator-framework/operator-sdk/releases/download/v1.2.0/operator-sdk-v1.2.0-x86_64-linux-gnu && chmod +x operator-sdk && sudo mv operator-sdk /usr/local/bin/

.PHONY: code/compile
code/compile:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${OUTPUT_BIN_NAME} ${TARGET_BIN}

.PHONY: code/gen
code/gen:
	$(CONTROLLER_GEN) rbac:roleName=manager-role webhook paths="./..."
	@go generate ./...

.PHONY: code/check
code/check:
	@diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)
	golint ./pkg/... | (! egrep -vi 'comment on|or be unexported')
	go vet ./...

.PHONY: code/fix
code/fix:
	@gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: image/build
image/build: code/compile
	@mkdir -p ${OUT_STATIC_DIR}/deploy/template
	@cp ${DEPLOY_DIR}/template/tutorial-web-app.yml ${OUT_STATIC_DIR}/deploy/template
	echo "build image $(OPERATOR_IMAGE)"
	docker build . -t ${OPERATOR_IMAGE}

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

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif