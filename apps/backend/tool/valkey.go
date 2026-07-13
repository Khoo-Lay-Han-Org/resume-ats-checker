package tool

import (
	"crypto/tls"
	"log"

	"github.com/valkey-io/valkey-go"
	systemconfig "resuming/system-config"
)

var Valkey valkey.Client

func SetupValkey() error {
	var tlsConfig *tls.Config
	if systemconfig.ValkeyTLS {
		tlsConfig = &tls.Config{}
	}
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{systemconfig.ValkeyUri},
		Password:    systemconfig.ValkeyPassword,
		TLSConfig:   tlsConfig,
		// Upstash does not support client-side caching (no functionality is loss)
		DisableCache: true,
	})

	if err != nil {
		log.Println("Failed to connect to Valkey Server:", err)
		return err
	}
	Valkey = client
	return nil
}
