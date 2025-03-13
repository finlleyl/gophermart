package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	// Простой обработчик, который возвращает "Hello, World!"
	handler := GzipMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	handler(w, req)
	resp := w.Result()

	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"), "Должен быть установлен заголовок Content-Encoding: gzip")

	// Чтение и распаковка ответа
	gzReader, err := gzip.NewReader(resp.Body)
	assert.NoError(t, err, "Не удалось создать gzip.Reader")
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	assert.NoError(t, err, "Ошибка чтения распакованного ответа")
	assert.Equal(t, "Hello, World!", string(decompressed))
}
