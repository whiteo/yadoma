// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamWriter(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		sendFunc  func([]byte) error
		expectErr bool
		expectLen int
	}{
		{
			name: "successful write",
			data: []byte("hello world"),
			sendFunc: func(data []byte) error {
				return nil
			},
			expectErr: false,
			expectLen: 11,
		},
		{
			name: "send function error",
			data: []byte("test data"),
			sendFunc: func(data []byte) error {
				return errors.New("send error")
			},
			expectErr: true,
			expectLen: 0,
		},
		{
			name: "empty data",
			data: []byte{},
			sendFunc: func(data []byte) error {
				return nil
			},
			expectErr: false,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &streamWriter{send: tt.sendFunc}
			n, err := writer.Write(tt.data)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectLen, n)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectLen, n)
			}
		})
	}
}

func TestStreamReader(t *testing.T) {
	createSuccessfulSendFunc := func() (func([]byte) error, *[][]byte) {
		var received [][]byte
		return func(data []byte) error {
			received = append(received, append([]byte(nil), data...))
			return nil
		}, &received
	}

	createErrorSendFunc := func() (func([]byte) error, *[][]byte) {
		return func(data []byte) error {
			return errors.New("send failed")
		}, nil
	}

	tests := []struct {
		name      string
		input     string
		sendFunc  func([]byte) error
		expectErr bool
		setupSend func() (func([]byte) error, *[][]byte)
	}{
		{
			name:      "successful stream reading",
			input:     "hello world test data",
			setupSend: createSuccessfulSendFunc,
			expectErr: false,
		},
		{
			name:      "empty reader",
			input:     "",
			setupSend: createSuccessfulSendFunc,
			expectErr: false,
		},
		{
			name:      "send function error",
			input:     "test data",
			setupSend: createErrorSendFunc,
			expectErr: true,
		},
		{
			name:      "large data stream",
			input:     strings.Repeat("a", 2048),
			setupSend: createSuccessfulSendFunc,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			sendFunc, received := tt.setupSend()

			err := StreamReader(reader, sendFunc)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if received != nil && tt.input != "" {
					var totalReceived bytes.Buffer
					for _, chunk := range *received {
						totalReceived.Write(chunk)
					}
					assert.Equal(t, tt.input, totalReceived.String())
				}
			}
		})
	}
}

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestStreamDecoder(t *testing.T) {
	createSuccessfulTestStructSendFunc := func() (func(TestStruct) error, *[]TestStruct) {
		var received []TestStruct
		return func(data TestStruct) error {
			received = append(received, data)
			return nil
		}, &received
	}

	createErrorTestStructSendFunc := func() (func(TestStruct) error, *[]TestStruct) {
		return func(data TestStruct) error {
			return errors.New("send failed")
		}, nil
	}

	createNilTestStructSendFunc := func() (func(TestStruct) error, *[]TestStruct) {
		return func(data TestStruct) error {
			return nil
		}, nil
	}

	tests := []struct {
		name      string
		input     string
		sendFunc  func(TestStruct) error
		expectErr bool
		setupSend func() (func(TestStruct) error, *[]TestStruct)
	}{
		{
			name:      "successful JSON decoding",
			input:     `{"name":"test1","value":1}{"name":"test2","value":2}`,
			setupSend: createSuccessfulTestStructSendFunc,
			expectErr: false,
		},
		{
			name:      "single JSON object",
			input:     `{"name":"single","value":42}`,
			setupSend: createSuccessfulTestStructSendFunc,
			expectErr: false,
		},
		{
			name:      "empty input",
			input:     "",
			setupSend: createSuccessfulTestStructSendFunc,
			expectErr: false,
		},
		{
			name:      "invalid JSON",
			input:     `{"name":"test","invalid":}`,
			setupSend: createNilTestStructSendFunc,
			expectErr: true,
		},
		{
			name:      "send function error",
			input:     `{"name":"test","value":1}`,
			setupSend: createErrorTestStructSendFunc,
			expectErr: true,
		},
		{
			name:      "multiple JSON objects with newlines",
			input:     "{\"name\":\"test1\",\"value\":1}\n{\"name\":\"test2\",\"value\":2}\n",
			setupSend: createSuccessfulTestStructSendFunc,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			sendFunc, received := tt.setupSend()

			err := StreamDecoder[TestStruct](reader, sendFunc)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if received == nil || tt.input == "" {
				return
			}
			assertReceivedTestStructs(t, received, tt.input)
		})
	}
}

func countTestStructObjects(input string) int {
	decoder := json.NewDecoder(strings.NewReader(input))
	count := 0
	for {
		var temp TestStruct
		if err := decoder.Decode(&temp); err == io.EOF {
			break
		} else if err != nil {
			break
		}
		count++
	}
	return count
}

func assertReceivedTestStructs(t *testing.T, received *[]TestStruct, input string) {
	expectedCount := countTestStructObjects(input)
	if expectedCount == 0 {
		return
	}
	assert.Len(t, *received, expectedCount)
	if len(*received) > 0 {
		first := (*received)[0]
		assert.NotEmpty(t, first.Name)
	}
}

func TestStreamDecoderWithDifferentTypes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "string type",
			input:     `"hello""world"`,
			expectErr: false,
		},
		{
			name:      "int type",
			input:     `12345`,
			expectErr: false,
		},
		{
			name:      "array type",
			input:     `[1,2,3][4,5,6]`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var received []interface{}

			sendFunc := func(data interface{}) error {
				received = append(received, data)
				return nil
			}

			err := StreamDecoder[interface{}](reader, sendFunc)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, received)
			}
		})
	}
}

func TestStreamReaderErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		reader   io.Reader
		sendFunc func([]byte) error
		wantErr  bool
	}{
		{
			name:   "reader error",
			reader: &errorReader{err: errors.New("read error")},
			sendFunc: func(data []byte) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StreamReader(tt.reader, tt.sendFunc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type errorReader struct {
	err error
}

func (r *errorReader) Read(_ []byte) (n int, err error) {
	return 0, r.err
}

func TestStreamDecoderErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		reader   io.Reader
		sendFunc func(TestStruct) error
		wantErr  bool
	}{
		{
			name:   "reader error",
			reader: &errorReader{err: errors.New("read error")},
			sendFunc: func(data TestStruct) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StreamDecoder[TestStruct](tt.reader, tt.sendFunc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
