













package expfmt


type Format string


const (
	TextVersion   = "0.0.4"
	ProtoType     = `application/vnd.google.protobuf`
	ProtoProtocol = `io.prometheus.client.MetricFamily`
	ProtoFmt      = ProtoType + "; proto=" + ProtoProtocol + ";"

	
	FmtUnknown      Format = `<unknown>`
	FmtText         Format = `text/plain; version=` + TextVersion + `; charset=utf-8`
	FmtProtoDelim   Format = ProtoFmt + ` encoding=delimited`
	FmtProtoText    Format = ProtoFmt + ` encoding=text`
	FmtProtoCompact Format = ProtoFmt + ` encoding=compact-text`
)

const (
	hdrContentType = "Content-Type"
	hdrAccept      = "Accept"
)
