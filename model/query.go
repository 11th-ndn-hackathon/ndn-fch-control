package model

import (
	"net/url"
	"sort"
	"strconv"

	"github.com/pkg/math"
)

// Query represents an API query.
type Query struct {
	Count     int
	Transport TransportType
	IPv4      bool
	IPv6      bool
	Position  LonLat
}

func (q Query) match(router RouterAvail) bool {
	return (q.IPv4 && router.Available[TransportIPFamily{q.Transport, IPv4}]) ||
		(q.IPv6 && router.Available[TransportIPFamily{q.Transport, IPv6}])
}

// Execute executes a query.
func (q Query) Execute(avail []RouterAvail) (res []RouterAvail) {
	for _, router := range avail {
		if q.match(router) {
			res = append(res, router)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return Distance(q.Position, res[i].Position) < Distance(q.Position, res[j].Position)
	})

	if len(res) > q.Count {
		res = res[:q.Count]
	}
	return res
}

// ParseQueries constructs a list of Query from URL query string.
func ParseQueries(qs string) (list []Query) {
	v, _ := url.ParseQuery(qs)

	q := Query{
		Count:     1,
		Transport: TransportUDP,
		IPv4:      v.Get("ipv4") != "0",
		IPv6:      v.Get("ipv6") != "0",
	}
	if k, e := strconv.ParseUint(v.Get("k"), 10, 32); e == nil {
		q.Count = math.MaxInt(int(k), q.Count)
	}
	q.Position[0], _ = strconv.ParseFloat(v.Get("lon"), 64)
	q.Position[1], _ = strconv.ParseFloat(v.Get("lat"), 64)

	for _, tr := range v["cap"] {
		q.Transport = TransportType(tr)
		list = append(list, q)
	}
	if len(list) == 0 {
		list = append(list, q)
	}
	return list
}