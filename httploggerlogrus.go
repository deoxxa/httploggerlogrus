package httploggerlogrus // import "fknsrs.biz/p/httploggerlogrus"

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

type Logger struct {
	logger logrus.FieldLogger
}

type Option func(l *Logger)

func New(options ...Option) *Logger {
	r := Logger{logger: logrus.StandardLogger()}

	for _, option := range options {
		option(&r)
	}

	return &r
}

func SetLogger(logger logrus.FieldLogger) Option {
	return func(l *Logger) {
		l.logger = logger
	}
}

func (l *Logger) LogRequest(req *http.Request) {
	f := logrus.Fields{
		"http.method": req.Method,
		"http.path":   req.URL.Path,
		"http.host":   req.URL.Host,
	}

	for k := range req.Header {
		f["http.request.header."+k] = req.Header.Get(k)
	}

	l.logger.WithFields(f).Debug("http request")
}

func (l *Logger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	f := logrus.Fields{
		"http.method":   req.Method,
		"http.path":     req.URL.Path,
		"http.host":     req.URL.Host,
		"http.duration": duration.String(),
	}

	if res != nil {
		f["http.status"] = res.Status
	}

	if err != nil {
		f["http.error"] = err.Error()
	}

	for k := range req.Header {
		f["http.request.header."+k] = req.Header.Get(k)
	}
	if res != nil {
		for k := range res.Header {
			f["http.response.header."+k] = res.Header.Get(k)
		}
	}

	l.logger.WithFields(f).Debug("http response")
}
