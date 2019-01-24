package service

import (
	"expvar"
	"fmt"
	"sync"
	"time"

	"github.com/buehlert/netem-pub/hping"
	"github.com/buehlert/netem-pub/netem"
	"github.com/buehlert/netem-pub/netemd/config"
	"github.com/buehlert/netem-pub/plus"
)

type ifaceExpVars struct {
	// add all new variables from the plus package
	// problem, in one iteration we can get multiple values from the observer output...
	// we just have to add a timestamp variable, that should be analyzed by the dashboard
	// then just update multiple times in a row?

	PktCount     *expvar.Int
	PktDropped   *expvar.Int
	PktReordered *expvar.Int
	BytesCount   *expvar.Int
	ForwardDelay *expvar.Int
	ReverseDelay *expvar.Int
	xValue       *expvar.Int
	yValue       *expvar.Int
}

type expVars struct {
	Map map[string]ifaceExpVars
	Mtx sync.Mutex
}

var ev expVars

func updateNetemExpVars(iface config.Interface, d *netem.NetemData) {
	ev.Mtx.Lock()
	defer ev.Mtx.Unlock()

	v := ev.Map[iface.Name]

	v.PktCount.Set(d.Total)
	v.PktDropped.Set(d.Dropped)
	v.PktReordered.Set(d.Reordered)
	v.BytesCount.Set(d.Bytes)
}

func updateHpingExpVars(iface config.Interface, d *hping.HpingData) {
	ev.Mtx.Lock()
	defer ev.Mtx.Unlock()

	v := ev.Map[iface.Name]

	v.ForwardDelay.Set(d.ForwardDelay)
	v.ReverseDelay.Set(d.ReverseDelay)
}

func updatePlusExpVars(iface config.Interface, d *plus.PlusData) {
	ev.Mtx.Lock()
	defer ev.Mtx.Unlock()

	v := ev.Map[iface.Name]

	v.xValue.Set(d.XData)
	v.yValue.Set(d.YData)
}

// func updatePlusExpVars()

func initExpVars(cfg *config.Config) {
	ev.Map = make(map[string]ifaceExpVars)

	for _, iface := range cfg.Interfaces {
		v := ifaceExpVars{
			PktCount:     expvar.NewInt(fmt.Sprintf("%s.packet.count", iface.Tag)),
			PktDropped:   expvar.NewInt(fmt.Sprintf("%s.packet.dropped", iface.Tag)),
			PktReordered: expvar.NewInt(fmt.Sprintf("%s.packet.reordered", iface.Tag)),
			BytesCount:   expvar.NewInt(fmt.Sprintf("%s.bytes.count", iface.Tag)),
			ForwardDelay: expvar.NewInt(fmt.Sprintf("%s.delay.forward", iface.Tag)),
			ReverseDelay: expvar.NewInt(fmt.Sprintf("%s.delay.reverse", iface.Tag)),
			xValue:       expvar.NewInt(fmt.Sprintf("%s.x.value", iface.Tag)),
			yValue:       expvar.NewInt(fmt.Sprintf("%s.y.value", iface.Tag)),
			// init plus variables
		}

		ev.Map[iface.Name] = v
	}
}

func netemPoller(cfg *config.Config) {
	for {
		for _, iface := range cfg.Interfaces {
			out, err := netem.Fetch(iface.Name)
			if err != nil {
				continue
			}

			netemData, err := netem.Parse(out)
			if err != nil {
				continue
			}

			updateNetemExpVars(iface, netemData)

		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)
	}
}

func hpingPoller(cfg *config.Config) {
	for {
		for _, iface := range cfg.Interfaces {
			out, err := hping.Fetch(iface.PingHost)
			if err != nil {
				continue
			}

			hpingData, err := hping.Parse(out)
			if err != nil {
				continue
			}

			updateHpingExpVars(iface, hpingData)

		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)
	}
}

func plusPoller(cfg *config.Config) {
	for {
		for _, iface := range cfg.Interfaces {
			out, err := plus.Fetch(iface.Name)
			if err != nil {
				continue
			}

			plusData, err := plus.Parse(out)
			if err != nil {
				continue
			}

			updatePlusExpVars(iface, plusData)

		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)
	}
}

// plusPoller

func NetemPub(cfg *config.Config, noPing bool) {
	initExpVars(cfg)
	go netemPoller(cfg)
	go plusPoller(cfg)
	if !noPing {
		go hpingPoller(cfg)
	}
	// start plusPoller
}
