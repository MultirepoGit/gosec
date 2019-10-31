package redigo

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"
)

// Redigo redis
type Redigo struct {
	pool *redigo.Pool
}

// Config of connection
type Config struct {
	MaxActive int
	MaxIdle   int
	Timeout   int
}

// New redis connection using redigo library
func New(ctx context.Context, address string, config *Config) (*Redigo, error) {
	pool := &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", address)
		},
	}

	r := Redigo{
		pool: pool,
	}

	if _, err := r.Ping(ctx); err != nil {
		return nil, err
	}
	return &r, nil
}

// getConn return the connection of redigo
func (rdg *Redigo) getConn(ctx context.Context) (redigo.Conn, error) {
	return rdg.pool.GetContext(ctx)
}

func (rdg *Redigo) do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	conn, err := rdg.getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := conn.Do(cmd, args...)
	return resp, err
}

// Ping the redis
func (rdg *Redigo) Ping(ctx context.Context) (string, error) {
	val, err := rdg.do(ctx, "PING")
	return redigo.String(val, err)
}

// IsErrNil return true if error is nil
func IsErrNil(err error) bool {
	if !errors.Is(err, redigo.ErrNil) {
		return false
	}
	return true
}

// IsResponseOK return true if result value of command is ok
func IsResponseOK(result string) bool {
	if result != "OK" {
		return false
	}
	return true
}