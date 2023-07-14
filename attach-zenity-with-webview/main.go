package main

import (
	"app/localhttpd"
	"net/http"
	"net/url"

	"github.com/jchv/go-webview2"
	"github.com/ncruces/zenity"
)

type Handler struct {
	webview webview2.WebView
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/page":
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`
			<button id="button">dialog</button>
			<script>
				const button = document.getElementById("button");
				button.onclick = () => fetch("/dialog");
			</script>
		`))

	case "/dialog":
		go func() {
			options := []zenity.Option{
				zenity.Attach(uintptr(h.webview.Window())),
			}

			if path, err := zenity.SelectFile(options...); err != nil {
				println("canceled")

			} else {
				println("selected:", path)
			}
		}()
	}
}

func main() {
	handler := new(Handler)

	httpd, err := localhttpd.Launch(handler, 12345)
	if err != nil {
		panic(err)
	}

	pageUrl, err := url.JoinPath(httpd.Url(), "page")
	if err != nil {
		panic(err)
	}

	webview := webview2.New(true)
	handler.webview = webview

	webview.Navigate(pageUrl)
	webview.Run()
}
