FROM alpine:3.6

RUN adduser -D tutorial-web-app-operator
USER tutorial-web-app-operator

ADD tmp/_output/bin/tutorial-web-app-operator /usr/local/bin/tutorial-web-app-operator
ADD tmp/_output/deploy/template/tutorial-web-app.yml /home/tutorial-web-app-operator/deploy/template/tutorial-web-app.yml