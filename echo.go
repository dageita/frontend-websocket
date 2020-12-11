package main

import (
    "log"
    "os"
    "github.com/gorilla/websocket"
    "net/http"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/rest"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var f *os.File

func EchoMessage(w http.ResponseWriter, r *http.Request) {
    log.Println("&&&&&&&&&&&&&&test2************\n\n")
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        // 解决跨域问题
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    vars := r.URL.Query()
    namespace := vars.Get("namespace")
    pod := vars.Get("pod")
    container := vars.Get("container")

    cf, err := upgrader.Upgrade(w, r, nil) // 实际应用时记得做错误处理
    if err != nil {
        log.Fatalln(err)
        return
    }
    log.Println("&&&&&&&&&&&&&&test3************\n\n")
    var opts *ExecOptions = &ExecOptions{namespace,pod,container,[]string {"/bin/bash"},true,true}

    var kubeconfig string = "/etc/rancher/k3s/k3s.yaml"
    // kubeconfig = replaceHomePath(kubeconfig)
    log.Println("opts:",opts)
    
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalln(err)
    }
    log.Println("&&&&&&&&&&&&&&test4************\n\n")
    // log.Println("config:",config)
    tlsConfig, err := rest.TLSConfigFor(config)
    rt := &WebsocketRoundTripper{
		Callback:  WebsocketCallback,
        TLSConfig: tlsConfig,
        WebSocketCF: cf,
	}
    wrt,err := rest.HTTPWrappersForConfig(config, rt)
	// wrt, err := ExecRoundTripper(config, WebsocketCallback)
	// if err != nil {
	// 	log.Fatalln(err)
    // }

    log.Println("wrt:",wrt)
    log.Println("wrt:",wrt)

	req, err := ExecRequest(config, opts)
	if err != nil {
		log.Fatalln(err)
    }
    log.Println("req.Header:",req.Header)

	if _, err := wrt.RoundTrip(req); err != nil {
		log.Fatalln(err)
	}
}

func main() {
    f,err := os.OpenFile("/var/log/go-websocket.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)
    log.Println("&&&&&&&&&&&&&&test************\n\n")
    // http.Handle("/", websocket.Handler(EchoMessage))
    http.HandleFunc("/",EchoMessage)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

