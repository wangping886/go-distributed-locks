package go_distributed_locks

import (
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestMutexLock(t *testing.T) {
	mutex := NewMutex("resource-name", redis.Pool{})
	mutex.AcquireLock()
	defer mutex.ReleaseLock()
	//your process logic
}
