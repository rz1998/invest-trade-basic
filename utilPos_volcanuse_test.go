package tradeBasic

import (
	"testing"

	investBasic "github.com/rz1998/invest-basic"
)

// ---------------------------------------------------------------------------
// GetVolSellableFromAcPos: the can_use rule.
//
// The rule under test, in one sentence: when the counter says how much is
// sellable, that number wins -- including when it is zero -- and only an
// ABSENT field falls back to the yesterday/frozen formula.
//
// The zero case is the one that matters. On 2026-10-16 the al-013 account
// double-sold one position 12 times and had 11 of them rejected with
// 251005 (可用数量不足), because the code read a can_use of 0 as "no data"
// and re-derived 900 sellable shares from the stale yesterday/frozen pair.
// A fallback on zero is the bug, so the zero case is tested first and most
// thoroughly.
// ---------------------------------------------------------------------------

func int64Ptr(v int64) *int64 { return &v }

// A present can_use of 0 means "nothing is sellable" -- not "no data".
func TestVolSellable_CanUseZeroIsZeroNotFallback(t *testing.T) {
	pos := &SAcPos{
		UniqueCode:   "600000.SSE",
		VolTotal:     1000,
		VolYd:        900,
		VolTd:        100,
		VolFrozenYd:  0, // the formula would say 900 sellable
		VolCanUse:    int64Ptr(0),
	}
	got := GetVolSellableFromAcPos(pos)
	if got != 0 {
		t.Fatalf("can_use=0 must yield 0 sellable, got %d (the formula would have said %d)",
			got, pos.VolYd-pos.VolFrozenYd)
	}
}

// A present can_use wins over the formula even when it disagrees with it.
func TestVolSellable_CanUseWinsOverFormula(t *testing.T) {
	pos := &SAcPos{
		UniqueCode:  "600000.SSE",
		VolTotal:    1000,
		VolYd:       900,
		VolFrozenYd: 0,
		VolCanUse:   int64Ptr(700),
	}
	if got := GetVolSellableFromAcPos(pos); got != 700 {
		t.Fatalf("can_use=700 must win, got %d", got)
	}
}

// An ABSENT can_use keeps the original behaviour. This is the zero-regression
// guarantee for every other user of tradeBasic: the DBF / broker path never
// sets the field, so it must behave exactly as before.
func TestVolSellable_AbsentCanUseKeepsFormula(t *testing.T) {
	cases := []struct {
		name string
		pos  *SAcPos
		want int64
	}{
		{
			name: "sse yesterday minus frozen",
			pos: &SAcPos{
				UniqueCode:  "600000.SSE",
				VolTotal:    1000,
				VolYd:       900,
				VolFrozenYd: 100,
			},
			want: 800,
		},
		{
			name: "sse fully frozen",
			pos: &SAcPos{
				UniqueCode:  "600000.SSE",
				VolYd:       900,
				VolFrozenYd: 900,
			},
			want: 0,
		},
		{
			name: "non-sse uses total minus frozen",
			pos: &SAcPos{
				UniqueCode:     "IF2401.CFFEX",
				VolTotal:       1000,
				VolFrozenTotal: 250,
			},
			want: 750,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.pos.VolCanUse != nil {
				t.Fatalf("this case is about the absent field, but VolCanUse is set")
			}
			if got := GetVolSellableFromAcPos(c.pos); got != c.want {
				t.Fatalf("want %d, got %d", c.want, got)
			}
		})
	}
}

// A negative can_use cannot happen at the counter, but if it somehow does it
// must not turn into a negative sellable quantity -- a negative volume would
// be sent to the broker as an order.
func TestVolSellable_NegativeCanUseClampsToZero(t *testing.T) {
	pos := &SAcPos{
		UniqueCode: "600000.SSE",
		VolYd:      900,
		VolCanUse:  int64Ptr(-5),
	}
	if got := GetVolSellableFromAcPos(pos); got != 0 {
		t.Fatalf("negative can_use must clamp to 0, got %d", got)
	}
}

// The can_use rule applies to every market, not only A shares. The old
// formula had a per-market branch; the new one short-circuits before it, and
// this test makes sure that is true.
func TestVolSellable_CanUseAppliesToNonSSE(t *testing.T) {
	pos := &SAcPos{
		UniqueCode:     "IF2401.CFFEX",
		VolTotal:       1000,
		VolFrozenTotal: 0,
		VolCanUse:      int64Ptr(123),
	}
	if got := GetVolSellableFromAcPos(pos); got != 123 {
		t.Fatalf("can_use must win for non-SSE too, got %d", got)
	}
}

// A nil position is still nil-safe.
func TestVolSellable_NilPosition(t *testing.T) {
	if got := GetVolSellableFromAcPos(nil); got != 0 {
		t.Fatalf("nil position must yield 0, got %d", got)
	}
}

// The local position bookkeeping must NOT touch VolCanUse.
//
// This is the second half of the fix. handleTrade and handleCancel rewrite
// VolFrozenYd and VolYd whenever a fill or a cancel arrives. If the can_use
// value had been folded into VolFrozenYd, the first fill would have
// overwritten it and the sellable number would silently have gone back to
// being a derived quantity. VolCanUse is a separate field precisely so that
// those routines cannot reach it.
func TestVolCanUse_SurvivesLocalBookkeeping(t *testing.T) {
	pos := &SAcPos{
		UniqueCode:  "600000.SSE",
		TradeDir:    investBasic.LONG,
		VolTotal:    1000,
		VolYd:       900,
		VolTd:       100,
		VolCanUse:   int64Ptr(0),
	}
	// A sell of 100 shares arrives. handleTrade rewrites the frozen/yesterday
	// numbers; the can_use field must come out untouched.
	handleTrade(pos, Close, 100)
	if pos.VolCanUse == nil {
		t.Fatalf("handleTrade cleared VolCanUse")
	}
	if *pos.VolCanUse != 0 {
		t.Fatalf("handleTrade changed VolCanUse to %d, want 0", *pos.VolCanUse)
	}
	// And the sellable number still comes from can_use.
	if got := GetVolSellableFromAcPos(pos); got != 0 {
		t.Fatalf("sellable must still be 0, got %d", got)
	}
}

// handleCancel also must not reach the can_use field.
func TestVolCanUse_SurvivesCancel(t *testing.T) {
	pos := &SAcPos{
		UniqueCode:  "600000.SSE",
		VolTotal:    1000,
		VolYd:       900,
		VolFrozenYd: 900,
		VolCanUse:   int64Ptr(0),
	}
	handleCancel(pos, Close, 900)
	if pos.VolCanUse == nil || *pos.VolCanUse != 0 {
		t.Fatalf("handleCancel disturbed VolCanUse: %v", pos.VolCanUse)
	}
}
