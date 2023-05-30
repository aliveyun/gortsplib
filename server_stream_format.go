package gortsplib

import (
	"github.com/aliveyun/gortsplib/pkg/formats"
	"github.com/aliveyun/gortsplib/pkg/rtcpsender"
)

type serverStreamFormat struct {
	format     formats.Format
	rtcpSender *rtcpsender.RTCPSender
}
