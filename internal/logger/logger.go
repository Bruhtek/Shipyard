package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func Init(isDev bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if isDev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Str("log level", "DBG").Msg("Logger initialized")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Str("log level", "INF").Msg("Logger initialized")
	}
}

func HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		tStart := time.Now()

		defer func() {
			log.Debug().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", ww.Status()).
				Int("size", ww.BytesWritten()).
				Dur("duration", time.Since(tStart)).
				Msg("Request")
		}()

		next.ServeHTTP(ww, r)
	})
}

func HandleSimpleRecoverPanic(r any, message string) {
	if r == nil {
		return
	}
	err, ok := r.(error)
	if ok {
		log.Err(err).Msg(message)
	} else {
		log.Error().Msg(message + " - unable to cast to error")
	}
}
