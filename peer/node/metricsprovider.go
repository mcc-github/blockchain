/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	kitstatsd "github.com/go-kit/kit/metrics/statsd"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/common/metrics/prometheus"
	"github.com/mcc-github/blockchain/common/metrics/statsd"
	"github.com/mcc-github/blockchain/common/metrics/statsd/goruntime"
	"github.com/mcc-github/blockchain/core/comm"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)


type logFunc func(keyvals ...interface{}) error


func (l logFunc) Log(keyvals ...interface{}) error {
	return l(keyvals...)
}




func initializeMetrics() (provider metrics.Provider, shutdown func(), err error) {
	logger := flogging.MustGetLogger("metrics.provider")
	kitLogger := logFunc(func(keyvals ...interface{}) error {
		logger.Warn(keyvals...)
		return nil
	})

	providerType := viper.GetString("metrics.provider")
	switch providerType {
	case "statsd":
		network := viper.GetString("metrics.statsd.network")               
		address := viper.GetString("metrics.statsd.address")               
		writeInterval := viper.GetDuration("metrics.statsd.writeInterval") 
		prefix := viper.GetString("metrics.statsd.prefix")                 
		if prefix != "" {
			prefix = prefix + "."
		}

		c, err := net.Dial(network, address)
		if err != nil {
			return nil, nil, err
		}
		c.Close()

		ks := kitstatsd.New(prefix, kitLogger)
		statsdProvider := &statsd.Provider{Statsd: ks}
		goCollector := goruntime.NewCollector(statsdProvider)

		collectorTicker := time.NewTicker(writeInterval / 2)
		go goCollector.CollectAndPublish(collectorTicker.C)

		sendTicker := time.NewTicker(writeInterval)
		go ks.SendLoop(sendTicker.C, network, address)

		shutdown := func() {
			sendTicker.Stop()
			collectorTicker.Stop()
		}
		return statsdProvider, shutdown, nil

	case "prometheus":
		prometheusProvider := &prometheus.Provider{}

		address := viper.GetString("metrics.prometheus.listenAddress")                   
		handlerPath := viper.GetString("metrics.prometheus.handlerPath")                 
		tlsEnabled := viper.GetBool("metrics.prometheus.tls.enabled")                    
		certificate := viper.GetString("metrics.prometheus.tls.cert.file")               
		key := viper.GetString("metrics.prometheus.tls.key.file")                        
		clientCertRequired := viper.GetBool("metrics.prometheus.tls.clientAuthRequired") 
		caCerts := viper.GetStringSlice("metrics.prometheus.tls.clientRootCAs.files")    

		mux := http.NewServeMux()
		mux.Handle(handlerPath, prom.Handler())
		httpServer := &http.Server{
			Addr:         address,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 2 * time.Minute,
		}

		var tlsConfig *tls.Config
		if tlsEnabled {
			cert, err := tls.LoadX509KeyPair(certificate, key)
			if err != nil {
				return nil, nil, err
			}
			caCertPool := x509.NewCertPool()
			for _, caPath := range caCerts {
				caPem, err := ioutil.ReadFile(caPath)
				if err != nil {
					return nil, nil, err
				}
				caCertPool.AppendCertsFromPEM(caPem)
			}
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				CipherSuites: comm.DefaultTLSCipherSuites,
				ClientCAs:    caCertPool,
				NextProtos:   []string{"h2", "http/1.1"},
			}
			if clientCertRequired {
				tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			}
		}

		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, nil, err
		}
		if tlsConfig != nil {
			listener = tls.NewListener(listener, tlsConfig)
		}

		go httpServer.Serve(listener)
		shutdown := func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			httpServer.Shutdown(ctx)
		}
		return prometheusProvider, shutdown, nil

	default:
		if providerType != "disabled" {
			logger.Warnf("Unknown provider type: %s; metrics disabled", providerType)
		}

		disabledProvider := &disabled.Provider{}
		return disabledProvider, func() {}, nil
	}
}
