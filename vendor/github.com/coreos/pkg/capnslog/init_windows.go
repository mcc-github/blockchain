













package capnslog

import "os"

func init() {
	initHijack()

	
	SetFormatter(NewPrettyFormatter(os.Stderr, false))
	SetGlobalLogLevel(INFO)
}
