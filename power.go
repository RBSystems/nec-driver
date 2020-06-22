package nec

import (
	"context"
	"fmt"
)

var (
	// PowerStatus gets the projector's power status
	PowerStatus = []byte{0x00, 0x85, 0x00, 0x00, 0x01, 0x01, 0x87}

	// PowerOn powers on the projector
	PowerOn = []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}

	// PowerStandby powers off the projector
	PowerStandby = []byte{0x02, 0x01, 0x00, 0x00, 0x00, 0x03}
)

func (p *Projector) GetPower(ctx context.Context) (bool, error) {
	resp, err := p.SendCommand(ctx, PowerStatus)
	switch {
	case err != nil:
		return false, err
	case len(resp) < 8:
		return false, fmt.Errorf("bad response from projector: 0x%x", resp)
	case resp[7] == 0b1:
		return true, nil
	default:
		return false, nil
	}
}

func (p *Projector) SetPower(ctx context.Context, power bool) error {
	var cmd []byte
	switch {
	case power:
		cmd = PowerOn
	case !power:
		cmd = PowerStandby
	default:
		return fmt.Errorf("unable to set power state to %v: must be true or false", power)
	}

	_, err := p.SendCommand(ctx, cmd)
	return err
}
