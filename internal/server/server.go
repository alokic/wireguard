package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/qubole/wireguard/internal/logger"
	"bitbucket.org/qubole/wireguard/internal/router"
	"github.com/go-kit/kit/log"
)

var (
	defaultServerPort  = 4000
	defaultRouter      = router.CreateRouter("gorilla")
	defaultLoggerLevel = "info"
	defaultDrainTime   = 5

)

// server struct
type server struct {
	port            int
	router          router.Router
	logger          log.Logger
	notFoundHandler http.Handler
	drainTime       int
	routes          []Route
}

// Option to set params.
type Option func(*server)

// RunFn type to encapsulate a goroutine.
type RunFn func(<-chan struct{}) error

// New is contructor.
func New(options ...Option) *server {
	s := &server{}

	// set defaults
	opts := append([]Option{}, Port(defaultServerPort), Router(defaultRouter),
		Routes(RouteTable), Logger(defaultLoggerLevel), DrainTime(defaultDrainTime),
	)

	// set overrides
	opts = append(opts, options...)

	// apply
	for _, opt := range opts {
		opt(s)
	}

	return s
}



// Runnables returns all parallely executing goroutine for servers.
func (s *server) Runnables() []RunFn {
	h := &http.Server{Addr: ":" + fmt.Sprint(s.port), Handler: s.router}

	fs := []RunFn{}
	fs = append(fs, func(stop <-chan struct{}) error {
		return h.ListenAndServe()
	})

	// shutdown http
	fs = append(fs, func(stop <-chan struct{}) error {
		return s.shutdownServer(h, stop)
	})

	return fs
}

// Port sets port.
func Port(port int) Option {
	return func(s *server) {
		s.port = port
	}
}

// Router sets router.
func Router(r router.Router) Option {
	return func(s *server) {
		s.router = r
	}
}

// DrainTime sets drainTime.
func DrainTime(seconds int) Option {
	return func(s *server) {
		s.drainTime = seconds
	}
}

// Routes sets routes.
func Routes(rs []Route) Option {
	return func(s *server) {
		for _, r := range rs {
			s.handle(r.Method, r.Path, r.Handler(s))
		}

		//	s.router.NotFound(s.notFoundHandler)
	}
}

// Logger sets logger.
// labels are key, val pair .. so even in number always.
func Logger(level string, labels ...interface{}) Option {
	return func(s *server) {
		s.logger = logger.Create(level)

		s.logger = log.With(s.logger, labels...)
	}
}

// NotFoundHandler sets notFoundHandler.
func NotFoundHandler(hn http.Handler) Option {
	return func(s *server) {
		s.notFoundHandler = s.wrapHandlers(hn)
		s.router.NotFound(s.notFoundHandler)
	}
}


func (s *server) shutdownServer(server *http.Server, stop <-chan struct{}) error {
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.drainTime)*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	s.logger.Log("shutdown", "completed", "for", server.Addr, "error", err)
	return err
}

func (s *server) handle(method, path string, handler http.Handler) {
	s.router.Handle(method, path, s.wrapHandlers(handler))
}

func (s *server) recoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				s.notifyException(fmt.Errorf("%v", r))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte{})
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func (s *server) accessControlHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s *server) loggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w}

		defer func(begin time.Time, r *http.Request) {
			s.logger.Log("host", r.Host, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "method", r.Method, "status", sw.status, "content-len", sw.length, "took", time.Since(begin))
		}(time.Now(), r)

		next.ServeHTTP(&sw, r)
	})
}

func (s *server) wrapHandlers(h http.Handler) http.Handler {
	return s.recoverHandler(s.loggerHandler(s.accessControlHandler(h)))
}

func (s *server) notifyException(err error) {

}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}
