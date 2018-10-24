REG=quay.io
ORG=integreatly
IMAGE=tutorial-web-app-operator
TAG=latest
KUBE_CMD=oc apply -f
DEPLOY_DIR=deploy
OUT_STATIC_DIR=tmp/_output


test-unit:
	go test -v -race -cover ./pkg/...

test: test-unit

template-copy:
	mkdir -p ${OUT_STATIC_DIR}/deploy/template
	cp ${DEPLOY_DIR}/template/tutorial-web-app.yml ${OUT_STATIC_DIR}/deploy/template

sdk-build:
	operator-sdk build ${REG}/${ORG}/${IMAGE}:${TAG}

build: template-copy sdk-build

push:
	docker push ${REG}/${ORG}/${IMAGE}:${TAG}

prepare:
	${KUBE_CMD} ${DEPLOY_DIR}/rbac.yaml
	${KUBE_CMD} ${DEPLOY_DIR}/sa.yaml
	${KUBE_CMD} ${DEPLOY_DIR}/crd.yaml
	${KUBE_CMD} ${DEPLOY_DIR}/cr.yaml

deploy:
	${KUBE_CMD} ${DEPLOY_DIR}/operator.yaml
