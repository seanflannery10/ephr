package main

import (
	"expvar"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/seanflannery10/ossa/version"
	"runtime"
	"time"
)

func publishVars(db *pgxpool.Pool) {
	expvar.NewString("version").Set(version.Get())

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stat()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
}
