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
	PktCount     *expvar.Int
	PktDropped   *expvar.Int
	PktReordered *expvar.Int
	BytesCount   *expvar.Int
	ForwardDelay *expvar.Int
	ReverseDelay *expvar.Int
	PsnPse       *expvar.Float
	Spin         *expvar.Float
	Valid        *expvar.Int
	Invalid      *expvar.Int
	ValidTs      *expvar.Int
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

	v.PsnPse.Set(d.PsnPse)
	v.Spin.Set(d.Spin)
	v.Valid.Set(d.Valid)
	v.Invalid.Set(d.Invalid)
	v.ValidTs.Set(d.ValidTs)
}

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
			PsnPse:       expvar.NewFloat(fmt.Sprintf("%s.delay.PsnPse", iface.Tag)),
			Spin:         expvar.NewFloat(fmt.Sprintf("%s.delay.Spin", iface.Tag)),
			Valid:        expvar.NewInt(fmt.Sprintf("%s.delay.Valid", iface.Tag)),
			Invalid:      expvar.NewInt(fmt.Sprintf("%s.count.Invalid", iface.Tag)),
			ValidTs:      expvar.NewInt(fmt.Sprintf("%s.count.ValidTs", iface.Tag)),
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
	oldValid := []int64{0, 0}
	oldInvalid := []int64{0, 0}
	oldValidTs := []int64{0, 0}
	temp := []int64{0, 0}

	dataFile := "/root/share/vagrant_test/vagrant/rtt_data.out"
	counterFile := "/root/share/vagrant_test/vagrant/rtt_counters.out"

	for {
		for i, iface := range cfg.Interfaces {

			plusData, err := plus.Parse(dataFile, counterFile)
			if err != nil {
				continue
			}

			if plusData.Valid == 0 || plusData.PsnPse == 0 {
				continue
			}

			temp[i] = plusData.Valid
			plusData.Valid = plusData.Valid - oldValid[i]
			oldValid[i] = temp[i]

			temp[i] = plusData.Invalid
			plusData.Invalid = plusData.Invalid - oldInvalid[i]
			oldInvalid[i] = temp[i]

			temp[i] = plusData.ValidTs
			plusData.ValidTs = plusData.ValidTs - oldValidTs[i]
			oldValidTs[i] = temp[i]

			updatePlusExpVars(iface, plusData)
		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)
	}
}

func NetemPub(cfg *config.Config, noPing bool) {
	initExpVars(cfg)
	go netemPoller(cfg)
	go plusPoller(cfg)
	if !noPing {
		go hpingPoller(cfg)
	}
}
