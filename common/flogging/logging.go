/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging/fabenc"
	logging "github.com/op/go-logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type Config struct {
	
	
	
	
	
	
	
	Format string

	
	
	
	
	LogSpec string

	
	
	
	Writer io.Writer
}




type Logging struct {
	*ModuleLevels

	mutex          sync.RWMutex
	encoding       Encoding
	encoderConfig  zapcore.EncoderConfig
	multiFormatter *fabenc.MultiFormatter
	writer         zapcore.WriteSyncer
}



func New(c Config) (*Logging, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.NameKey = "name"

	s := &Logging{
		ModuleLevels: &ModuleLevels{
			defaultLevel: defaultLevel,
		},
		encoderConfig:  encoderConfig,
		multiFormatter: fabenc.NewMultiFormatter(),
	}

	err := s.Apply(c)
	if err != nil {
		return nil, err
	}
	return s, nil
}


func (s *Logging) Apply(c Config) error {
	err := s.SetFormat(c.Format)
	if err != nil {
		return err
	}

	if c.LogSpec == "" {
		c.LogSpec = "INFO"
	}

	err = s.ModuleLevels.ActivateSpec(c.LogSpec)
	if err != nil {
		return err
	}

	if c.Writer == nil {
		c.Writer = os.Stderr
	}
	s.SetWriter(c.Writer)

	var formatter logging.Formatter
	if s.Encoding() == JSON {
		formatter = SetFormat(defaultFormat)
	} else {
		formatter = SetFormat(c.Format)
	}

	InitBackend(formatter, c.Writer)

	return nil
}





func (s *Logging) SetFormat(format string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if format == "" {
		format = defaultFormat
	}

	if format == "json" {
		s.encoding = JSON
		return nil
	}

	formatters, err := fabenc.ParseFormat(format)
	if err != nil {
		return err
	}
	s.multiFormatter.SetFormatters(formatters)
	s.encoding = CONSOLE

	return nil
}




func (s *Logging) SetWriter(w io.Writer) {
	var sw zapcore.WriteSyncer
	switch t := w.(type) {
	case *os.File:
		sw = zapcore.Lock(t)
	case zapcore.WriteSyncer:
		sw = t
	default:
		sw = zapcore.AddSync(w)
	}

	s.mutex.Lock()
	s.writer = sw
	s.mutex.Unlock()
}




func (s *Logging) Write(b []byte) (int, error) {
	s.mutex.RLock()
	w := s.writer
	s.mutex.RUnlock()

	return w.Write(b)
}



func (s *Logging) Sync() error {
	s.mutex.RLock()
	w := s.writer
	s.mutex.RUnlock()

	return w.Sync()
}



func (s *Logging) Encoding() Encoding {
	s.mutex.RLock()
	e := s.encoding
	s.mutex.RUnlock()
	return e
}




func (s *Logging) ZapLogger(module string) *zap.Logger {
	s.mutex.RLock()
	module = strings.Replace(module, "/", ".", -1)
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		
		
		
		return true
	})
	core := &Core{
		LevelEnabler: levelEnabler,
		Levels:       s.ModuleLevels,
		Encoders: map[Encoding]zapcore.Encoder{
			JSON:    zapcore.NewJSONEncoder(s.encoderConfig),
			CONSOLE: fabenc.NewFormatEncoder(s.multiFormatter),
		},
		Selector: s,
		Output:   s,
	}
	s.mutex.RUnlock()

	return NewZapLogger(core).Named(module)
}




func (s *Logging) Logger(module string) *FabricLogger {
	zl := s.ZapLogger(module)
	return NewFabricLogger(zl)
}
