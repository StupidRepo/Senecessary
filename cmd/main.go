package main

import (
	"crypto/tls"
	"fmt"
	"github.com/StupidRepo/Senecessary/pkg/shared"
	"github.com/StupidRepo/Senecessary/pkg/web"
	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Get the environment variables
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is not set")
	}

	// Initialize a HTTP Client
	transport := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "localhost:8082",
		}),
		DisableKeepAlives: false,
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_AES_128_GCM_SHA256,
				tls.VersionTLS13,
				tls.VersionTLS10,
			},
			MinVersion:         tls.VersionTLS13,
			InsecureSkipVerify: true,
		},
	}

	shared.Client = http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	fmt.Println("Checking if we can login to Seneca with provided token!")

	if user := shared.Login(); user != nil {
		fmt.Printf("Successfully logged in as %s!\n", cases.Title(language.Und).String(user.DisplayName))
	} else {
		fmt.Println("Failed to login!")
		return
	}

	web.StartMux()
}
