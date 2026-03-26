package memory

import (
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Ticker struct {
	core   cacheCore
	ticker *time.Ticker
	stopCh chan struct{}
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newTicker(core cacheCore) Ticker {
	return Ticker{
		core:   core,
		ticker: nil,
		stopCh: nil,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (t Ticker) Start() {
	if t.ticker != nil || t.stopCh != nil {
		t.Stop()
	}

	t.ticker = time.NewTicker(tickerPeriod)
	t.stopCh = make(chan struct{})

	for {
		select {
		case now := <-t.ticker.C:
			prevTickerPoint := calcTickerPointPrev(now)
			keys, _ := t.core.tickerPointToKeys.DeleteAndGet(prevTickerPoint)
			for k := range keys {
				t.core.m.Delete(k)
			}
		case <-t.stopCh:
			return
		}
	}
}

func (t Ticker) Stop() {
	if t.stopCh != nil {
		close(t.stopCh)
		t.stopCh = nil
	}

	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------
