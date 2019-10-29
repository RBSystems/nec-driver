package nec

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type Projector struct {
	Address string
	Log     Logger

	pool *connpool.Pool
}

var (
	_defaultTTL   = 30 * time.Second
	_defaultDelay = 500 * time.Millisecond
)

type options struct {
	ttl    time.Duration
	delay  time.Duration
	logger Logger
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithConnTTL(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.ttl = t
	})
}

func WithDelay(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.delay = t
	})
}

func WithLogger(l Logger) Option {
	return optionFunc(func(o *options) {
		o.logger = l
	})
}

func NewProjector(addr string, opts ...Option) *Projector {
	options := options{
		ttl:   _defaultTTL,
		delay: _defaultDelay,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	p := &Projector{
		Address: addr,
		pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
	}

	p.pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		d := net.Dialer{}
		return d.DialContext(ctx, "tcp", p.Address+":7142")
	}

	return p
}

// SendCommand sends the byte array to the desired address of projector
func (p *Projector) SendCommand(ctx context.Context, cmd []byte) ([]byte, error) {
	var resp []byte

	err := p.pool.Do(ctx, func(conn connpool.Conn) error {
		conn.SetWriteDeadline(time.Now().Add(3 * time.Second))

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return err
		case n != len(cmd):
			return fmt.Errorf("wrote %v/%v bytes of command 0x%x", n, len(cmd), cmd)
		}

		resp = make([]byte, 5)
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))

		n, err = conn.Read(resp)
		switch {
		case err != nil:
			return err
		case n != len(resp):
			return fmt.Errorf("read %v/%v bytes (read: 0x%x)", n, len(resp), resp)
		}

		rest := make([]byte, (uint8)(resp[4])+1)
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))

		n, err = conn.Read(rest)
		switch {
		case err != nil:
			return err
		case n != len(rest):
			return fmt.Errorf("read %v/%v bytes (read: 0x%x)", n, len(rest), rest)
		}

		resp = append(resp, rest...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// getChecksum returns the checksum value for the end of the hex array
func getChecksum(command []byte) byte {
	var checksum byte
	for i := 0; i < len(command); i++ {
		checksum += command[i]
	}

	return checksum
}
