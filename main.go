package main

import (
	"encoding/json"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"log"
	"os"
	"strings"
)

func ToJsonString(v interface{}, pretty bool) string {
	var buf []byte
	if pretty {
		buf, _ = json.MarshalIndent(v, "", "\t")
	} else {
		buf, _ = json.Marshal(v)
	}
	return string(buf)
}

func main() {

	if len(os.Args) <= 1 {
		log.Println("[info] 需要指定nginx配置文件地址")
		return
	}
	var f = os.Args[1]

	p, err := parser.NewParser(f)
	if err != nil {
		log.Println("[parser.NewParser.Error]", err.Error())
		return
	}
	defer func() { _ = p.Close() }()
	conf, err := p.Parse()
	if err != nil {
		log.Println("[p.Parse.Error]", err.Error())
		return
	}
	var nHttp NHttp
	var nServer NServer
	var nLocation NLocation
	for _, directive := range conf.GetDirectives() {
		//log.Println("[pp]", directive.GetName(), directive.GetParameters())
		//continue
		if directive.GetName() == "http" {
			nHttp = parseHttp(directive.GetBlock().GetDirectives())
			log.Println("[nHttp]", ToJsonString(nHttp, true))
		}
		if directive.GetName() == "server" {
			nServer = parseServer(directive.GetBlock().GetDirectives())
			log.Println("[nServer]", ToJsonString(nServer, true))
		}
		if directive.GetName() == "location" {
			nLocation = parseLocation(directive.GetBlock().GetDirectives())
			log.Println("[nLocation]", ToJsonString(nLocation, true))
		}
	}
	log.Println("[file]", f)
	//log.Println("[nHttp]", ToJsonString(nHttp, true))
	//log.Println("[nServer]", ToJsonString(nServer, true))
	//log.Println("[nLocation]", ToJsonString(nLocation, true))

}

func parseServer(directives []gonginx.IDirective) NServer {
	var nServer NServer

	nServer.Listen = make([]string, 0)
	nServer.ServerName = make([]string, 0)
	nServer.Includes = make([]string, 0)
	nServer.ErrorPages = make([]string, 0)
	nServer.Locations = make([]NLocation, 0)
	nServer.SslCiphers = make([]string, 0)
	nServer.SslProtocols = make([]string, 0)

	for _, d := range directives {
		switch d.GetName() {
		case "listen":
			nServer.Listen = append(nServer.Listen, d.GetParameters()...)
		case "server_name":
			nServer.ServerName = append(nServer.ServerName, d.GetParameters()...)
		case "root":
			nServer.Root = d.GetParameters()[0]
		case "include":
			nServer.Includes = append(nServer.Includes, d.GetParameters()...)
		case "error_page":
			nServer.ErrorPages = append(nServer.ErrorPages, strings.Join(d.GetParameters(), " "))
		case "location":
			var tmpLocation = parseLocation(d.GetBlock().GetDirectives())
			tmpLocation.Path = strings.Join(d.GetParameters(), "")
			nServer.Locations = append(nServer.Locations, tmpLocation)
		case "access_log":
			nServer.AccessLog = d.GetParameters()[0]
		case "ssl_certificate":
			nServer.SslCertificate = d.GetParameters()[0]
		case "ssl_certificate_key":
			nServer.SslCertificateKey = d.GetParameters()[0]
		case "ssl_protocols":
			nServer.SslProtocols = append(nServer.SslProtocols, d.GetParameters()...)
		case "ssl_session_cache":
			nServer.SslSessionCache = d.GetParameters()[0]
		case "ssl_session_timeout":
			nServer.SslSessionTimeout = d.GetParameters()[0]
		case "ssl_ciphers":
			nServer.SslCiphers = append(nServer.SslCiphers, d.GetParameters()...)
		case "ssl_prefer_server_ciphers":
			nServer.SslPreferServerCiphers = d.GetParameters()[0]
		}
	}

	return nServer
}

func parseHttp(directives []gonginx.IDirective) NHttp {
	var nHttp NHttp

	nHttp.Servers = make([]NServer, 0)

	for _, d := range directives {
		switch d.GetName() {
		case "log_format":
			nHttp.LogFormat = d.GetParameters()
		case "access_log":
			nHttp.AccessLog = d.GetParameters()
		case "include":
			nHttp.Includes = d.GetParameters()
		case "sendfile":
			nHttp.Sendfile = d.GetParameters()
		case "server":
			nHttp.Servers = append(nHttp.Servers, parseServer(d.GetBlock().GetDirectives()))
		}
	}

	return nHttp
}

func parseLocation(directives []gonginx.IDirective) NLocation {
	var nLocation NLocation

	nLocation.ProxySetHeaders = make([][]string, 0)
	nLocation.ProxyCacheUseStale = make([]string, 0)
	nLocation.Indexs = make([]string, 0)
	nLocation.TryFiles = make([]string, 0)
	nLocation.Includes = make([]string, 0)
	nLocation.AddHeaders = make([][]string, 0)
	nLocation.Rewrites = make([]string, 0)

	for _, d := range directives {
		switch d.GetName() {
		case "index":
			nLocation.Indexs = append(nLocation.Indexs, d.GetParameters()...)
		case "try_files":
			nLocation.TryFiles = append(nLocation.TryFiles, strings.Join(d.GetParameters(), " "))
		case "root":
			nLocation.Root = d.GetParameters()[0]
		case "alias":
			nLocation.Alias = d.GetParameters()[0]
		case "add_header":
			nLocation.AddHeaders = append(nLocation.AddHeaders, d.GetParameters())
		case "autoindex":
			nLocation.AutoIndex = d.GetParameters()[0]
		case "rewrite":
			nLocation.Rewrites = append(nLocation.Rewrites, strings.Join(d.GetParameters(), " "))
		case "proxy_set_header":
			nLocation.ProxySetHeaders = append(nLocation.ProxySetHeaders, d.GetParameters())
		case "proxy_http_version":
			nLocation.ProxyHttpVersion = d.GetParameters()[0]
		case "proxy_pass":
			nLocation.IsProxy = true
			nLocation.ProxyPass = d.GetParameters()[0]
		case "proxy_cache_valid":
			nLocation.ProxyCacheValids = append(nLocation.ProxyCacheValids, d.GetParameters())
		case "proxy_cache_use_stale":
			nLocation.ProxyCacheUseStale = append(nLocation.ProxyCacheUseStale, d.GetParameters()...)
		case "fastcgi_pass":
			nLocation.IsFastCgi = true
			nLocation.FastCgiPass = d.GetParameters()[0]
		case "fastcgi_index":
			nLocation.FastCgiIndex = d.GetParameters()[0]
		case "fastcgi_param":
			nLocation.FastCgiParam = append(nLocation.FastCgiParam, strings.Join(d.GetParameters(), " "))
		case "include":
			nLocation.Includes = append(nLocation.Includes, d.GetParameters()...)
		}
	}

	return nLocation
}
