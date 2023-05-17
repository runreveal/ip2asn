package ip2asn

import (
	"net/netip"
	"os"
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupTable(t *testing.T) {
	f, err := os.Open("ip2asn-combined.tsv")
	require.NoError(t, err)
	defer f.Close()

	lt, err := New(f)
	require.NoError(t, err)

	res, ok := lt.Find(netip.MustParseAddr("1.1.1.1"))
	assert.True(t, ok, "value should exist in the database")
	expected := []Result{
		{ASN: 13335, ASOrg: "CLOUDFLARENET", CountryCode: "US"},
	}
	assert.Equal(t, expected, res, "value should exist in the database")

	res, ok = lt.Find(netip.MustParseAddr("1.1.1.1"))
	assert.True(t, ok, "value should exist in the database")

	res, ok = lt.Find(netip.MustParseAddr("2a03:2880:f11c:8183:face:b00c:0:25de"))
	assert.True(t, ok, "value should exist in the database")
	expected = []Result{{ASN: 32934, ASOrg: "FACEBOOK", CountryCode: "US"}, {ASN: 32934, ASOrg: "FACEBOOK", CountryCode: "US"}}
	assert.Equal(t, expected, res, "value should exist in the database")

	res, ok = lt.Find(netip.MustParseAddr("1.82.212.254"))
	assert.True(t, ok, "value should exist in the database")
	expected = []Result{
		{ASN: 0x1026, ASOrg: "CHINANET-BACKBONE No.31,Jin-rong Street", CountryCode: "CN"},
		{ASN: 0x20e70, ASOrg: "CHINANET-SHAANXI-CLOUD-BASE CHINANET SHAANXI province Cloud Base network", CountryCode: "CN"},
	}
	assert.Equal(t, expected, res, "value should exist in the database")
}

func BenchmarkFind(b *testing.B) {
	testCases := []string{
		"154.235.212.204",
		"2a03:2880:f11c:8183:face:b00c:0:25de",
		"204.201.100.123",
		"179.175.27.104",
		"254.198.33.71",
		"16.128.182.33",
	}

	f, err := os.Open("ip2asn-combined.tsv")
	require.NoError(b, err)
	defer f.Close()
	c, err := lru.New[netip.Addr, []Result](30)
	require.NoError(b, err)

	lookupTable, err := New(f, WithLRUCache(c))
	if err != nil {
		b.Fatalf("Error in initializing LookupTable: %v", err)
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			addr := netip.MustParseAddr(tc)
			runBench(b, lookupTable, addr)
		})
	}
}

func runBench(b *testing.B, lt *LookupTable, ip netip.Addr) {
	for i := 0; i < b.N; i++ {
		_, _ = lt.Find(ip)
	}
}
