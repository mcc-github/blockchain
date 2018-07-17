/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/mcc-github/blockchain/common/flogging"
	logging "github.com/op/go-logging"
)

const (
	pkgLogID    = "orderer/consensus/kafka"
	saramaLogID = pkgLogID + "/sarama"
)

var logger *logging.Logger

var saramaLogger eventLogger


func init() {
	logger = flogging.MustGetLogger(pkgLogID)
}


func init() {
	loggingProvider := flogging.MustGetLogger(saramaLogID)
	loggingProvider.ExtraCalldepth = 3
	saramaEventLogger := &saramaLoggerImpl{
		logger: loggingProvider,
		eventListenerSupport: &eventListenerSupport{
			listeners: make(map[string][]chan string),
		},
	}
	sarama.Logger = saramaEventLogger
	saramaLogger = saramaEventLogger
}


func init() {
	listener := saramaLogger.NewListener("insufficient data to decode packet")
	go func() {
		for {
			select {
			case <-listener:
				logger.Critical("Unable to decode a Kafka packet. Usually, this " +
					"indicates that the Kafka.Version specified in the orderer " +
					"configuration is incorrectly set to a version which is newer than " +
					"the actual Kafka broker version.")
			}
		}
	}()
}




type eventLogger interface {
	sarama.StdLogger
	NewListener(substr string) <-chan string
	RemoveListener(substr string, listener <-chan string)
}

type saramaLoggerImpl struct {
	logger               *logging.Logger
	eventListenerSupport *eventListenerSupport
}

func (l saramaLoggerImpl) Print(args ...interface{}) {
	l.print(fmt.Sprint(args...))
}

func (l saramaLoggerImpl) Printf(format string, args ...interface{}) {
	l.print(fmt.Sprintf(format, args...))
}

func (l saramaLoggerImpl) Println(args ...interface{}) {
	l.print(fmt.Sprintln(args...))
}

func (l saramaLoggerImpl) print(message string) {
	l.eventListenerSupport.fire(message)
	l.logger.Debug(message)
}


const listenerChanSize = 100

func (l saramaLoggerImpl) NewListener(substr string) <-chan string {
	listener := make(chan string, listenerChanSize)
	l.eventListenerSupport.addListener(substr, listener)
	return listener
}

func (l saramaLoggerImpl) RemoveListener(substr string, listener <-chan string) {
	l.eventListenerSupport.removeListener(substr, listener)
}



type eventListenerSupport struct {
	sync.Mutex
	listeners map[string][]chan string
}


func (b *eventListenerSupport) addListener(substr string, listener chan string) {
	b.Lock()
	defer b.Unlock()
	if listeners, ok := b.listeners[substr]; ok {
		b.listeners[substr] = append(listeners, listener)
	} else {
		b.listeners[substr] = []chan string{listener}
	}
}



func (b *eventListenerSupport) fire(message string) {
	b.Lock()
	defer b.Unlock()
	for substr, listeners := range b.listeners {
		if strings.Contains(message, substr) {
			for _, listener := range listeners {
				listener <- message
			}
		}
	}
}


func (b *eventListenerSupport) removeListener(substr string, listener <-chan string) {
	b.Lock()
	defer b.Unlock()
	if listeners, ok := b.listeners[substr]; ok {
		for i, l := range listeners {
			if l == listener {
				copy(listeners[i:], listeners[i+1:])
				listeners[len(listeners)-1] = nil
				b.listeners[substr] = listeners[:len(listeners)-1]
			}
		}
	}
}
