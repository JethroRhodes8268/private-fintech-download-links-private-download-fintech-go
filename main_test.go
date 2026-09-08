package main

import "testing"

func TestAllowDownload(t *testing.T) {
	tests := []struct {
		name  string
		event PaymentEvent
		want  bool
	}{
		{"funded account", PaymentEvent{AccountID: "acct-7", AmountCents: 2500}, true},
		{"zero payment", PaymentEvent{AccountID: "acct-7", AmountCents: 0}, false},
		{"missing account", PaymentEvent{AmountCents: 2500}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowDownload(tt.event); got != tt.want {
				t.Fatalf("allowDownload()=%v, want %v", got, tt.want)
			}
		})
	}
}
