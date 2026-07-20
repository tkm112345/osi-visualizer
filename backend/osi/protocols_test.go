package osi

import "testing"

func TestProtocolByKeyDefaultsToHTTP(t *testing.T) {
	if ProtocolByKey("").Key != "http" {
		t.Errorf("empty key should default to http, got %s", ProtocolByKey("").Key)
	}
	if ProtocolByKey("nonexistent").Key != "http" {
		t.Errorf("unknown key should default to http")
	}
}

func finalStep(steps []Step) Step { return steps[len(steps)-1] }

// DNS は UDP(8B) を使うので TCP(20B) より小さくなる。
func TestEncapsulateDNSUsesUDP(t *testing.T) {
	msg := "Hello"
	steps := Encapsulate(Request{Message: msg, Protocol: "dns"})
	want := len(msg) + udpHeaderBytes + ipHeaderBytes + ethHeaderBytes + ethTrailer
	if got := finalStep(steps).TotalBytes; got != want {
		t.Errorf("DNS final bytes = %d, want %d", got, want)
	}
	for _, s := range steps {
		if s.Level == 4 && s.HeaderBytes != udpHeaderBytes {
			t.Errorf("L4 for DNS should be UDP(%d), got %d", udpHeaderBytes, s.HeaderBytes)
		}
	}
}

// Ping(ICMP) は L4/L7 を使わず、L3 に ICMP が乗る。
func TestEncapsulatePingSkipsL4AndL7(t *testing.T) {
	steps := Encapsulate(Request{Message: "ping", Protocol: "ping"})
	for _, s := range steps {
		switch s.Level {
		case 4, 5, 6, 7:
			if s.Active {
				t.Errorf("L%d should be inactive for ping", s.Level)
			}
		case 3:
			if !contains(s.Structure, "ICMP") {
				t.Errorf("L3 structure should contain ICMP, got %s", s.Structure)
			}
		}
	}
}

func TestEncapsulatePingBytes(t *testing.T) {
	msg := "ping"
	steps := Encapsulate(Request{Message: msg, Protocol: "ping"})
	want := len(msg) + icmpHeaderBytes + ipHeaderBytes + ethHeaderBytes + ethTrailer
	if got := finalStep(steps).TotalBytes; got != want {
		t.Errorf("ping final bytes = %d, want %d", got, want)
	}
}

// 各プロトコルで encap 最終 == decap 開始（対称性）を確認。
func TestEncapDecapSymmetryAllProtocols(t *testing.T) {
	for _, p := range Protocols {
		req := Request{Message: "Hello", Protocol: p.Key}
		enc := Encapsulate(req)
		dec := Decapsulate(req)
		if finalStep(enc).TotalBytes != dec[0].TotalBytes {
			t.Errorf("%s: encap final %d != decap start %d",
				p.Key, finalStep(enc).TotalBytes, dec[0].TotalBytes)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
