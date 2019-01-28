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
	PsnPse       *expvar.Float
	Spin         *expvar.Float
	NSpin        *expvar.Float
	Valid        *expvar.Int
	Invalid      *expvar.Int
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
	v.NSpin.Set(d.NSpin)
	v.Valid.Set(d.Valid)
	v.Invalid.Set(d.Invalid)
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
			PsnPse:       expvar.NewFloat(fmt.Sprintf("%s.delay.PsnPse", iface.Tag)),
			Spin:         expvar.NewFloat(fmt.Sprintf("%s.delay.Spin", iface.Tag)),
			NSpin:        expvar.NewFloat(fmt.Sprintf("%s.count.NSpin", iface.Tag)),
			Valid:        expvar.NewInt(fmt.Sprintf("%s.delay.Valid", iface.Tag)),
			Invalid:      expvar.NewInt(fmt.Sprintf("%s.count.Invalid", iface.Tag)),
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
	currentLine := 0
	countSame := 0
	var oldValid, oldInvalid, temp int64
	for {
		for _, iface := range cfg.Interfaces {

			// err := plus.Fetch("/root/share/vagrant_test/vagrant/spinbit_plus_printf.out", currentLine)
			// err := plus.Fetch(currentLine)
			// if err != nil {
			// 	continue
			// }

			plusData, newLine, err := plus.Parse("/root/share/vagrant_test/vagrant/mod.out", "/root/share/vagrant_test/vagrant/mod2.out", 0)
			if err != nil {
				continue
			}

			if plusData.Valid > oldValid {
				temp = plusData.Valid
				plusData.Valid = plusData.Valid - oldValid
				oldValid = temp
			} else {
				oldValid = plusData.Valid
			}

			if plusData.Invalid > oldInvalid {
				temp = plusData.Invalid
				plusData.Invalid = plusData.Invalid - oldInvalid
				oldInvalid = temp
			} else {
				oldInvalid = plusData.Invalid
			}

			if currentLine == newLine {
				countSame++
			} else {
				currentLine = newLine
				countSame = 0
			}

			if countSame == 5 {
				currentLine = 0
				countSame = 0
			}

			updatePlusExpVars(iface, plusData)

			// if currentLine > 50 {
			// 	currentLine = 0
			// }

			// fOut, err := os.OpenFile("/root/share/test_output_count.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			// if err == nil {

			// 	defer fOut.Close()

			// 	_, _ = fOut.WriteString(strconv.Itoa(currentLine))
			// 	_, _ = fOut.WriteString("\n")
			// 	_, _ = fOut.WriteString(strconv.Itoa(countSame))
			// 	_, _ = fOut.WriteString("\n")
			// }
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
