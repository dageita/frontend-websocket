Check Dockerfile.
```
docker built -t go-websocket:0.1 .
docker save -o /var/lib/rancher/k3s/agent/images/go-websocket-0.1.tar go-websocket:0.1
systemctl restart k3s
'''
Set cluster.server.ip (/etc/rancher/k3s/k3s.yaml) to local host public IP.
Write containers.image = docker.io/library/go-websocket:0.1
'''
kubectl apply -f go-websocket.yaml
'''