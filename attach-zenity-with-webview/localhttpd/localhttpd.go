package localhttpd

import (
	"fmt"
	"net"
	"net/http"
)

type LocalHttpd interface {
	Url() string
	Wait()
	Close()
}

type _localttpd struct {
	url     string
	closeCh chan struct{}
	doneCh  chan struct{}
}

func Launch(handler http.Handler, port int) (LocalHttpd, error) {
	ret := new(_localttpd)
	ret.closeCh = make(chan struct{})
	ret.doneCh = make(chan struct{})

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, err
	}

	ret.url = "http://" + listener.Addr().String()

	errorCh := make(chan error)

	go func() {
		errorCh <- http.Serve(listener, handler)
	}()
	go func() {
		select {
		case <-errorCh:
		case <-ret.closeCh:
		}
		listener.Close()
		close(ret.doneCh)
	}()

	return ret, nil
}

func (h *_localttpd) Url() string {
	return h.url
}

func (h *_localttpd) Wait() {
	<-h.doneCh
}

func (h *_localttpd) Close() {
	close(h.closeCh)
}
