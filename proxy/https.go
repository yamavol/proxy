package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

func (p *Proxy) httpsHandler(w http.ResponseWriter, r *http.Request) {
	// Create TCP connection with backend
	var backConn net.Conn
	var err error
	if p.isProxySet() {
		backConn, err = p.relayCONNECT(w, r)
		// TODO: if proxy refuses, return error
	} else {
		host := r.URL.Host
		backConn, err = net.Dial("tcp", host)
		// TODO: handle error
	}

	if err != nil {
		println(err.Error())
		return
	}

	w.Write([]byte("HTTP/1.0 200 Connection established\r\n\r\n"))

	// Hijack client-proxy connection to take ownership of the connection.
	conn, _, hjErr := w.(http.Hijacker).Hijack()
	if errors.Is(hjErr, http.ErrNotSupported) {
		println(hjErr.Error())
	}
	if hjErr != nil {
		return
	}
	targetTcp, targetOK := backConn.(*net.TCPConn)
	clientTcp, clientOK := conn.(*net.TCPConn)

	if targetOK && clientOK {
		go copyAndClose(targetTcp, clientTcp)
		go copyAndClose(clientTcp, targetTcp)
	} else {
		// TODO: fix?
		go copyAndClose(targetTcp, clientTcp)
		go copyAndClose(clientTcp, targetTcp)
	}
}

// If this proxy is behind another proxy, the CONNECT request has to be re-sent.
//
// http Client, Transport does not exposes net.Conn, nor an interface to
// send CONNECT request directly, so we have to do it manually.
// https://github.com/golang/go/issues/22554
//
// TLS handshake between proxy is currently unsupported.
func (p *Proxy) relayCONNECT(_ http.ResponseWriter, r *http.Request) (net.Conn, error) {

	proxyHost := p.Options.Proxy

	conn, err := net.Dial("tcp", proxyHost)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(conn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", r.URL.Host, proxyHost)
	br := bufio.NewReader(conn)
	res, err := http.ReadResponse(br, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server refused CONNECT")
	}
	if br.Buffered() > 0 {
		return nil, fmt.Errorf("unexpected body data in CONNECT response")
	}
	return conn, nil

}

// read from SRC, copy to DST, and do half-close on EOF
func copyAndClose(dst, src *net.TCPConn) {
	if _, err := io.Copy(dst, src); err != nil {
		println(err.Error())
	}
	dst.CloseWrite()
	src.CloseRead()
}
