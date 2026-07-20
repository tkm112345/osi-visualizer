package osi

import "testing"

func TestDecapsulateStepCount(t *testing.T) {
	steps := Decapsulate(Request{Message: "Hello"})
	if len(steps) != len(Layers) {
		t.Fatalf("want %d steps, got %d", len(Layers), len(steps))
	}
}

// 受信側は L1 → L7 の順で返る（level が昇順）。
func TestDecapsulateOrderIsL1ToL7(t *testing.T) {
	steps := Decapsulate(Request{Message: "Hello"})
	for i, s := range steps {
		if s.Level != i+1 {
			t.Errorf("step %d has level %d, want %d", i, s.Level, i+1)
		}
	}
}

// ヘッダを外していくので TotalBytes は単調減少するはず。
func TestDecapsulateTotalBytesMonotonicallyDecreases(t *testing.T) {
	steps := Decapsulate(Request{Message: "Hello"})
	for i := 1; i < len(steps); i++ {
		if steps[i].TotalBytes > steps[i-1].TotalBytes {
			t.Errorf("totalBytes increased at level %d: %d > %d",
				steps[i].Level, steps[i].TotalBytes, steps[i-1].TotalBytes)
		}
	}
}

// L1 は full フレーム、L7 は payload のみに戻るはず。
func TestDecapsulateEndpoints(t *testing.T) {
	msg := "Hello"
	steps := Decapsulate(Request{Message: msg})
	full := len(msg) + tcpHeaderBytes + ipHeaderBytes + ethHeaderBytes + ethTrailer
	if steps[0].TotalBytes != full {
		t.Errorf("L1 totalBytes = %d, want %d", steps[0].TotalBytes, full)
	}
	last := steps[len(steps)-1]
	if last.TotalBytes != len(msg) {
		t.Errorf("L7 totalBytes = %d, want %d", last.TotalBytes, len(msg))
	}
}

// 送信側の最終バイト数 == 受信側の開始バイト数（対称性）。
func TestEncapDecapSymmetry(t *testing.T) {
	req := Request{Message: "Hello"}
	enc := Encapsulate(req)
	dec := Decapsulate(req)
	if enc[len(enc)-1].TotalBytes != dec[0].TotalBytes {
		t.Errorf("encap final %d != decap start %d",
			enc[len(enc)-1].TotalBytes, dec[0].TotalBytes)
	}
}
