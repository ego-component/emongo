package emongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Authentication struct {
	// TLS authentication
	TLS *TLSConfig
}

func (config *Authentication) ConfigureAuthentication(opts *options.ClientOptions) (err error) {
	if config.TLS != nil {
		if err = configureTLS(config.TLS, opts); err != nil {
			return err
		}
	}
	return nil
}

func configureTLS(config *TLSConfig, opts *options.ClientOptions) error {
	tlsConfig, err := config.LoadTLSConfig()
	if err != nil {
		return fmt.Errorf("error loading tls config: %w", err)
	}
	if tlsConfig != nil && tlsConfig.InsecureSkipVerify {
		opts.TLSConfig = tlsConfig
		return nil
	}
	return nil
}
