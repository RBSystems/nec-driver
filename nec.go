package necdriver

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/byuoitav/pooled"
)

type Projector struct {
	poolInit sync.Once
	pool     *pooled.Pool
}

/*
var commands = map[string][]byte{
	"MuteOn":      {0x02, 0x12, 0x00, 0x00, 0x00, 0x14},
	"MuteOff":     {0x02, 0x13, 0x00, 0x00, 0x00, 0x15},
	"Volume":      {0x03, 0x10, 0x00, 0x00, 0x05, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00}, // Used for changing the volume level of the projector
	"VolumeLevel": {0x03, 0x04, 0x00, 0x00, 0x03, 0x05, 0x00, 0x00, 0x00},             // Used for getting the volume level
}
*/

func getConnection(key interface{}) (pooled.Conn, error) {
	address, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	conn, err := net.DialTimeout("tcp", address+":7142", 10*time.Second)
	if err != nil {
		return nil, err
	}

	return pooled.Wrap(conn), nil
}

// SendCommand sends the byte array to the desired address of projector
func (p *Projector) SendCommand(ctx context.Context, addr string, cmd []byte) ([]byte, error) {
	p.poolInit.Do(func() {
		// create the pool
		p.pool = pooled.NewPool(45*time.Second, 400*time.Millisecond, getConnection)
	})

	var resp []byte
	err := p.pool.Do(addr, func(conn pooled.Conn) error {
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
