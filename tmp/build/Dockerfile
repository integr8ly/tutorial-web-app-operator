FROM registry.access.redhat.com/ubi8/ubi:8.4

# Add user and delete its password
RUN adduser tutorial-web-app-operator && passwd -d tutorial-web-app-operator
USER tutorial-web-app-operator

ADD tmp/_output/bin/tutorial-web-app-operator /usr/local/bin/tutorial-web-app-operator
ADD tmp/_output/deploy/template/tutorial-web-app.yml /home/tutorial-web-app-operator/deploy/template/tutorial-web-app.yml

