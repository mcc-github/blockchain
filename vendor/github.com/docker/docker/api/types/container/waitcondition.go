package container 



type WaitCondition string












const (
	WaitConditionNotRunning WaitCondition = "not-running"
	WaitConditionNextExit   WaitCondition = "next-exit"
	WaitConditionRemoved    WaitCondition = "removed"
)
