package log15handler

import (
	"bufio"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/inconshreveable/log15.v2"
	"os"
	"testing"
)

func TestRollingFileHandler(t *testing.T) {
	Convey("rolling file handler test", t, func() {
		conf := `
filename = "rolling.log"
maxsize = 100
maxage = 10
maxbackups = 5
localtime = true`[1:]
		_, err := RollingFileHandler(conf, log15.LogfmtFormat())
		So(err, ShouldBeNil)
	})
}

func lineCount(filePath string) (cnt int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cnt++
	}

	err = scanner.Err()
	return
}

func TestSafeBufferChannel(t *testing.T) {
	Convey("safe buffer channel test", t, func() {
		var (
			err                        error
			filePath                   string
			fileHandler, bufferHandler log15.Handler
		)
		filePath = "test.log"
		_ = os.Remove(filePath)
		fileHandler, err = log15.FileHandler(filePath, log15.LogfmtFormat())
		bufferHandler = NewSafeBufferHandler(1000, fileHandler)
		log := log15.New()
		log.SetHandler(bufferHandler)

		cnt := 100000
		for i := 0; i < cnt; i++ {
			log.Info("test", "i", i)
		}
		bufferHandler.(*SafeBufferHandler).Exit()
		lineCnt, err := lineCount(filePath)
		So(err, ShouldBeNil)
		So(lineCnt, ShouldEqual, cnt)
	})
}
