package ippool

import (
	"net"
	"testing"
)

func TestAddIP(t *testing.T) {
	ip := net.IPv4(1, 2, 3, 4)
	newIP := addIP(ip, 300)
	if !newIP.Equal(net.IPv4(1, 2, 4, 48)) {
		t.Error(newIP)
	}
}

func TestRangeCount(t *testing.T) {
	r := Range{
		First: net.IPv4(1, 2, 0, 0),
		Last:  net.IPv4(1, 3, 0, 0),
	}
	expect := int64(65537)
	count := r.Count()
	if count != expect {
		t.Errorf("Range.Count() is %d (expect=%d)", expect, count)
	}
}

func TestRangeContainTrue(t *testing.T) {
	r := Range{
		First: net.IPv4(1, 2, 0, 0),
		Last:  net.IPv4(1, 3, 0, 0),
	}
	if !r.Contain(Range{
		net.IPv4(1, 2, 0, 0),
		net.IPv4(1, 2, 0, 4),
	}) {
		t.Error("Range.Contain() is not true")
	}
}

func TestRangeContainFalse(t *testing.T) {
	r := Range{
		First: net.IPv4(1, 2, 0, 0),
		Last:  net.IPv4(1, 3, 0, 0),
	}
	if r.Contain(Range{
		net.IPv4(1, 3, 0, 0),
		net.IPv4(1, 3, 0, 1),
	}) {
		t.Error("Range.Contain() is not false")
	}
}

func TestPrevIP(t *testing.T) {
	ip := net.IPv4(1, 0, 0, 0)
	if !prevIP(ip).Equal(net.IPv4(0, 255, 255, 255)) {
		t.Error()
	}
}

func TestNextIP(t *testing.T) {
	ip := net.IPv4(0, 255, 255, 255)
	if !nextIP(ip).Equal(net.IPv4(1, 0, 0, 0)) {
		t.Error()
	}
}

func TestRangeAllocate(t *testing.T) {
	pool := NewPool(
		net.IPv4(10, 224, 0, 0),
		net.IPv4(10, 224, 255, 255),
	)
	ip := Range{
		net.IPv4(10, 224, 100, 1),
		net.IPv4(10, 224, 100, 4),
	}
	if err := pool.Allocate(ip); err != nil {
		t.Error(err)
	}
	if !pool.IsAllocated(ip) {
		t.Error()
	}
}

func TestRangeDeallocate(t *testing.T) {
	pool := NewPool(
		net.IPv4(10, 224, 0, 0),
		net.IPv4(10, 224, 255, 255),
	)
	ip := Range{
		net.IPv4(10, 224, 100, 1),
		net.IPv4(10, 224, 100, 4),
	}
	if err := pool.Allocate(ip); err != nil {
		t.Error(err)
	}
	if err := pool.Deallocate(ip); err != nil {
		t.Error(err)
	}
	if pool.IsAllocated(ip) {
		t.Error()
	}
}

func TestIPv4Range(t *testing.T) {
	r := IPv4Range(net.IPv4(192, 168, 1, 0), 24)
	if !r.First.Equal(net.IPv4(192, 168, 1, 0)) || !r.Last.Equal(net.IPv4(192, 168, 1, 255)) {
		t.Error(r)
	}
}
