package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/geekypanda/httpcache"
	"gopkg.in/unrolled/secure.v1"

	"github.com/altran-nl/krakend/config"
	"github.com/altran-nl/krakend/config/viper"
	"github.com/altran-nl/krakend/logging/gologging"
	"github.com/altran-nl/krakend/proxy"
	"github.com/altran-nl/krakend/router/mux"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	logLevel := flag.String("l", "ERROR", "Logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/etc/krakend/configuration.json", "Path to the configuration filename")
	flag.Parse()

	parser := viper.New()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, err := gologging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"127.0.0.1:8080", "example.com", "ssl.example.com"},
		SSLRedirect:           false,
		SSLHost:               "ssl.example.com",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})

	// routerFactory := mux.DefaultFactory(proxy.DefaultFactory(logger), logger)

	routerFactory := mux.NewFactory(mux.Config{
		Engine:       mux.DefaultEngine(),
		ProxyFactory: proxy.DefaultFactory(logger),
		Middlewares:  []mux.HandlerMiddleware{secureMiddleware},
		Logger:       logger,
		HandlerFactory: func(cfg *config.EndpointConfig, p proxy.Proxy) http.HandlerFunc {
			return httpcache.CacheFunc(mux.EndpointHandler(cfg, p), time.Minute)
		},
	})

	routerFactory.New().Run(serviceConfig)
}
