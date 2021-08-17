package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	TLSConfig   *tls.Config
	Callback    RoundTripCallback
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

		for {
			_, buf, err := cf.ReadMessage()
			if err != nil {
				return
			}

			fmt.Println("Frontend Reading Message:", buf, err)
			if err != nil {
				errChan <- err
				return
			}
			newbuf := buf[:len(buf)]
			newbuf = append(newbuf, byte(10))
			newbuf = append([]byte{0}, newbuf...)
			fmt.Println("Creating New Buffer:", newbuf)
			cacheBuff.Write(newbuf[1 : len(newbuf)-1])
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
			fmt.Println("Backend Reading Message:", buf, err)
			if err != nil {
				errChan <- err
				return
			}

			if len(buf) > 1 {
				fmt.Println("Tranform starting...")
				s := " "
				s = strings.Replace(string(buf[1:]), cacheBuff.String(), "", -1)
				fmt.Println("Transform Backend Message to:", s)
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
	content1, _ := json.Marshal(&dialer.TLSClientConfig)
	content2, _ := json.Marshal(&r)
	fmt.Println("RoundTrip websocket dialer.TLSClientConfig:", content1)
	fmt.Println("RoundTrip websocket dialer.http.Request:", content2)
	cb, resp, err := dialer.Dial(r.URL.String(), r.Header)
	fmt.Println("RoundTrip websocket Dial done.")
	if err != nil {
		return nil, err
	}
	defer cb.Close()
	return resp, wrt.Callback(wrt.WebSocketCF, cb)
}

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
		return nil, err
	}

	u.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec", opts.Namespace, opts.Pod)

	rawQuery := "stdout=true&&stdin=true&tty=true"
	for _, c := range opts.Command {
		rawQuery += "&command=" + c
	}

	if opts.Container != "" {
		rawQuery += "&container=" + opts.Container
	}

	u.RawQuery = rawQuery

	return &http.Request{
		Method: http.MethodGet,
		URL:    u,
	}, nil
}
