/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging/fabenc"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type Config struct {
	
	
	
	
	
	
	
	Format string

	
	
	
	
	LogSpec string

	
	
	
	Writer io.Writer
}




type Logging struct {
	*LoggerLevels

	mutex          sync.RWMutex
	encoding       Encoding
	encoderConfig  zapcore.EncoderConfig
	multiFormatter *fabenc.MultiFormatter
	writer         zapcore.WriteSyncer
	observer       Observer
}



func New(c Config) (*Logging, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.NameKey = "name"

	l := &Logging{
		LoggerLevels: &LoggerLevels{
			defaultLevel: defaultLevel,
		},
		encoderConfig:  encoderConfig,
		multiFormatter: fabenc.NewMultiFormatter(),
	}

	err := l.Apply(c)
	if err != nil {
		return nil, err
	}
	return l, nil
}


func (l *Logging) Apply(c Config) error {
	err := l.SetFormat(c.Format)
	if err != nil {
		return err
	}

	if c.LogSpec == "" {
		c.LogSpec = os.Getenv("FABRIC_LOGGING_SPEC")
	}
	if c.LogSpec == "" {
		c.LogSpec = defaultLevel.String()
	}

	err = l.LoggerLevels.ActivateSpec(c.LogSpec)
	if err != nil {
		return err
	}

	if c.Writer == nil {
		c.Writer = os.Stderr
	}
	l.SetWriter(c.Writer)

	return nil
}





func (l *Logging) SetFormat(format string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if format == "" {
		format = defaultFormat
	}

	if format == "json" {
		l.encoding = JSON
		return nil
	}

	if format == "logfmt" {
		l.encoding = LOGFMT
		return nil
	}

	formatters, err := fabenc.ParseFormat(format)
	if err != nil {
		return err
	}
	l.multiFormatter.SetFormatters(formatters)
	l.encoding = CONSOLE

	return nil
}




func (l *Logging) SetWriter(w io.Writer) io.Writer {
	var sw zapcore.WriteSyncer
	switch t := w.(type) {
	case *os.File:
		sw = zapcore.Lock(t)
	case zapcore.WriteSyncer:
		sw = t
	default:
		sw = zapcore.AddSync(w)
	}

	l.mutex.Lock()
	ow := l.writer
	l.writer = sw
	l.mutex.Unlock()

	return ow
}



func (l *Logging) SetObserver(observer Observer) Observer {
	l.mutex.Lock()
	so := l.observer
	l.observer = observer
	l.mutex.Unlock()

	return so
}




func (l *Logging) Write(b []byte) (int, error) {
	l.mutex.RLock()
	w := l.writer
	l.mutex.RUnlock()

	return w.Write(b)
}



func (l *Logging) Sync() error {
	l.mutex.RLock()
	w := l.writer
	l.mutex.RUnlock()

	return w.Sync()
}



func (l *Logging) Encoding() Encoding {
	l.mutex.RLock()
	e := l.encoding
	l.mutex.RUnlock()
	return e
}



func (l *Logging) ZapLogger(name string) *zap.Logger {
	if !isValidLoggerName(name) {
		panic(fmt.Sprintf("invalid logger name: %s", name))
	}

	l.mutex.RLock()
	core := &Core{
		LevelEnabler: l.LoggerLevels,
		Levels:       l.LoggerLevels,
		Encoders: map[Encoding]zapcore.Encoder{
			JSON:    zapcore.NewJSONEncoder(l.encoderConfig),
			CONSOLE: fabenc.NewFormatEncoder(l.multiFormatter),
			LOGFMT:  zaplogfmt.NewEncoder(l.encoderConfig),
		},
		Selector: l,
		Output:   l,
		Observer: l,
	}
	l.mutex.RUnlock()

	return NewZapLogger(core).Named(name)
}

func (l *Logging) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) {
	l.mutex.RLock()
	observer := l.observer
	l.mutex.RUnlock()

	if observer != nil {
		observer.Check(e, ce)
	}
}

func (l *Logging) WriteEntry(e zapcore.Entry, fields []zapcore.Field) {
	l.mutex.RLock()
	observer := l.observer
	l.mutex.RUnlock()

	if observer != nil {
		observer.WriteEntry(e, fields)
	}
}



func (l *Logging) Logger(name string) *FabricLogger {
	zl := l.ZapLogger(name)
	return NewFabricLogger(zl)
}
