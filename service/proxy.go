package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Serve a reverse proxy for a given url
func Proxy(res http.ResponseWriter, req *http.Request) {
	// parse the url
	targetUrl, _ := url.Parse(Configs.Server.Uri)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// Prevent go transport layer from adding gzip itself
	req.Header.Add("Accept-Encoding", "identity")

	// Update the headers to allow for SSL redirection
	req.URL.Host = targetUrl.Host
	req.URL.Scheme = targetUrl.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = targetUrl.Host

	go processIp(req)

	//Adding Header example
	test := req.Header.Get("GoProx")
	res.Header().Set("GoProx", "0.1" + test)

	proxy.ServeHTTP(res, req)
}

func processIp(r *http.Request) {
	//ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	CacheClient.SaveVisit(r.RemoteAddr)
}

//func changeIPAddress(r *http.Request) {
//	r.RemoteAddr = "127.0.0.1"
//}