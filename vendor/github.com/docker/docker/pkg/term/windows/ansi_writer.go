

package windowsconsole 

import (
	"io"
	"os"

	ansiterm "github.com/Azure/go-ansiterm"
	"github.com/Azure/go-ansiterm/winterm"
)


type ansiWriter struct {
	file           *os.File
	fd             uintptr
	infoReset      *winterm.CONSOLE_SCREEN_BUFFER_INFO
	command        []byte
	escapeSequence []byte
	inAnsiSequence bool
	parser         *ansiterm.AnsiParser
}



func NewAnsiWriter(nFile int) io.Writer {
	initLogger()
	file, fd := winterm.GetStdFile(nFile)
	info, err := winterm.GetConsoleScreenBufferInfo(fd)
	if err != nil {
		return nil
	}

	parser := ansiterm.CreateParser("Ground", winterm.CreateWinEventHandler(fd, file))
	logger.Infof("newAnsiWriter: parser %p", parser)

	aw := &ansiWriter{
		file:           file,
		fd:             fd,
		infoReset:      info,
		command:        make([]byte, 0, ansiterm.ANSI_MAX_CMD_LENGTH),
		escapeSequence: []byte(ansiterm.KEY_ESC_CSI),
		parser:         parser,
	}

	logger.Infof("newAnsiWriter: aw.parser %p", aw.parser)
	logger.Infof("newAnsiWriter: %v", aw)
	return aw
}

func (aw *ansiWriter) Fd() uintptr {
	return aw.fd
}


func (aw *ansiWriter) Write(p []byte) (total int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	logger.Infof("Write: % x", p)
	logger.Infof("Write: %s", string(p))
	return aw.parser.Parse(p)
}
