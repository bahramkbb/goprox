package service

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Serve a reverse proxy for a given url
func  Proxy(res http.ResponseWriter, req *http.Request) {
	//Process ip visit and create history
	go processIpStats(req)

	//Process blocked ips
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	if PermanentBlackListIPs[ip] {
		errResponse(res, http.StatusTooManyRequests,
			"Too many requests in short time detected! You have been permanently blocked!")
		return
	} else if BlackListIPs[ip] {
		errResponse(res, http.StatusTooManyRequests,
			"Too many requests in short time detected! You have been temporary blocked!")
		return
	}

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

	//Adding GoProx Header
	test := req.Header.Get("GoProx")
	res.Header().Set("GoProx", "0.1" + test)

	proxy.ServeHTTP(res, req)
}

func processIpStats(r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	CacheClient.SaveVisit(ip)
}

func errResponse(res http.ResponseWriter, code int, body string) {
	res.WriteHeader(code)
	res.Write([]byte(body))
}