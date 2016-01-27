package muxtest

import (
	multiplex "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-stream-muxer/multiplex"
	multistream "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-stream-muxer/multistream"
	muxado "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-stream-muxer/muxado"
	spdy "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-stream-muxer/spdystream"
	yamux "github.com/noffle/ipget/Godeps/_workspace/src/github.com/jbenet/go-stream-muxer/yamux"
)

var _ = multiplex.DefaultTransport
var _ = multistream.NewTransport
var _ = muxado.Transport
var _ = spdy.Transport
var _ = yamux.DefaultTransport
