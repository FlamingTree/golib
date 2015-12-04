package log15handler

import (
	"errors"
	"github.com/BurntSushi/toml"
	. "github.com/FlamingTree/golib/waitquit"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"sync/atomic"
)

// conf := `
//	filename = "foo"
//	maxsize = 5
//	maxage = 10
//	maxbackups = 3
//	localtime = true`[1:]
func RollingFileHandler(conf string, fmtr log15.Format) (log15.Handler, error) {
	var (
		lj  lumberjack.Logger
		err error
	)
	if _, err = toml.Decode(conf, &lj); err != nil {
		return nil, err
	} else if _, err = lj.Write([]byte("init test")); err != nil {
		return nil, err
	}
	return log15.StreamHandler(&lj, fmtr), nil
}

type SafeBufferHandler struct {
	bufSize  int
	recs     chan *log15.Record
	wq       WaitQuit
	exitFlag int32
	h        log15.Handler
}

func NewSafeBufferHandler(bufSize int, h log15.Handler) log15.Handler {
	bufferHandler := &SafeBufferHandler{
		bufSize: bufSize,
		recs:    make(chan *log15.Record, bufSize),
		wq:      NewWaitQuit(),
		h:       h,
	}
	bufferHandler.asyncLoop()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return bufferHandler
}

var ErrorProcessExit = errors.New("process exiting")

func (h *SafeBufferHandler) Log(r *log15.Record) error {
	if h.Exiting() {
		log.Println(ErrorProcessExit)
		return ErrorProcessExit
	}

	h.recs <- r
	return nil
}

func (h *SafeBufferHandler) Exiting() bool {
	return atomic.LoadInt32(&h.exitFlag) == 1
}

func (h *SafeBufferHandler) asyncLoop() {
	f := func() {
		exitChan := h.wq.Done()
		for {
			select {
			case m, ok := <-h.recs:
				if !ok {
					return
				}
				_ = h.h.Log(m)
			case <-exitChan:
				if atomic.CompareAndSwapInt32(&h.exitFlag, 0, 1) {
					close(h.recs)
				}
				exitChan = nil
			}
		}
	}
	h.wq.Wrap(f)
}

func (h *SafeBufferHandler) Exit() {
	h.wq.Exit()
}
