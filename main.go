package main

import (
	"encoding/json"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func HashString(data string) string {
	return fmt.Sprintf("%x", xxhash.Sum64String(data))
}

func main() {

	if len(os.Args) <= 1 {
		log.Println("[info] 需要指定nginx配置文件地址")
		return
	}
	var configFile = os.Args[1]

	directives, err := parseConfig(configFile)
	if err != nil {
		log.Println("[parser.NewParser.Error]", err.Error())
		return
	}

	var nHttp NHttp
	var nServer NServer
	var nLocation NLocation
	for _, directive := range directives {
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

func parseConfig(configFile string) (directives []gonginx.IDirective, err error) {
	p, err := parser.NewParser(configFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = p.Close() }()
	conf, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return conf.GetDirectives(), nil
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

	if len(nServer.Listen) > 0 || len(nServer.ServerName) > 0 {
		nServer.Id = HashString(fmt.Sprintf("%d", time.Now().UnixNano()))
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
			//  按照上下文推断出应该是server模块等数据
			var files = parseIncludeFiles(d.GetParameters()[0])
			nHttp.Servers = append(nHttp.Servers, parseIncludeConfig(files)...)
			nHttp.Includes = append(nHttp.Includes, d.GetParameters()...)
		case "sendfile":
			nHttp.Sendfile = d.GetParameters()
		case "server":
			nHttp.Servers = append(nHttp.Servers, parseServer(d.GetBlock().GetDirectives()))
		}
	}
	if len(nHttp.Servers) > 0 {
		nHttp.Id = HashString(fmt.Sprintf("%d", time.Now().UnixNano()))
	}

	return nHttp
}

func parseLocation(directives []gonginx.IDirective) NLocation {
	var nLocation NLocation

	nLocation.ProxySetHeaders = make([][]string, 0)
	nLocation.ProxyCacheUseStale = make([]string, 0)
	nLocation.ProxyCacheValid = make([][]string, 0)
	nLocation.Index = make([]string, 0)
	nLocation.TryFiles = make([]string, 0)
	nLocation.Includes = make([]string, 0)
	nLocation.AddHeaders = make([][]string, 0)
	nLocation.Rewrites = make([]string, 0)

	for _, d := range directives {
		switch d.GetName() {
		case "index":
			nLocation.Index = append(nLocation.Index, d.GetParameters()...)
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
			nLocation.ProxyCacheValid = append(nLocation.ProxyCacheValid, d.GetParameters())
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
		case "deny":
			nLocation.Deny = d.GetParameters()[0]
		}
	}
	nLocation.Id = HashString(fmt.Sprintf("%d", time.Now().UnixNano()))

	return nLocation
}

func parseIncludeConfig(configFiles []string) []NServer {
	var nServers []NServer
	for _, tmpFile := range configFiles {
		log.Println("[include]", tmpFile)
		parsedDirectives, err := parseConfig(tmpFile)
		if err != nil {
			continue
		}
		for _, directive := range parsedDirectives {
			if directive.GetName() == "server" {
				log.Println("[parse server]")
				tmpServer := parseServer(directive.GetBlock().GetDirectives())
				if len(tmpServer.Id) > 0 {
					nServers = append(nServers, tmpServer)
				}
			}
		}
	}
	return nServers
}

func parseIncludeFiles(pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return make([]string, 0)
	}
	return matches
}
