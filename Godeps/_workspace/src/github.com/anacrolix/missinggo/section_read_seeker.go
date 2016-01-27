package missinggo

import (
	"fmt"
	"io"
	"os"
)

type sectionReadSeeker struct {
	base      io.ReadSeeker
	off, size int64
}

func NewSectionReadSeeker(base io.ReadSeeker, off, size int64) (ret io.ReadSeeker) {
	ret = &sectionReadSeeker{
		base: base,
		off:  off,
		size: size,
	}
	ret.Seek(0, os.SEEK_SET)
	return
}

func (me *sectionReadSeeker) Seek(off int64, whence int) (ret int64, err error) {
	switch whence {
	case os.SEEK_SET:
		off += me.off
	case os.SEEK_CUR:
	case os.SEEK_END:
		off += me.off + me.size
		whence = os.SEEK_SET
	default:
		err = fmt.Errorf("unhandled whence: %d", whence)
		return
	}
	return me.base.Seek(off, whence)
}

func (me *sectionReadSeeker) Read(b []byte) (n int, err error) {
	off, err := me.Seek(0, os.SEEK_CUR)
	if err != nil {
		return
	}
	left := me.off + me.size - off
	if left <= 0 {
		err = io.EOF
		return
	}
	if int64(len(b)) > left {
		b = b[:left]
	}
	return me.base.Read(b)
}
