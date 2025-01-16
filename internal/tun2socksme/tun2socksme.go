package tun2socksme

import (
	"fmt"
	"strings"
	"tun2socksme/internal/tun"
	shell "tun2socksme/pkg"
)

type Gateway struct {
	device  string
	address string
}

type Tun2socksme struct {
	tun         tun.Tun
	defgate     *Gateway
	excludenets []string
	metric      string
}

func New(
	_tun tun.Tun,
	_excludenets []string,
	_metric int,
) Tun2socksme {
	return Tun2socksme{
		tun:         _tun,
		excludenets: _excludenets,
		metric:      fmt.Sprint(_metric),
	}
}

func (t *Tun2socksme) Run() error {
	if err := t.tun.Run(); err != nil {
		return fmt.Errorf("run tun2socks error: %w", err)
	}
	if err := t.setDefGate(); err != nil {
		return fmt.Errorf("gateway error: %w", err)
	}
	if err := t.setExcludeNets(); err != nil {
		return fmt.Errorf("route error: %w", err)
	}
	if err := t.setDefGateToTun(); err != nil {
		return fmt.Errorf("default route to proxy error: %w", err)
	}
	return nil
}

func (t *Tun2socksme) setDefGate() error {
	out, err := shell.New("ip", "ro", "sh").Run()
	if err != nil {
		return fmt.Errorf("failed to get default gateway: %w", err)
	}
	s := strings.Fields(strings.TrimSpace(string(out)))
	if len(s) < 6 {
		return fmt.Errorf("failed to get default gateway")
	}
	t.defgate = &Gateway{
		address: s[2],
		device:  s[4],
	}
	return nil
}

func (t *Tun2socksme) setExcludeNets() error {
	if _, err := shell.New("ip", "ro", "add", t.tun.Host, "via", t.defgate.address, "dev", t.defgate.device).Run(); err != nil {
		return fmt.Errorf("failed to set route %s via %s", t.tun.Host, t.defgate.device)
	}
	for _, net := range t.excludenets {
		if _, err := shell.New("ip", "ro", "add", net, "via", t.defgate.address, "dev", t.defgate.device).Run(); err != nil {
			return fmt.Errorf("failed to set route %s via %s", net, t.defgate.device)
		}
	}
	return nil
}

func (t *Tun2socksme) setDefGateToTun() error {
	if _, err := shell.New("ip", "link", "set", t.tun.Device, "up").Run(); err != nil {
		return fmt.Errorf("failed to up %s device: %w", t.tun.Device, err)
	}
	if _, err := shell.New("ip", "route", "add", "default", "dev", t.tun.Device, "metric", t.metric).Run(); err != nil {
		return fmt.Errorf("failed to set default route via %s: %w", t.tun.Device, err)
	}
	return nil
}
