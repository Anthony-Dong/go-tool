package pool

import "bytes"

type BufferPool interface {
	Get() *bytes.Buffer
	Put(buffer *bytes.Buffer)
}
