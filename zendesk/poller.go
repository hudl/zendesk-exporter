// Package zendesk simply provides a method which regularly calls
// the zendesk incremental ticket API
package zendesk

import (
	"github.com/adamar/ZeGo/zego"
	"time"

	"github.com/hudl/zendesk-exporter/ticketwriter"
)

const (
	// Zendesk allows this endpoint to be hit no more than 10 times/minute
	// so our minimum wait time should be 6 seconds
	minIntervalSec = 6

	// This may have to be tweaked to find a sweet spot
	maxIntervalSec = 60
)

// A Poller contains all the information necessary to regularly
// poll the zendesk API such as credentials and startTime
type Poller struct {
	Auth      zego.Auth
	StartTime string
}

// Poll continuously hits the IncrementalTicket API starting from
// StartTime and working up in batches of 1000 tickets following the
// guidelines located in Zendesk's documentation:
// https://developer.zendesk.com/rest_api/docs/core/incremental_export#tickets
func (p *Poller) Poll() {
	for {
		results, err := p.Auth.IncrementalTicket(p.StartTime)
		if err != nil {
			log.Error("Error when polling for zendesk tickets: %+v", err)
			return
		}
		log.Info("Fetched %d tickets.", results.Count)

		ticketwriter.Write(results.Tickets)

		sleepTime := time.Duration(interpSleep(float32(results.Count)))
		log.Info("Sleeping for %d seconds", sleepTime)
		time.Sleep(sleepTime * time.Second)
	}
}

// This is a vanilla interpolation between minIntervalSec and maxIntervalSec
// count is number of tickets returned (0 - 1000)
// Returns the length of time to sleep
func interpSleep(count float32) int {
	return int(maxIntervalSec + ((maxIntervalSec - minIntervalSec) * (count / (-1000))))
}
