












package expfmt

import (
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/prometheus/common/internal/bitbucket.org/ww/goautoneg"

	dto "github.com/prometheus/client_model/go"
)


type Encoder interface {
	Encode(*dto.MetricFamily) error
}

type encoder func(*dto.MetricFamily) error

func (e encoder) Encode(v *dto.MetricFamily) error {
	return e(v)
}



func Negotiate(h http.Header) Format {
	for _, ac := range goautoneg.ParseAccept(h.Get(hdrAccept)) {
		
		if ac.Type+"/"+ac.SubType == ProtoType && ac.Params["proto"] == ProtoProtocol {
			switch ac.Params["encoding"] {
			case "delimited":
				return FmtProtoDelim
			case "text":
				return FmtProtoText
			case "compact-text":
				return FmtProtoCompact
			}
		}
		
		ver := ac.Params["version"]
		if ac.Type == "text" && ac.SubType == "plain" && (ver == TextVersion || ver == "") {
			return FmtText
		}
	}
	return FmtText
}


func NewEncoder(w io.Writer, format Format) Encoder {
	switch format {
	case FmtProtoDelim:
		return encoder(func(v *dto.MetricFamily) error {
			_, err := pbutil.WriteDelimited(w, v)
			return err
		})
	case FmtProtoCompact:
		return encoder(func(v *dto.MetricFamily) error {
			_, err := fmt.Fprintln(w, v.String())
			return err
		})
	case FmtProtoText:
		return encoder(func(v *dto.MetricFamily) error {
			_, err := fmt.Fprintln(w, proto.MarshalTextString(v))
			return err
		})
	case FmtText:
		return encoder(func(v *dto.MetricFamily) error {
			_, err := MetricFamilyToText(w, v)
			return err
		})
	}
	panic("expfmt.NewEncoder: unknown format")
}
