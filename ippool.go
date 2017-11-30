package ippool

import (
	"fmt"
	"net"
	"sort"
)

type Pool struct {
	Entire  Range
	Remains []Range
}

type Range struct {
	First net.IP
	Last  net.IP
}

func addIP(ip net.IP, n int64) net.IP {
	res := make(net.IP, len(ip))
	copy(res, ip)
	carry := n
	for i := len(res) - 1; i > 0; i-- {
		tmp := int64(res[i]) + carry
		res[i] = byte(tmp)
		carry = tmp >> 8
		if carry == 0 {
			break
		}
	}
	if carry > 0 {
		return nil
	}
	return res
}

func prevIP(ip net.IP) net.IP {
	prev := make(net.IP, len(ip))
	copy(prev, ip)
	for i := len(prev) - 1; i > 0; i-- {
		if prev[i] > 0 {
			prev[i]--
			return prev
		}
		prev[i] = 0xff
	}
	return nil
}

func nextIP(ip net.IP) net.IP {
	next := make(net.IP, len(ip))
	copy(next, ip)
	for i := len(next) - 1; i > 0; i-- {
		if next[i] < 0xff {
			next[i]++
			return next
		}
		next[i] = 0
	}
	return nil
}

func compareIP(a, b net.IP) int {
	x := a.To16()
	y := b.To16()
	for i := 0; i < len(x); i++ {
		if x[i] < y[i] {
			return -1
		}
		if x[i] > y[i] {
			return 1
		}
	}
	return 0
}

func NewPool(first, last net.IP) *Pool {
	if compareIP(first, last) > 0 {
		return nil
	}
	return &Pool{
		Entire: Range{
			First: first.To16(),
			Last:  last.To16(),
		},
		Remains: []Range{
			Range{
				First: first.To16(),
				Last:  last.To16(),
			},
		},
	}
}

func IPv4Range(ip net.IP, size int64) Range {
	return Range{
		First: ip,
		Last:  addIP(ip, int64(1<<uint(32-size))-1),
	}
}

func (p *Pool) Clean() {
	tmp := make([]Range, 0)
	sort.Slice(p.Remains, func(i, j int) bool {
		return compareIP(p.Remains[i].First, p.Remains[j].First) < 0
	})
	for _, r := range p.Remains {
		if r.Count() > 0 {
			tmp = append(tmp, r)
		}
	}
	p.Remains = tmp
}

func (p *Pool) IsAllocated(ip Range) bool {
	if !p.Entire.Contain(ip) {
		return false
	}
	for _, r := range p.Remains {
		if r.Contain(ip) {
			return false
		}
	}
	return true
}

func (p *Pool) FindFirst(size int64) (Range, bool) {
	for _, r := range p.Remains {
		if size < r.Count() {
			return Range{
				First: r.First,
				Last:  addIP(r.First, size-1),
			}, true
		}
	}
	return Range{}, false
}

func (p *Pool) Allocate(ip Range) error {
	if !p.Entire.Contain(ip) {
		return fmt.Errorf("%s is out of pool range", ip.String())
	}
	for i, r := range p.Remains {
		if r.Contain(ip) {
			p.Remains = append(p.Remains, Range{
				First: nextIP(ip.Last),
				Last:  r.Last,
			})
			p.Remains[i].Last = prevIP(ip.First)
			p.Clean()
			return nil
		}
	}
	return fmt.Errorf("%s is already allocated", ip.String())
}

func (p *Pool) Deallocate(ip Range) error {
	if !p.Entire.Contain(ip) {
		return fmt.Errorf("%s is out of pool range", ip.String())
	}
	if !p.IsAllocated(ip) {
		return fmt.Errorf("%s is not yet allocated", ip.String())
	}
	p.Remains = append(p.Remains, ip)
	p.Clean()
	return nil
}

func (r *Range) Count() int64 {
	count := int64(0)
	for i := 0; i < net.IPv6len; i++ {
		count = (count << 8) + int64(r.Last[i]-r.First[i])
	}
	return count + 1
}

func (r *Range) Contain(ip Range) bool {
	return compareIP(ip.First, r.First) >= 0 && compareIP(r.Last, ip.Last) >= 0
}

func (r *Range) String() string {
	return fmt.Sprintf("%s-%s", r.First.String(), r.Last.String())
}
