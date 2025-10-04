// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

// Package service provides shared utilities for the agent's gRPC service layer.
// It offers helpers to stream bytes and JSON-decoded messages from io.Reader
// sources to server-side send callbacks.
//
// Helpers normalize I/O termination (io.EOF is treated as a clean close),
// propagate context cancellation and deadlines, and avoid spawning goroutines.
// The package does not manage the lifetime of readers or streams; callers are
// responsible for closing resources and consuming streams.
//
// Intended for internal use by higher-level handlers (e.g., container, image,
// network, and system).
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

// Write sends the provided bytes to the underlying send callback.
// It implements io.Writer by forwarding p as a single chunk; on success it
// returns len(p) to report that all bytes were written.
// If the callback returns an error, Write returns that error and reports that
// zero bytes were written. The input slice is neither retained nor modified.
func (w *streamWriter) Write(p []byte) (int, error) {
	if err := w.send(p); err != nil {
		return 0, err
	}
	return len(p), nil
}

// StreamReader streams raw bytes from r to the provided send callback.
// It copies data using a fixed 1 KiB buffer, invoking send once per chunk.
// Reaching io.EOF is treated as a clean termination and returns nil.
// Any other error from the reader or the send callback is returned as-is.
// The function does not close r, retains no buffers, and spawns no goroutines.
func StreamReader(r io.Reader, send func([]byte) error) error {
	_, err := io.CopyBuffer(&streamWriter{send}, r, make([]byte, 1024))
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

// StreamDecoder streams JSON-decoded values of type T from r to the send callback.
// It uses encoding/json.Decoder to read a sequence of concatenated JSON values and
// invokes send once per decoded item.
// Reaching io.EOF is treated as normal completion and returns nil.
// Any other decode error, or an error returned by send (including context cancellation
// or deadline expiration), is propagated to the caller.
// The function does not close r, keep long-lived buffers, or spawn goroutines.
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
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			return err
		}
	}
}
