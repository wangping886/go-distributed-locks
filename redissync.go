package go_distributed_locks

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/garyburd/redigo/redis"
	"time"
)

//see doc https://redis.io/commands/set
const expiry = 3 * time.Second

type Mutex struct {
	resorceName  string
	value        string
	genValueFunc func() (string, error)
	pool         redis.Pool
}

func NewMutex(name string, pool redis.Pool) *Mutex {
	return &Mutex{
		resorceName:  name,
		genValueFunc: genValue,
		pool:         pool,
	}
}
func (mux *Mutex) AcquireLock() bool {
	conn := mux.pool.Get()
	defer conn.Close()
	value, err := mux.genValueFunc()
	if err != nil {
		return false
	}
	mux.value = value
	reply, err := redis.String(conn.Do("SET", mux.resorceName, value, "NX", "PX", int(expiry/time.Millisecond)))
	return err == nil && reply == "OK"
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (mux *Mutex) ReleaseLock() bool {
	conn := mux.pool.Get()
	defer conn.Close()
	status, err := redis.Int64(deleteScript.Do(conn, mux.resorceName, mux.value))

	return err == nil && status != 0
}

func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
