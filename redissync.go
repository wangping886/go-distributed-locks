package go_distributed_locks

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/garyburd/redigo/redis"
	"time"
)

//see doc https://redis.io/commands/set
const expiry = 3 * time.Second

func AcquireLock(redisPool redis.Pool, resourceName string, randValue string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("SET", resourceName, randValue, "NX", "PX", int(expiry/time.Millisecond)))
	return err == nil && reply == "OK"
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func ReleaseLock(redisPool redis.Pool, resourceName string, randValue string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	status, err := redis.Int64(deleteScript.Do(conn, resourceName, randValue))

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
