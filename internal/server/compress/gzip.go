// Package compress предоставляет middleware для сжатия HTTP-ответов и декомпрессии HTTP-запросов
// с использованием алгоритма gzip. Пакет автоматически определяет поддержку gzip клиентом
// и прозрачно обрабатывает сжатие/декомпрессию данных.
package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	acceptEncoding  = "Accept-Encoding"
	contentEncoding = "Content-Encoding"
	contentType     = "Content-Type"
)

// gzipWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type gzipWriter struct {
	w           http.ResponseWriter
	zw          *gzip.Writer
	headersSent bool
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		w:           w,
		zw:          gzip.NewWriter(w),
		headersSent: false,
	}
}

// http.ResponseWriter implementation
func (c *gzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *gzipWriter) Write(p []byte) (int, error) {
	if !c.headersSent {
		c.w.Header().Set(contentEncoding, "gzip")
		c.headersSent = true
	}
	return c.zw.Write(p)
}

func (c *gzipWriter) WriteHeader(statusCode int) {
	if !c.headersSent {
		c.w.Header().Set(contentEncoding, "gzip")
		c.headersSent = true
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *gzipWriter) Close() error {
	return c.zw.Close()
}

// gzipReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// io.ReadCloser implementation
func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// WithGzipCompression создает HTTP-middleware для автоматического сжатия ответов и декомпрессии запросов.
// Middleware проверяет заголовки Accept-Encoding и Content-Encoding для определения
// поддержки gzip клиентом и сервером соответственно.
//
// Функция автоматически:
// - Сжимает ответы сервера, если клиент поддерживает gzip
// - Декомпрессирует запросы клиента, если они сжаты в gzip
// - Устанавливает правильные HTTP-заголовки для сжатых данных
//
// Параметр h - HTTP-обработчик, который будет обернут middleware.
// Возвращает новый HTTP-обработчик с поддержкой gzip-сжатия.
func WithGzipCompression(h http.Handler) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get(acceptEncoding)
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newGzipWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get(contentEncoding)
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(compressFn)
}
