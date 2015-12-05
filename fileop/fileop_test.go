package fileop

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TestLineCount(t *testing.T) {
	Convey("line count test", t, func() {
		Convey("not exit file test", func() {
			cnt, err := LineCount("not exist")
			So(err, ShouldNotBeNil)
			So(cnt, ShouldBeZeroValue)
		})

		Convey("empty file test", func() {
			fileName := "empty"
			file, err := os.Create(fileName)
			file.Close()

			cnt, err := LineCount(fileName)
			So(err, ShouldBeNil)
			So(cnt, ShouldBeZeroValue)
		})

		Convey("non empty file test", func() {
			fileName := "notempty"
			file, _ := os.Create(fileName)
			cnt := 100000
			for i := 0; i < cnt; i++ {
				file.WriteString(fmt.Sprintf("i = %d\n", i))
			}
			file.Close()

			lineCnt, err := LineCount(fileName)
			So(err, ShouldBeNil)
			So(lineCnt, ShouldEqual, cnt)
		})
	})
}
