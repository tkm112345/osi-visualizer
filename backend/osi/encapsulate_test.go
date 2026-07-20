package osi

import "testing"

func TestEncapsulateStepCount(t *testing.T) {
	steps := Encapsulate(Request{Message: "Hello"})
	if len(steps) != len(Layers) {
		t.Fatalf("want %d steps, got %d", len(Layers), len(steps))
	}
}

// TotalBytes は L7 → L1 へ進むにつれて単調増加（ヘッダが積み上がる）するはず。
func TestTotalBytesMonotonicallyIncreases(t *testing.T) {
	steps := Encapsulate(Request{Message: "Hello"})
	for i := 1; i < len(steps); i++ {
		if steps[i].TotalBytes < steps[i-1].TotalBytes {
			t.Errorf("totalBytes decreased at level %d: %d < %d",
				steps[i].Level, steps[i].TotalBytes, steps[i-1].TotalBytes)
		}
	}
}

// L1 の累積バイト数は payload + TCP20 + IP20 + Eth(14+4) になるはず。
func TestFinalTotalBytes(t *testing.T) {
	msg := "Hello" // 5 bytes
	steps := Encapsulate(Request{Message: msg})
	last := steps[len(steps)-1]
	want := len(msg) + tcpHeaderBytes + ipHeaderBytes + ethHeaderBytes + ethTrailer
	if last.TotalBytes != want {
		t.Errorf("final totalBytes = %d, want %d", last.TotalBytes, want)
	}
}

// ヘッダを付与するレイヤーには Headers が、付与しないレイヤーには付かないこと。
func TestHeaderPresenceMatchesAddsHeader(t *testing.T) {
	steps := Encapsulate(Request{Message: "Hi"})
	for _, s := range steps {
		if s.AddsHeader && s.HeaderBytes == 0 && s.Level != 7 {
			t.Errorf("level %d adds header but HeaderBytes is 0", s.Level)
		}
		if !s.AddsHeader && s.HeaderBytes != 0 {
			t.Errorf("level %d should not add header bytes, got %d", s.Level, s.HeaderBytes)
		}
	}
}

func TestDefaultIPsApplied(t *testing.T) {
	steps := Encapsulate(Request{Message: "x"})
	for _, s := range steps {
		if s.Level == 3 {
			if s.Headers["srcIp"] == "" || s.Headers["dstIp"] == "" {
				t.Errorf("L3 missing default IPs: %+v", s.Headers)
			}
		}
	}
}
