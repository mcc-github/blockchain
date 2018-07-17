package kingpin

import (
	"net"
	"net/url"
	"os"
	"time"

	"github.com/alecthomas/units"
)

type Settings interface {
	SetValue(value Value)
}

type parserMixin struct {
	value    Value
	required bool
}

func (p *parserMixin) SetValue(value Value) {
	p.value = value
}


func (p *parserMixin) StringMap() (target *map[string]string) {
	target = &(map[string]string{})
	p.StringMapVar(target)
	return
}


func (p *parserMixin) Duration() (target *time.Duration) {
	target = new(time.Duration)
	p.DurationVar(target)
	return
}


func (p *parserMixin) Bytes() (target *units.Base2Bytes) {
	target = new(units.Base2Bytes)
	p.BytesVar(target)
	return
}


func (p *parserMixin) IP() (target *net.IP) {
	target = new(net.IP)
	p.IPVar(target)
	return
}


func (p *parserMixin) TCP() (target **net.TCPAddr) {
	target = new(*net.TCPAddr)
	p.TCPVar(target)
	return
}


func (p *parserMixin) TCPVar(target **net.TCPAddr) {
	p.SetValue(newTCPAddrValue(target))
}


func (p *parserMixin) ExistingFile() (target *string) {
	target = new(string)
	p.ExistingFileVar(target)
	return
}


func (p *parserMixin) ExistingDir() (target *string) {
	target = new(string)
	p.ExistingDirVar(target)
	return
}


func (p *parserMixin) ExistingFileOrDir() (target *string) {
	target = new(string)
	p.ExistingFileOrDirVar(target)
	return
}


func (p *parserMixin) File() (target **os.File) {
	target = new(*os.File)
	p.FileVar(target)
	return
}


func (p *parserMixin) OpenFile(flag int, perm os.FileMode) (target **os.File) {
	target = new(*os.File)
	p.OpenFileVar(target, flag, perm)
	return
}


func (p *parserMixin) URL() (target **url.URL) {
	target = new(*url.URL)
	p.URLVar(target)
	return
}


func (p *parserMixin) StringMapVar(target *map[string]string) {
	p.SetValue(newStringMapValue(target))
}


func (p *parserMixin) Float() (target *float64) {
	return p.Float64()
}


func (p *parserMixin) FloatVar(target *float64) {
	p.Float64Var(target)
}


func (p *parserMixin) DurationVar(target *time.Duration) {
	p.SetValue(newDurationValue(target))
}


func (p *parserMixin) BytesVar(target *units.Base2Bytes) {
	p.SetValue(newBytesValue(target))
}


func (p *parserMixin) IPVar(target *net.IP) {
	p.SetValue(newIPValue(target))
}


func (p *parserMixin) ExistingFileVar(target *string) {
	p.SetValue(newExistingFileValue(target))
}


func (p *parserMixin) ExistingDirVar(target *string) {
	p.SetValue(newExistingDirValue(target))
}


func (p *parserMixin) ExistingFileOrDirVar(target *string) {
	p.SetValue(newExistingFileOrDirValue(target))
}


func (p *parserMixin) FileVar(target **os.File) {
	p.SetValue(newFileValue(target, os.O_RDONLY, 0))
}


func (p *parserMixin) OpenFileVar(target **os.File, flag int, perm os.FileMode) {
	p.SetValue(newFileValue(target, flag, perm))
}


func (p *parserMixin) URLVar(target **url.URL) {
	p.SetValue(newURLValue(target))
}


func (p *parserMixin) URLList() (target *[]*url.URL) {
	target = new([]*url.URL)
	p.URLListVar(target)
	return
}


func (p *parserMixin) URLListVar(target *[]*url.URL) {
	p.SetValue(newURLListValue(target))
}


func (p *parserMixin) Enum(options ...string) (target *string) {
	target = new(string)
	p.EnumVar(target, options...)
	return
}


func (p *parserMixin) EnumVar(target *string, options ...string) {
	p.SetValue(newEnumFlag(target, options...))
}


func (p *parserMixin) Enums(options ...string) (target *[]string) {
	target = new([]string)
	p.EnumsVar(target, options...)
	return
}


func (p *parserMixin) EnumsVar(target *[]string, options ...string) {
	p.SetValue(newEnumsFlag(target, options...))
}


func (p *parserMixin) Counter() (target *int) {
	target = new(int)
	p.CounterVar(target)
	return
}

func (p *parserMixin) CounterVar(target *int) {
	p.SetValue(newCounterValue(target))
}
