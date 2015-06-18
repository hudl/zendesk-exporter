// Package zendesk simply provides a method which regularly calls
// the zendesk incremental ticket API
package zendesk

import (
	"fmt"
	"github.com/hudl/ZeGo/zego"
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
	TickWrt   ticketwriter.TicketWriter

	// This is to prevent duplicates during low activity times
	PrevTickId uint64
}

// Poll continuously hits the IncrementalTicket API starting from
// StartTime and working up in batches of 1000 tickets following the
// guidelines located in Zendesk's documentation:
// https://developer.zendesk.com/rest_api/docs/core/incremental_export#tickets
func (p *Poller) Poll() {
	for {
		results, err := p.Auth.IncrementalTicket(p.StartTime)
		if err != nil {
			log.Error("Error when polling for zendesk tickets with StartTime=%s: %+v", p.StartTime, err)
			return
		}
		log.Info("Fetched %d tickets.", results.Count)

		if results.Count == 0 || results.Tickets[len(results.Tickets)-1].Id == p.PrevTickId || results.Next_page == "" {
			log.Info("Sleeping for 5 minutes then using the same start time")
			time.Sleep(5 * time.Minute)
			continue
		}

		p.TickWrt.Write(results.Tickets, p.StartTime)
		p.PrevTickId = results.Tickets[len(results.Tickets)-1].Id
		p.StartTime = fmt.Sprintf("%d", results.EndTime)
		startTime := time.Unix(int64(results.EndTime), 0)
		log.Info("Next start time is: %s (%s)", p.StartTime, startTime.Format("2006-01-02 15:04:05"))

		//And we sleep for a reasonable amount of time
		sleepTime := time.Duration(interpSleep(float32(results.Count))) * time.Second
		log.Info("Sleeping for %f seconds", sleepTime.Seconds())
		time.Sleep(sleepTime)
	}
}

// This is a vanilla interpolation between minIntervalSec and maxIntervalSec
// count is number of tickets returned (0 - 1000)
// Returns the length of time to sleep
func interpSleep(count float32) int {
	return int(maxIntervalSec + ((maxIntervalSec - minIntervalSec) * (count / (-1000))))
}
