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

// シリアル通信(UART/I2C/SPI)は L3〜L7 を使わず L1/L2 のみ。
func TestSerialProtocolsUseOnlyL1L2(t *testing.T) {
	for _, key := range []string{"uart", "i2c", "spi"} {
		steps := Encapsulate(Request{Message: "Hi", Protocol: key})
		for _, s := range steps {
			active := s.Level == 1 || s.Level == 2
			if s.Active != active {
				t.Errorf("%s L%d active=%v, want %v", key, s.Level, s.Active, active)
			}
		}
	}
}

// すべてのプロトコルにサンプルペイロードがある。
func TestAllProtocolsHaveSamplePayload(t *testing.T) {
	for _, p := range Protocols {
		if p.SamplePayload == "" {
			t.Errorf("protocol %s has empty SamplePayload", p.Key)
		}
	}
}

// HTTP/HTTPS のサンプルは HTML テンプレート。
func TestHTTPSamplePayloadIsHTML(t *testing.T) {
	for _, key := range []string{"http", "https"} {
		p := ProtocolByKey(key)
		if !contains(p.SamplePayload, "<!DOCTYPE html>") {
			t.Errorf("%s SamplePayload should be HTML, got %q", key, p.SamplePayload)
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
