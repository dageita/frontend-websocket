# Vars
```
VERSION=${1.0}
NAME=${frontend-websocket}
IP=...
```

# Note
* Websocket thread port: 1234, you can customize it in code.
* MEC cluster frontend-websocket port: 32080, which can be customize in ${NAME}.yaml.
* Make sure image has presented in node which started ${NAME} service.
* ${NAME}-test.yaml is used to debug code which provide a golang enviroment includes vim.

# Make Docker Image
```
docker build -t ${NAME}:${VERSION} .
docker save -o /var/lib/rancher/k3s/agent/images/${NAME}-${VERSION}.tar ${NAME}:${VERSION}
systemctl restart k3s
```
<!-- Set cluster.server.ip (/etc/rancher/k3s/k3s.yaml) to local host public IP. -->

# Check YAML
```
image: docker.io/library/${NAME}:${VERSION}
nodeSelector: # You can add nodeSlector in ${NAME}.yaml to make sure this pod is started in the node with ${NAME}-${VERSION}.tar
kubectl apply -f ${NAME}.yaml
```

# Test Websocket Connection
Install Websocket test tool `wscat`(centos):
```
curl -sL https://rpm.nodesource.com/setup_10.x | sudo bash -
sudo yum install nodejs
npm install -g wscat
wscat -n -c 'ws://${IP}:32080/exec?namespace=default&pod=device-plugin-kvm-79b9x&container=device-plugin-kvm&command=sh'
```