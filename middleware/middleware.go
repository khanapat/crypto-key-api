package middleware

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestInfoMsg    = "Request Information"
	responseInfoMsg   = "Response Information"
	xRequestID        = "X-Request-ID"
	contentType       = "Content-Type"
	applicationJSON   = "application/json"
	accessControl     = "Access-Control-Allow-Origin"
	contentLength     = "Content-Length"
	contentLengthByte = 512
	key               = "logger"
)

var (
	tempBody string
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
		r.Body = &HackReqBody{
			ReadCloser: r.Body,
			method:     r.Method,
			host:       r.Host,
			requestURI: r.RequestURI,
			remoteAddr: r.RemoteAddr,
			header:     r.Header,
			logger:     log,
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
			if n < contentLengthByte {
				h.logger.Debug(requestInfoMsg,
					zap.String("body", tempBody),
					zap.String("method", h.method),
					zap.String("host", h.host),
					zap.String("path_uri", h.requestURI),
					zap.String("remote_addr", h.remoteAddr),
					zap.String("content_type", h.header.Get(contentType)),
					zap.String(xRequestID, h.header.Get(xRequestID)),
				)
				tempBody = ""
			}
		} else {
			h.logger.Debug(requestInfoMsg,
				zap.String("body", string(body[:n])),
				zap.String("method", h.method),
				zap.String("host", h.host),
				zap.String("path_uri", h.requestURI),
				zap.String("remote_addr", h.remoteAddr),
				zap.String("content_type", h.header.Get(contentType)),
				zap.String(xRequestID, h.header.Get(xRequestID)),
			)
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
			zap.String(xRequestID, h.Header().Get(xRequestID)),
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
