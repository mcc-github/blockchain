



















package zap

import (
	"sort"
	"time"

	"go.uber.org/zap/zapcore"
)






type SamplingConfig struct {
	Initial    int `json:"initial" yaml:"initial"`
	Thereafter int `json:"thereafter" yaml:"thereafter"`
}














type Config struct {
	
	
	
	Level AtomicLevel `json:"level" yaml:"level"`
	
	
	Development bool `json:"development" yaml:"development"`
	
	
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	
	
	
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`
	
	
	
	Encoding string `json:"encoding" yaml:"encoding"`
	
	
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	
	
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	
	
	
	
	
	
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}



func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}






func NewProductionConfig() Config {
	return Config{
		Level:       NewAtomicLevelAt(InfoLevel),
		Development: false,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}



func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}







func NewDevelopmentConfig() Config {
	return Config{
		Level:            NewAtomicLevelAt(DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}


func (cfg Config) Build(opts ...Option) (*Logger, error) {
	enc, err := cfg.buildEncoder()
	if err != nil {
		return nil, err
	}

	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	log := New(
		zapcore.NewCore(enc, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

func (cfg Config) buildOptions(errSink zapcore.WriteSyncer) []Option {
	opts := []Option{ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, AddCaller())
	}

	stackLevel := ErrorLevel
	if cfg.Development {
		stackLevel = WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, AddStacktrace(stackLevel))
	}

	if cfg.Sampling != nil {
		opts = append(opts, WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, Fields(fs...))
	}

	return opts
}

func (cfg Config) openSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := Open(cfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}
	errSink, _, err := Open(cfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}
	return sink, errSink, nil
}

func (cfg Config) buildEncoder() (zapcore.Encoder, error) {
	return newEncoder(cfg.Encoding, cfg.EncoderConfig)
}
