package sockets

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func runTest(t *testing.T, path string, l net.Listener, echoStr string) {
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			_, _ = conn.Write([]byte(echoStr))
			_ = conn.Close()
		}
	}()

	conn, err := net.Dial("unix", path)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 5)
	if _, err := conn.Read(buf); err != nil {
		t.Fatal(err)
	} else if string(buf) != echoStr {
		t.Fatal(fmt.Errorf("msg may lost"))
	}
}

// TestNewUnixSocket run under root user.
func TestNewUnixSocket(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("requires root")
	}
	gid := os.Getgid()
	path := "/tmp/test.sock"
	echoStr := "hello"
	l, err := NewUnixSocket(path, gid)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	runTest(t, path, l, echoStr)
}

func TestUnixSocketWithOpts(t *testing.T) {
	socketFile, err := os.CreateTemp("", "test*.sock")
	if err != nil {
		t.Fatal(err)
	}
	_ = socketFile.Close()
	defer func() { _ = os.Remove(socketFile.Name()) }()

	l := createTestUnixSocket(t, socketFile.Name())
	defer func() { _ = l.Close() }()

	echoStr := "hello"
	runTest(t, socketFile.Name(), l, echoStr)
}
