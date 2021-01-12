package middleware

import (
	"context"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestInfoMsg    string = "Request Information"
	responseInfoMsg   string = "Response Information"
	xRequestID        string = "X-Request-ID"
	contentType       string = "Content-Type"
	applicationJSON   string = "application/json"
	accessControl     string = "Access-Control-Allow-Origin"
	contentLength     string = "Content-Length"
	contentLengthByte int    = 512
	key               string = "logger"
)

var (
	tempBody  string
	tempCount int
)

type middleware struct {
	ZapLogger *zap.Logger
}

func NewMiddleware(zapLogger *zap.Logger) *middleware {
	return &middleware{
		ZapLogger: zapLogger,
	}
}

func ContextData(ctx context.Context) *zap.Logger {
	v := ctx.Value(key)
	if v == nil {
		return nil
	}
	l, ok := v.(*zap.Logger)
	if ok {
		return l
	}
	return nil
}

func (m *middleware) ContextLogAndLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.ZapLogger.With(
			zap.String(xRequestID, uuid.New().String()),
		)
		if r.Method == "GET" {
			log.Debug(requestInfoMsg,
				zap.String("body", tempBody),
				zap.String("method", r.Method),
				zap.String("host", r.Host),
				zap.String("path_uri", r.RequestURI),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("content_type", r.Header.Get(contentType)),
			)
		} else {
			r.Body = &HackReqBody{
				ReadCloser: r.Body,
				method:     r.Method,
				host:       r.Host,
				requestURI: r.RequestURI,
				remoteAddr: r.RemoteAddr,
				header:     r.Header,
				logger:     log,
			}
		}
		w.Header().Set(contentType, applicationJSON)
		next.ServeHTTP(&HackResBody{w, log}, r.WithContext(context.WithValue(r.Context(), key, log)))
	})
}

type HackReqBody struct {
	io.ReadCloser
	method     string
	host       string
	requestURI string
	remoteAddr string
	header     http.Header
	logger     *zap.Logger
}

func (h *HackReqBody) Read(body []byte) (int, error) {
	var n int
	var err error
	defer func() {
		if stringToInt(h.header.Get(contentLength)) > contentLengthByte {
			tempBody += string(body[:n])
			if n < (contentLengthByte * int(math.Pow(2.0, float64(tempCount)))) {
				h.logger.Debug(requestInfoMsg,
					zap.String("body", tempBody),
					zap.String("method", h.method),
					zap.String("host", h.host),
					zap.String("path_uri", h.requestURI),
					zap.String("remote_addr", h.remoteAddr),
					zap.String("content_type", h.header.Get(contentType)),
				)
				tempBody = ""
				tempCount = 0
			} else {
				tempCount++
			}
		} else {
			h.logger.Debug(requestInfoMsg,
				zap.String("body", string(body[:n])),
				zap.String("method", h.method),
				zap.String("host", h.host),
				zap.String("path_uri", h.requestURI),
				zap.String("remote_addr", h.remoteAddr),
				zap.String("content_type", h.header.Get(contentType)),
			)
			tempBody = ""
			tempCount = 0
		}
	}()
	n, err = h.ReadCloser.Read(body)
	return n, err
}

type HackResBody struct {
	http.ResponseWriter
	logger *zap.Logger
}

func (h *HackResBody) Write(b []byte) (int, error) {
	defer func() {
		h.logger.Debug(responseInfoMsg,
			zap.String("body", string(b)),
			zap.String("content_type", h.Header().Get(contentType)),
		)
	}()

	return h.ResponseWriter.Write(b)
}

func stringToInt(numberStr string) int {
	numberInt, err := strconv.Atoi(numberStr)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return numberInt
}
