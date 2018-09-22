package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// NewPool creates a new redis connection pool.
func NewPool(config Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		Wait:    true, // Wait for the connection pool, no connection pool exhausted error
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				config.Server(),
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to dial redis server")
			}

			if len(config.Password) > 0 {
				var err error

				for _, password := range config.Password {
					_, err = c.Do("AUTH", password)
					if err == nil {
						break
					}
				}

				if err != nil {
					c.Close()

					return nil, errors.Wrap(err, "none of the provided passwords were accepted by the server")
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			return err
		},
	}
}
