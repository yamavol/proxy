package main

import (
	"flag"
	"os"

	"github.com/yamavol/proxy/proxy"
)

func main() {

	pPort := flag.Int("port", 3128, "Port to listen on")
	pPac := flag.String("pac", "proxy.pac", "PAC file to serve")
	pProxy := flag.String("proxy", "", "Upstream proxy to use")
	pProxyEnv := flag.String("proxy-env", "", "Environment variable to use for upstream proxy")
	pProxyPac := flag.String("proxy-pac", "", "PAC file to use for upstream proxy")
	flag.Parse()

	if *pProxyEnv != "" {
		envProxy := os.Getenv(*pProxyEnv)
		if *pProxy == "" {
			*pProxy = envProxy
		} else {
			println("Warning: -proxy and -proxy-env are both set. Using -proxy")
		}
	}

	p := proxy.NewProxy()
	p.Options.Port = *pPort
	p.Options.PacFile = *pPac
	p.Options.Proxy = *pProxy
	p.Options.ProxyPac = *pProxyPac

	p.ProxyStart()
}
