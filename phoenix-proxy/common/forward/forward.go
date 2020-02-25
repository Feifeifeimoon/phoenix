package forward

import (
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

func Forward(f io.ReadWriteCloser, s io.ReadWriteCloser) (inCount, outCount int64) {
	var wg sync.WaitGroup
	pipe := func(dst io.ReadWriteCloser, src io.ReadWriteCloser, count *int64) {
		defer f.Close()
		defer s.Close()
		defer wg.Done()
		*count, _ = io.Copy(dst, src)
	}
	wg.Add(2)
	go pipe(f, s, &inCount)
	go pipe(s, f, &outCount)
	wg.Wait()
	logrus.Debug("Forward Over", inCount, outCount)
	return
}
