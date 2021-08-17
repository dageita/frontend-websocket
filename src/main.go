package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/rest"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var f *os.File

func EchoMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting Echo Message.")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// To solve fronend cross-platform problem.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	vars := r.URL.Query()
	fmt.Println("URL vars:", vars)
	namespace := vars.Get("namespace")
	pod := vars.Get("pod")
	container := vars.Get("container")
	command := vars.Get("command")

	cf, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println("Parsing URI from args.")
	var opts *ExecOptions = &ExecOptions{namespace, pod, container, []string{command}, true, true}
	fmt.Println("Frontend Websocket URL arguments:", opts)

	// var kubeconfig string = "/etc/rancher/k3s/k3s.yaml"
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	config, err := rest.InClusterConfig()

	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println("Preparing TLS Config.")
	tlsConfig, _ := rest.TLSConfigFor(config)
	rt := &WebsocketRoundTripper{
		Callback:    WebsocketCallback,
		TLSConfig:   tlsConfig,
		WebSocketCF: cf,
	}
	wrt, err := rest.HTTPWrappersForConfig(config, rt)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	fmt.Println("HTTP Wrap Config:", wrt)

	req, err := ExecRequest(config, opts)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println("Exec Request:", req)

	if _, err := wrt.RoundTrip(req); err != nil {
		fmt.Println("ERROR:", err)
	}
}

func ParseClusterInfo() (string, string) {
	clusterInfo, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return clusterInfo.Host, clusterInfo.BearerTokenFile
}

func main() {
	fmt.Println("Staring websocket Handler.")
	http.HandleFunc("/", EchoMessage)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		fmt.Println("ERROR:ListenAndServe:", err)
	}
}
