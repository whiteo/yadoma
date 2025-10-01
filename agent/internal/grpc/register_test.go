// @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)

package grpcserver

import (
	"testing"

	"google.golang.org/grpc"
)

type fakeRegistrator struct {
	called int
}

func (f *fakeRegistrator) Register(_ *grpc.Server) {
	f.called++
}

func TestRegisterAllCallsRegisterOnEach(t *testing.T) {
	srv := grpc.NewServer()
	t.Cleanup(srv.Stop)

	f1 := &fakeRegistrator{}
	f2 := &fakeRegistrator{}

	RegisterAll(srv, f1, f2)

	if f1.called != 1 || f2.called != 1 {
		t.Fatalf("expected each registrator to be called once, got f1=%d f2=%d", f1.called, f2.called)
	}
}

func TestRegisterAllNoServicesNoPanic(t *testing.T) {
	srv := grpc.NewServer()
	t.Cleanup(srv.Stop)

	RegisterAll(srv)
}
