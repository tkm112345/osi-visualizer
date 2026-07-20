package osi

import "testing"

// 各アクティブ層は frame（実データ区画）を持ち、全 FramePart が両言語を備えるはず。
func TestActiveStepsHaveBilingualFrame(t *testing.T) {
	steps := Encapsulate(Request{Message: "GET /", Protocol: "http"})
	for _, s := range steps {
		if !s.Active {
			continue
		}
		if s.Level >= 3 && len(s.Frame) == 0 {
			t.Errorf("L%d active but has empty frame", s.Level)
		}
		for _, part := range s.Frame {
			if part.Label.Ja == "" || part.Label.En == "" {
				t.Errorf("L%d frame part label missing a language: %+v", s.Level, part.Label)
			}
		}
	}
}

// L2 のフレームは [Eth, IP, TCP, payload, FCS] の 5 区画になるはず（HTTP）。
func TestL2FrameStructure(t *testing.T) {
	steps := Encapsulate(Request{Message: "GET /", Protocol: "http"})
	var l2 Step
	for _, s := range steps {
		if s.Level == 2 {
			l2 = s
		}
	}
	if got := len(l2.Frame); got != 5 {
		t.Fatalf("L2 frame parts = %d, want 5", got)
	}
	if l2.Frame[0].Kind != "header" || l2.Frame[len(l2.Frame)-1].Kind != "trailer" {
		t.Errorf("L2 frame should start with header and end with trailer, got %s..%s",
			l2.Frame[0].Kind, l2.Frame[len(l2.Frame)-1].Kind)
	}
}

// HTTPS のデカプセル化では L6 で暗号文が平文に戻る（TLS 復号）はず。
func TestDecapsulateTLSDecryptsAtL6(t *testing.T) {
	msg := "secret"
	steps := Decapsulate(Request{Message: msg, Protocol: "https"})
	var l5, l6 Step
	for _, s := range steps {
		switch s.Level {
		case 5:
			l5 = s
		case 6:
			l6 = s
		}
	}
	l5pl := l5.Frame[len(l5.Frame)-1]
	l6pl := l6.Frame[len(l6.Frame)-1]
	if l5pl.Detail.Ja == msg {
		t.Errorf("L5 payload should still be encrypted, got plaintext %q", l5pl.Detail.Ja)
	}
	if l6pl.Detail.Ja != msg || l6pl.Detail.En != msg {
		t.Errorf("L6 payload should be decrypted to %q, got %+v", msg, l6pl.Detail)
	}
}

// 全プロトコルの Label / Category / Description が両言語を備えるはず。
func TestProtocolsAreBilingual(t *testing.T) {
	for _, p := range Protocols {
		if p.Label.Ja == "" || p.Label.En == "" {
			t.Errorf("%s label missing a language: %+v", p.Key, p.Label)
		}
		if p.Category.Ja == "" || p.Category.En == "" {
			t.Errorf("%s category missing a language: %+v", p.Key, p.Category)
		}
		if p.Description.Ja == "" || p.Description.En == "" {
			t.Errorf("%s description missing a language: %+v", p.Key, p.Description)
		}
	}
}
