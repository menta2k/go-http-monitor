package notifier

import (
	"context"
	"fmt"
	"log"

	"github.com/sko/go-http-monitor/domain"
	"github.com/sko/go-http-monitor/notification"
	"github.com/sko/go-http-monitor/result"
)

type Sender interface {
	Send(ctx context.Context, target string, subject string, body string) error
}

type Notifier struct {
	notifRepo  notification.Repository
	resultRepo result.Repository
	senders    map[domain.NotificationType]Sender
}

func New(notifRepo notification.Repository, resultRepo result.Repository, senders map[domain.NotificationType]Sender) *Notifier {
	return &Notifier{
		notifRepo:  notifRepo,
		resultRepo: resultRepo,
		senders:    senders,
	}
}

// Notify checks whether a state transition occurred and sends alerts accordingly.
// - OK  -> FAIL: send alert
// - FAIL -> OK:  send recovery
// - FAIL -> FAIL: no notification (avoid flood)
// - OK  -> OK:   no notification
func (n *Notifier) Notify(ctx context.Context, m domain.Monitor, cr domain.CheckResult) {
	currentFailing := isFailing(m, cr)

	// Get the last 2 results (including the one just saved) to detect transitions
	recent, err := n.resultRepo.FindByMonitorID(ctx, m.ID, 2, 0)
	if err != nil {
		log.Printf("[notifier] failed to load recent results for monitor %d: %v", m.ID, err)
		return
	}

	// First ever check — notify only if failing
	if len(recent) < 2 {
		if currentFailing {
			n.dispatch(ctx, m, cr, eventAlert)
		}
		return
	}

	// previous is the second item (most recent is first due to DESC ordering)
	previousFailing := isFailing(m, recent[1])

	switch {
	case !previousFailing && currentFailing:
		// Transition: OK -> FAIL — send alert
		n.dispatch(ctx, m, cr, eventAlert)
	case previousFailing && !currentFailing:
		// Transition: FAIL -> OK — send recovery
		n.dispatch(ctx, m, cr, eventRecovery)
	default:
		// No transition (FAIL->FAIL or OK->OK) — no notification
	}
}

type eventType int

const (
	eventAlert    eventType = iota
	eventRecovery eventType = iota
)

func (n *Notifier) dispatch(ctx context.Context, m domain.Monitor, cr domain.CheckResult, event eventType) {
	configs, err := n.notifRepo.FindEnabledByMonitorID(ctx, m.ID)
	if err != nil {
		log.Printf("[notifier] failed to load notifications for monitor %d: %v", m.ID, err)
		return
	}

	if len(configs) == 0 {
		return
	}

	subject := buildSubject(m, cr, event)
	body := buildBody(m, cr, event)

	for _, cfg := range configs {
		sender, ok := n.senders[cfg.Type]
		if !ok {
			log.Printf("[notifier] no sender for type %q", cfg.Type)
			continue
		}
		if err := sender.Send(ctx, cfg.Target, subject, body); err != nil {
			log.Printf("[notifier] failed to send %s to %s: %v", cfg.Type, cfg.Target, err)
		} else {
			log.Printf("[notifier] sent %s %s to %s for monitor %d", eventLabel(event), cfg.Type, cfg.Target, m.ID)
		}
	}
}

func isFailing(m domain.Monitor, cr domain.CheckResult) bool {
	if cr.Error != "" {
		return true
	}
	if cr.StatusCode != m.ExpectedStatus {
		return true
	}
	if cr.BodyMatched != nil && !*cr.BodyMatched {
		return true
	}
	return false
}

func eventLabel(e eventType) string {
	if e == eventRecovery {
		return "RECOVERY"
	}
	return "ALERT"
}

func buildSubject(m domain.Monitor, cr domain.CheckResult, event eventType) string {
	if event == eventRecovery {
		return fmt.Sprintf("[RECOVERED] Monitor %s is back to normal", m.URL)
	}
	if cr.Error != "" {
		return fmt.Sprintf("[ALERT] Monitor %s - Error", m.URL)
	}
	if cr.StatusCode != m.ExpectedStatus {
		return fmt.Sprintf("[ALERT] Monitor %s - Status %d (expected %d)", m.URL, cr.StatusCode, m.ExpectedStatus)
	}
	return fmt.Sprintf("[ALERT] Monitor %s - Body mismatch", m.URL)
}

func buildBody(m domain.Monitor, cr domain.CheckResult, event eventType) string {
	var header string
	if event == eventRecovery {
		header = "RECOVERED - Monitor is healthy again\n\n"
	} else {
		header = "ALERT - Monitor check failed\n\n"
	}

	msg := header + fmt.Sprintf("Monitor: %s\nExpected Status: %d\nActual Status: %d\nResponse Time: %dms\n",
		m.URL, m.ExpectedStatus, cr.StatusCode, cr.ResponseTimeMs)

	if cr.Error != "" {
		msg += fmt.Sprintf("Error: %s\n", cr.Error)
	}
	if cr.BodyMatched != nil {
		msg += fmt.Sprintf("Body Match: %v (looking for %q)\n", *cr.BodyMatched, m.BodyContains)
	}
	msg += fmt.Sprintf("Checked At: %s\n", cr.CheckedAt.Format("2006-01-02 15:04:05 UTC"))
	return msg
}
