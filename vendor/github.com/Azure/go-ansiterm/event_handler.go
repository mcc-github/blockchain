package ansiterm

type AnsiEventHandler interface {
	
	Print(b byte) error

	
	Execute(b byte) error

	
	CUU(int) error

	
	CUD(int) error

	
	CUF(int) error

	
	CUB(int) error

	
	CNL(int) error

	
	CPL(int) error

	
	CHA(int) error

	
	VPA(int) error

	
	CUP(int, int) error

	
	HVP(int, int) error

	
	DECTCEM(bool) error

	
	DECOM(bool) error

	
	DECCOLM(bool) error

	
	ED(int) error

	
	EL(int) error

	
	IL(int) error

	
	DL(int) error

	
	ICH(int) error

	
	DCH(int) error

	
	SGR([]int) error

	
	SU(int) error

	
	SD(int) error

	
	DA([]string) error

	
	DECSTBM(int, int) error

	
	IND() error

	
	RI() error

	
	Flush() error
}
