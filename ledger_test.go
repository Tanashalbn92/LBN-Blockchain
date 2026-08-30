package main

import "testing"

func TestBlockReward(t *testing.T) {

	tests := []struct {
		height  uint64
		wantLBN uint64
	}{
		{
			height:  0,
			wantLBN: 0,
		},
		{
			height:  1,
			wantLBN: 50_000_000_000,
		},
		{
			height:  2,
			wantLBN: 50_000_000_000,
		},
		{
			height: 100_000,
			wantLBN: 50_000_000_000,
		},
		{
			height: 100_001,
			wantLBN: 25_000_000_000,
		},
		{
			height: 200_000,
			wantLBN: 25_000_000_000,
		},
		{
			height: 200_001,
			wantLBN: 12_500_000_000,
		},
		{
			height: 300_001,
			wantLBN: 6_250_000_000,
		},
	}

	for _, test := range tests {

		got := BlockReward(test.height)

		if got != test.wantLBN {

			t.Fatalf(
				"BlockReward(%d) = %d, want %d",
				test.height,
				got,
				test.wantLBN,
			)
		}
	}
}

func TestMaxSupply(t *testing.T) {

	expected := uint64(
		1_000_000_000_000_000_000,
	)

	got := MaxSupplyAtomicUnits()

	if got != expected {

		t.Fatalf(
			"MaxSupplyAtomicUnits() = %d, want %d",
			got,
			expected,
		)
	}
}

func TestSupplyCalculation(t *testing.T) {

	set := NewUTXOSet()

	set.Add(
		CreateUTXO(
			"TEST1",
			0,
			"address1",
			50_000_000_000,
		),
	)

	set.Add(
		CreateUTXO(
			"TEST2",
			0,
			"address2",
			25_000_000_000,
		),
	)

	set.Add(
		CreateUTXO(
			"TEST3",
			0,
			"address3",
			25_000_000_000,
		),
	)

	expected := uint64(
		100_000_000_000,
	)

	got := CalculateSupply(set)

	if got != expected {

		t.Fatalf(
			"CalculateSupply() = %d, want %d",
			got,
			expected,
		)
	}
}