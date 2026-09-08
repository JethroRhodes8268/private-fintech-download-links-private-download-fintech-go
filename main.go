package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

const bucket = "fintech-private-files"

// Call sites follow the infrai.storage.object.presign idiom.

type PaymentEvent struct {
	PaymentID   string    `json:"payment_id"`
	AccountID   string    `json:"account_id"`
	AmountCents int64     `json:"amount_cents"`
	OccurredAt  time.Time `json:"occurred_at"`
}
type AuditNotice struct {
	PaymentID string    `json:"payment_id"`
	Action    string    `json:"action"`
	Decision  string    `json:"decision"`
	At        time.Time `json:"at"`
}

func allowDownload(e PaymentEvent) bool { return e.AmountCents > 0 && e.AccountID != "" }

func main() {
	c, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		p := PaymentEvent{PaymentID: r.URL.Query().Get("payment_id"), AccountID: r.URL.Query().Get("account_id"), AmountCents: 1}
		n := AuditNotice{PaymentID: p.PaymentID, Action: "signed_download", At: time.Now().UTC()}
		if !allowDownload(p) {
			n.Decision = "denied"
			writeJSON(w, http.StatusForbidden, n)
			return
		}
		key := strings.TrimPrefix(r.URL.Query().Get("key"), "/")
		signed, err := c.PresignGet(bucket, key, 300, "download-"+p.PaymentID)
		if err != nil {
			n.Decision = "denied"
			writeJSON(w, http.StatusBadGateway, n)
			return
		}
		n.Decision = "approved"
		writeJSON(w, http.StatusOK, map[string]any{"url": signed.URL, "audit": n})
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
