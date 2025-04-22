package entities

import (
	"errors"
	"strings"
)

// TLS related configuration - from command line arguments
type TlsConfiguration interface {
	TlsEnabled() bool           // If true, the service will be secured using TLS; both a certificate and key must be provided if enabled
	TlsCertificateFile() string // Path to .crt file containing certificate of server (in PEM format)
	TlsKeyFile() string         // Path to .key file containing private key of server (in PEM format)
}

type tlsConfiguration struct {
	tlsEnabled         bool
	tlsCertificateFile string
	tlsKeyFile         string
}

func NewTlsConfiguration(tlsEnabled bool, tlsCertificateFile string, tlsKeyFile string) (TlsConfiguration, error) {

	if tlsEnabled {
		if tlsCertificateFile == "" {
			return nil, errors.New("parameter 'tlsCertificateFile' can't be empty")
		}

		if !strings.HasSuffix(tlsCertificateFile, ".crt") {
			return nil, errors.New("parameter 'tlsCertificateFile' must have a .crt extension")
		}

		if tlsKeyFile == "" {
			return nil, errors.New("parameter 'tlsKeyFile' can't be empty")
		}
	}

	return &tlsConfiguration{
		tlsEnabled:         tlsEnabled,
		tlsCertificateFile: tlsCertificateFile,
		tlsKeyFile:         tlsKeyFile,
	}, nil
}

func (t *tlsConfiguration) TlsEnabled() bool {
	return t.tlsEnabled
}

func (t *tlsConfiguration) TlsCertificateFile() string {
	return t.tlsCertificateFile
}

func (t *tlsConfiguration) TlsKeyFile() string {
	return t.tlsKeyFile
}
