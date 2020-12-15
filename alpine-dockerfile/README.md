Check Dockerfile.
```
docker built -t go-websocket:0.1 .
docker save -o /var/lib/rancher/k3s/agent/images/go-websocket-0.1.tar go-websocket:0.1
systemctl restart k3s
```
Set cluster.server.ip (/etc/rancher/k3s/k3s.yaml) to local host public IP.

Write containers.image = docker.io/library/go-websocket:0.1
You can add nodeSlector in go-websocket.yaml to make sure this pod is started in the node with go-websocket-0.1.tar
```
kubectl apply -f go-websocket.yaml
```
Test websocket tool: wscat
```
wscat -n -c "ws://10.121.12.175:32080/exec?namespace=default&pod=nginx-deployment-66b6c48dd5-4lmp5&container=nginx&command=bash"
```