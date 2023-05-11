package ip2asn

import (
	"encoding/csv"
	"io"
	"net/netip"
	"strconv"

	"github.com/rdleal/intervalst/interval"
)

type Result struct {
	ASN         uint
	ASOrg       string
	CountryCode string
}

type LookupTable struct {
	tree *interval.SearchTree[Result, netip.Addr]
}

func New(ip2asnTSV io.Reader) (*LookupTable, error) {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(ip2asnTSV)
	csvReader.Comma = '\t'
	ret := &LookupTable{}

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
	return lt.tree.AllIntersections(addr, addr)
}
