package handler

import "net/http"

type WrapWriter struct {
	http.ResponseWriter
	status int
}

func NewWrapWriter(w http.ResponseWriter) *WrapWriter {
	return &WrapWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *WrapWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *WrapWriter) Status() int {
	return w.status
}
