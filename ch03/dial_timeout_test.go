package main

import (
	"net"
	"syscall"
	"testing"
	"time"
)

// DialTimeout net.Dialer.Control 함수를 이용하여 Dial 함수의 타임아웃을 테스트한다.
func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        address,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: false,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	// DialTimeout 함수는 추가적으로 타임아웃 기간에 대한 매개변수를 받는다.
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}
	// DialTimeout 함수는 net.DNSError 타입의 에러를 반환하므로 net.Error로 타입 어설션을 해야한다.
	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(err)
	}
	if !nErr.Timeout() {
		t.Fatal("error is not a timeout")
	}
}
