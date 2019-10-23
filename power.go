package necdriver

import (
	"context"
	"fmt"
	"strings"
)

var (
	// PowerStatus gets the projector's power status
	PowerStatus = []byte{0x00, 0x85, 0x00, 0x00, 0x01, 0x01, 0x87}

	// PowerOn powers on the projector
	PowerOn = []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}

	// PowerStandby powers off the projector
	PowerStandby = []byte{0x02, 0x01, 0x00, 0x00, 0x00, 0x03}
)

func (p *Projector) GetPower(ctx context.Context, addr string) (string, error) {
	resp, err := p.SendCommand(ctx, addr, PowerStatus)
	switch {
	case err != nil:
		return "", err
	case len(resp) < 8:
		return "", fmt.Errorf("bad response from projector: 0x%x", resp)
	case resp[7] == 0b1:
		return "on", nil
	default:
		return "standby", nil
	}
}

func (p *Projector) SetPower(ctx context.Context, addr string, power string) error {
	var cmd []byte
	switch {
	case strings.EqualFold(power, "on"):
		cmd = PowerOn
	case strings.EqualFold(power, "standby"):
		cmd = PowerStandby
	default:
		return fmt.Errorf("unable to set power state to %q: must be %q or %q", power, "on", "standby")
	}

	_, err := p.SendCommand(ctx, addr, cmd)
	return err
}
