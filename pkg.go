package ip2asn

import (
	"encoding/csv"
	"io"
	"net/netip"
	"strconv"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/rdleal/intervalst/interval"
)

type Result struct {
	ASN         uint
	ASOrg       string
	CountryCode string
}

type LookupTable struct {
	tree *interval.SearchTree[Result, netip.Addr]
	lru  *lru.Cache[netip.Addr, []Result]
}

type Option func(*LookupTable)

func WithLRUCache(c *lru.Cache[netip.Addr, []Result]) Option {
	return func(lt *LookupTable) {
		lt.lru = c
	}
}

func New(ip2asnTSV io.Reader, opts ...Option) (*LookupTable, error) {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(ip2asnTSV)
	csvReader.Comma = '\t'
	ret := &LookupTable{}

	for _, o := range opts {
		o(ret)
	}

	cmpFn := func(a, b netip.Addr) int {
		return a.Compare(b)
	}
	ret.tree = interval.NewSearchTree[Result](cmpFn)

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
		start := netip.MustParseAddr(rec[0])
		end := netip.MustParseAddr(rec[1])
		asn, _ := strconv.Atoi(rec[2])
		err = ret.tree.Insert(start, end, Result{
			ASN:         uint(asn),
			ASOrg:       rec[4],
			CountryCode: rec[3],
		})
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (lt LookupTable) Find(addr netip.Addr) ([]Result, bool) {
	if lt.lru != nil {
		if v, ok := lt.lru.Get(addr); ok {
			return v, ok
		}
	}
	res, ok := lt.tree.AllIntersections(addr, addr)
	if ok && lt.lru != nil {
		lt.lru.Add(addr, res)
	}
	return res, ok
}
