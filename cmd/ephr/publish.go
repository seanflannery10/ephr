package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/seanflannery10/ossa/helpers"
)

func publishVars(db *pgxpool.Pool) {
	expvar.NewString("version").Set(helpers.GetVersion())

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
