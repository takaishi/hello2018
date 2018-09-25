package main

import (
	"net/http"
	"log"
	"fmt"
	"io"
)

type HttpConnection struct {
	Request *http.Request
	Response *http.Response
}

type HttpConnectionChannel chan *HttpConnection

var connChannel = make(HttpConnectionChannel)

func PrintHTTP(conn *HttpConnection) {
	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
	for k, v := range conn.Request.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println("=========================================")
	fmt.Printf("HTTP/1.1 %v", conn.Response.Status)
	for k, v := range conn.Response.Header {
		fmt.Println(k, ":", v)
	}
	var body []byte
	conn.Request.Body.Read(body)
	fmt.Printf("Body: %v\n", body)
	fmt.Println("=========================================")
}

type Proxy struct {

}

func NewProxy() *Proxy { return &Proxy{} }

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}
	var body []byte
	r.Body.Read(body)
	fmt.Printf("%v\n", body)
	fmt.Printf("Body: %v\n", body)
	req, err = http.NewRequest(r.Method, fmt.Sprintf("http://127.0.0.1:8000%s", r.RequestURI), r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = client.Do(req)
	r.Body.Close()

	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	conn := &HttpConnection{r, resp}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}

	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
	resp.Body.Close()

	PrintHTTP(conn)

}

func main() {
	proxy := NewProxy()
	err := http.ListenAndServe(":12345", proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}