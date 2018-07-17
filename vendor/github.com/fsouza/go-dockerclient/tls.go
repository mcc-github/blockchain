






package docker

import (
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"time"
)

type tlsClientCon struct {
	*tls.Conn
	rawConn net.Conn
}

func (c *tlsClientCon) CloseWrite() error {
	
	
	if cwc, ok := c.rawConn.(interface {
		CloseWrite() error
	}); ok {
		return cwc.CloseWrite()
	}
	return nil
}

func tlsDialWithDialer(dialer *net.Dialer, network, addr string, config *tls.Config) (net.Conn, error) {
	
	
	
	timeout := dialer.Timeout

	if !dialer.Deadline.IsZero() {
		deadlineTimeout := dialer.Deadline.Sub(time.Now())
		if timeout == 0 || deadlineTimeout < timeout {
			timeout = deadlineTimeout
		}
	}

	var errChannel chan error

	if timeout != 0 {
		errChannel = make(chan error, 2)
		time.AfterFunc(timeout, func() {
			errChannel <- errors.New("")
		})
	}

	rawConn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	colonPos := strings.LastIndex(addr, ":")
	if colonPos == -1 {
		colonPos = len(addr)
	}
	hostname := addr[:colonPos]

	
	
	if config.ServerName == "" {
		
		config = copyTLSConfig(config)
		config.ServerName = hostname
	}

	conn := tls.Client(rawConn, config)

	if timeout == 0 {
		err = conn.Handshake()
	} else {
		go func() {
			errChannel <- conn.Handshake()
		}()

		err = <-errChannel
	}

	if err != nil {
		rawConn.Close()
		return nil, err
	}

	
	
	return &tlsClientCon{conn, rawConn}, nil
}


func copyTLSConfig(cfg *tls.Config) *tls.Config {
	return &tls.Config{
		Certificates:             cfg.Certificates,
		CipherSuites:             cfg.CipherSuites,
		ClientAuth:               cfg.ClientAuth,
		ClientCAs:                cfg.ClientCAs,
		ClientSessionCache:       cfg.ClientSessionCache,
		CurvePreferences:         cfg.CurvePreferences,
		InsecureSkipVerify:       cfg.InsecureSkipVerify,
		MaxVersion:               cfg.MaxVersion,
		MinVersion:               cfg.MinVersion,
		NameToCertificate:        cfg.NameToCertificate,
		NextProtos:               cfg.NextProtos,
		PreferServerCipherSuites: cfg.PreferServerCipherSuites,
		Rand:                   cfg.Rand,
		RootCAs:                cfg.RootCAs,
		ServerName:             cfg.ServerName,
		SessionTicketKey:       cfg.SessionTicketKey,
		SessionTicketsDisabled: cfg.SessionTicketsDisabled,
	}
}
