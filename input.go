package nec

import (
	"context"
	"fmt"
	"strings"
)

var (
	// ActiveSignal whether or not the projector's input has signal or not
	ActiveSignal = []byte{0x00, 0xBF, 0x00, 0x00, 0x01, 0x02, 0xC2}

	// InputStatus returns the current input of the projector
	InputStatus = []byte{0x00, 0x85, 0x00, 0x00, 0x01, 0x02, 0x88}

	// ChangeInput changes the input of the projector
	ChangeInput = []byte{0x02, 0x03, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00}
)

func (p *Projector) GetInput(ctx context.Context) (string, error) {
	resp, err := p.SendCommand(ctx, InputStatus)
	switch {
	case err != nil:
		return "", nil
	case len(resp) < 9:
		return "", fmt.Errorf("")
	}

	input := ""
	switch {
	case resp[8] == 0x21:
		input = "hdmi"
	case resp[9] == 0x27:
		input = "hdbaset"
	}

	// add on the number of the input (e.g., hdmi1)
	return fmt.Sprintf("%s%d", input, resp[7]), nil
}

func (p *Projector) SetInput(ctx context.Context, input string) error {
	// copy the change input command
	cmd := make([]byte, len(ChangeInput))
	copy(cmd, ChangeInput)

	// stick the correct input in
	switch {
	case strings.EqualFold("hdmi1", input):
		cmd[6] = 0xA1
	case strings.EqualFold("hdmi2", input):
		cmd[6] = 0xA2
	case strings.EqualFold("hdbaset1", input):
		cmd[6] = 0xBF
	default:
		return fmt.Errorf("unknown input %q", input)
	}

	// add in the checksum
	cmd[7] = getChecksum(cmd)

	// send the command
	_, err := p.SendCommand(ctx, cmd)
	return err
}

func (p *Projector) GetActiveSignal(ctx context.Context, addr string) (bool, error) {
	// TODO
	return false, fmt.Errorf("not implemented")
}
