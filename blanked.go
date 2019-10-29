package nec

import (
	"context"
	"fmt"
)

var (
	// BlankMuteStatus gets the blank & mute status of the projector
	BlankMuteStatus = []byte{0x00, 0x85, 0x00, 0x00, 0x01, 0x03, 0x89}

	// Blank blanks the projector
	Blank = []byte{0x02, 0x10, 0x00, 0x00, 0x00, 0x12}

	// Unblank unblanks the projector
	Unblank = []byte{0x02, 0x11, 0x00, 0x00, 0x00, 0x13}
)

func (p *Projector) GetBlanked(ctx context.Context) (bool, error) {
	resp, err := p.SendCommand(ctx, BlankMuteStatus)
	switch {
	case err != nil:
		return false, err
	case len(resp) < 6:
		return false, fmt.Errorf("bad response from projector: 0x%x", resp)
	case resp[5] == 0b1:
		return true, nil
	default:
		return false, nil
	}
}

func (p *Projector) SetBlanked(ctx context.Context, blanked bool) error {
	cmd := Unblank
	if blanked {
		cmd = Blank
	}

	_, err := p.SendCommand(ctx, cmd)
	return err
}
