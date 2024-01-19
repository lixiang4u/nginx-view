package main

type NServer struct {
	Id                     string   `json:"id"`
	Listen                 []string `json:"listen"`
	ServerName             []string `json:"server_name"`
	SslCertificate         string   `json:"ssl_certificate"`
	SslCertificateKey      string   `json:"ssl_certificate_key"`
	SslSessionCache        string   `json:"ssl_session_cache"`
	SslSessionTimeout      string   `json:"ssl_session_timeout"`
	SslCiphers             []string `json:"ssl_ciphers"`
	SslPreferServerCiphers string   `json:"ssl_prefer_server_ciphers"`
	SslProtocols           []string `json:"ssl_protocols"`

	Root       string      `json:"root"`
	Includes   []string    `json:"includes"`
	ErrorPages []string    `json:"error_pages"`
	Locations  []NLocation `json:"locations"`
	AccessLog  string      `json:"access_log"`
}

type NLocation struct {
	Id                 string     `json:"id"`
	IsProxy            bool       `json:"is_proxy"`
	IsFastCgi          bool       `json:"is_fast_cgi"`
	Path               string     `json:"path"`
	Root               string     `json:"root"`
	Alias              string     `json:"alias"`
	AddHeaders         [][]string `json:"add_headers"`
	AutoIndex          string     `json:"auto_index"`
	Rewrites           []string   `json:"rewrites"`
	ProxyPass          string     `json:"proxy_pass"`
	ProxySetHeaders    [][]string `json:"proxy_set_headers"`
	ProxyHttpVersion   string     `json:"proxy_http_version"`
	ProxyCacheValid    [][]string `json:"proxy_cache_valid"`
	ProxyCacheUseStale []string   `json:"proxy_cache_use_stale"`
	Index              []string   `json:"index"`
	TryFiles           []string   `json:"try_files"`
	FastCgiPass        string     `json:"fast_cgi_pass"`
	FastCgiIndex       string     `json:"fast_cgi_index"`
	FastCgiParam       []string   `json:"fast_cgi_param"`
	Includes           []string   `json:"includes"`
	Deny               string     `json:"deny"`
}

type NHttp struct {
	Id        string    `json:"id"`
	LogFormat []string  `json:"log_format"`
	AccessLog []string  `json:"access_log"`
	Includes  []string  `json:"includes"`
	Sendfile  []string  `json:"sendfile"`
	Servers   []NServer `json:"servers"`
}
