package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

//type HttpConnectionChannel chan *HttpConnection

//var connChannel = make(HttpConnectionChannel)

type Proxy struct {
}

func NewProxy() *Proxy { return &Proxy{} }

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}

	fmt.Printf("> %v %v\n", r.Method, r.RequestURI)
	for k, v := range r.Header {
		fmt.Printf("> %s: %s\n", k, v)
	}
	fmt.Printf("\n")
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(r.Body)
	body := bufBody.String()
	fmt.Printf("> Body: %v\n", body)
	fmt.Printf("\n")
	req, err = http.NewRequest(r.Method, fmt.Sprintf("http://127.0.0.1:8000%s", r.RequestURI), strings.NewReader(body))
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = client.Do(req)
	r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = &HttpConnection{r, resp}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}

	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

}

func main() {
	proxy := NewProxy()
	err := http.ListenAndServe(":12345", proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
