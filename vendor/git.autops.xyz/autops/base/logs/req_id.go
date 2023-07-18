package logs

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"time"
)

var pid = uint32(os.Getpid())

//GenReqID req id
func GenReqID() string {
	var b [12]byte
	binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UnixNano()))
	binary.BigEndian.PutUint32(b[8:], pid)
	return base64.URLEncoding.EncodeToString(b[:])
}
