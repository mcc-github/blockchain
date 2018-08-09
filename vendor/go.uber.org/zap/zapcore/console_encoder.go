



















package zapcore

import (
	"fmt"
	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/internal/bufferpool"
)

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type consoleEncoder struct {
	*jsonEncoder
}









func NewConsoleEncoder(cfg EncoderConfig) Encoder {
	return consoleEncoder{newJSONEncoder(cfg, true)}
}

func (c consoleEncoder) Clone() Encoder {
	return consoleEncoder{c.jsonEncoder.Clone().(*jsonEncoder)}
}

func (c consoleEncoder) EncodeEntry(ent Entry, fields []Field) (*buffer.Buffer, error) {
	line := bufferpool.Get()

	
	
	
	
	
	
	arr := getSliceEncoder()
	if c.TimeKey != "" && c.EncodeTime != nil {
		c.EncodeTime(ent.Time, arr)
	}
	if c.LevelKey != "" && c.EncodeLevel != nil {
		c.EncodeLevel(ent.Level, arr)
	}
	if ent.LoggerName != "" && c.NameKey != "" {
		nameEncoder := c.EncodeName

		if nameEncoder == nil {
			
			nameEncoder = FullNameEncoder
		}

		nameEncoder(ent.LoggerName, arr)
	}
	if ent.Caller.Defined && c.CallerKey != "" && c.EncodeCaller != nil {
		c.EncodeCaller(ent.Caller, arr)
	}
	for i := range arr.elems {
		if i > 0 {
			line.AppendByte('\t')
		}
		fmt.Fprint(line, arr.elems[i])
	}
	putSliceEncoder(arr)

	
	if c.MessageKey != "" {
		c.addTabIfNecessary(line)
		line.AppendString(ent.Message)
	}

	
	c.writeContext(line, fields)

	
	
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	if c.LineEnding != "" {
		line.AppendString(c.LineEnding)
	} else {
		line.AppendString(DefaultLineEnding)
	}
	return line, nil
}

func (c consoleEncoder) writeContext(line *buffer.Buffer, extra []Field) {
	context := c.jsonEncoder.Clone().(*jsonEncoder)
	defer context.buf.Free()

	addFields(context, extra)
	context.closeOpenNamespaces()
	if context.buf.Len() == 0 {
		return
	}

	c.addTabIfNecessary(line)
	line.AppendByte('{')
	line.Write(context.buf.Bytes())
	line.AppendByte('}')
}

func (c consoleEncoder) addTabIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendByte('\t')
	}
}
