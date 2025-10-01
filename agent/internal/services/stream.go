// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

type streamWriter struct {
	send func([]byte) error
}

func (w *streamWriter) Write(p []byte) (int, error) {
	if err := w.send(p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func StreamReader(r io.Reader, send func([]byte) error) error {
	_, err := io.CopyBuffer(&streamWriter{send}, r, make([]byte, 1024))
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func StreamDecoder[T any](r io.Reader, send func(T) error) error {
	dec := json.NewDecoder(r)
	for {
		var t T
		if err := dec.Decode(&t); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if err := send(t); err != nil {
			// Проверяем, не связана ли ошибка с отменой контекста
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			return err
		}
	}
}
