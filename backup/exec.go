package main

import (
	"log"
	"bytes"
	"crypto/tls"
	"fmt"
	// "io"
	"net/http"
	"net/url"
	// "os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/rest"
)

// ExecOptions describe a execute request args.
type ExecOptions struct {
	Namespace string
	Pod       string
	Container string
	Command   []string
	TTY       bool
	Stdin     bool
}

type RoundTripCallback func(cf *websocket.Conn, cb *websocket.Conn) error
// type RoundTripCallback func(c *websocket.Conn) error

type WebsocketRoundTripper struct {
	TLSConfig *tls.Config
	Callback  RoundTripCallback
	WebSocketCF *websocket.Conn
}

var cacheBuff bytes.Buffer

var protocols = []string{
	"v4.channel.k8s.io",
	"v3.channel.k8s.io",
	"v2.channel.k8s.io",
	"channel.k8s.io",
}

const (
	stdin = iota
	stdout
	stderr
)

func WebsocketCallback(cf *websocket.Conn, cb *websocket.Conn) error {
	errChan := make(chan error, 3)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		// buf := make([]byte, 1025)

		for {
			_, buf, err := cf.ReadMessage()
			if err != nil {
				return
			}

			fmt.Println("Frontend**read message:",buf,err)
			if err != nil {
				errChan <- err
				return
			}
			newbuf := buf[:len(buf)]
			newbuf = append(newbuf,byte(10))
			newbuf = append([]byte {0},newbuf...)
			fmt.Println("Create newbuf:",newbuf)
			cacheBuff.Write(newbuf[1:len(newbuf)-1])
			cacheBuff.Write([]byte{13, 10})
			if err := cb.WriteMessage(websocket.BinaryMessage, newbuf); err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			msgType, buf, err := cb.ReadMessage()
			fmt.Println("Backend**read message:",buf,err)
			if err != nil {
				errChan <- err
				return
			}

			if len(buf) > 1 {
				// var w io.Writer
				// switch buf[0] {
				// case stdout:
				// 	w = os.Stdout
				// case stderr:
				// 	w = os.Stderr
				// }

				// if w == nil {
				// 	continue
				// }

				log.Println("Tranform starting...")
				s := " "
				s = strings.Replace(string(buf[1:]), cacheBuff.String(), "", -1)
				log.Println("Transform backend message to:",s)
				if err = cf.WriteMessage(msgType, buf); err != nil {
					return
				}
			}
			cacheBuff.Reset()
		}
	}()

	wg.Wait()
	close(errChan)
	err := <-errChan
	return err
}


func (wrt *WebsocketRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	dialer := &websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: wrt.TLSConfig,
		Subprotocols:    protocols,
	}
	fmt.Println("RoundTrip websocket dialer.TLSClientConfig:",dialer.TLSClientConfig)
	cb, resp, err := dialer.Dial(r.URL.String(), r.Header)
	if err != nil {
		return nil, err
	}
	defer cb.Close()
	fmt.Println("RoundTrip 2")
	return resp, wrt.Callback(wrt.WebSocketCF,cb)
}


// func ExecRoundTripper(config *rest.Config, f RoundTripCallback) (http.RoundTripper, error) {
// 	tlsConfig, err := rest.TLSConfigFor(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rt := &WebsocketRoundTripper{
// 		Callback:  f,
// 		TLSConfig: tlsConfig,
// 	}

// 	return rest.HTTPWrappersForConfig(config, rt)
// }

func ExecRequest(config *rest.Config, opts *ExecOptions) (*http.Request, error) {
	u, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		return nil, fmt.Errorf("Unrecognised URL scheme in %v", u)
	}

	u.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec", opts.Namespace, opts.Pod)

	rawQuery := "stdout=true&tty=true"
	for _, c := range opts.Command {
		rawQuery += "&command=" + c
	}

	if opts.Container != "" {
		rawQuery += "&container=" + opts.Container
	}

	if opts.TTY {
		rawQuery += "&tty=true"
	}

	if opts.Stdin {
		rawQuery += "&stdin=true"
	}
	u.RawQuery = rawQuery

	return &http.Request{
		Method: http.MethodGet,
		URL:    u,
	}, nil
}
