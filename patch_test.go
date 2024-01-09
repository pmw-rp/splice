package splice

import (
	"bytes"
	"testing"
)

func ptr(data []byte) *[]byte {
	return &data
}

func TestReindex(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foo")), ptr([]byte("bar"))}
	splice := Splice{
		slices: slices,
		len:    6,
	}

	//_, d := splice.getPosition(-1)
	//if !(d != nil) {
	//	t.Fail()
	//}

	p, d := splice.getPosition(0)
	if !(p.slice == 0 && p.offset == 0 && d == nil) {
		t.Fail()
	}

	p, d = splice.getPosition(1)
	if !(p.slice == 0 && p.offset == 1 && d == nil) {
		t.Fail()
	}

	p, d = splice.getPosition(2)
	if !(p.slice == 0 && p.offset == 2 && d == nil) {
		t.Fail()
	}

	p, d = splice.getPosition(3)
	if !(p.slice == 1 && p.offset == 0 && d == nil) {
		t.Fail()
	}

	p, d = splice.getPosition(4)
	if !(p.slice == 1 && p.offset == 1 && d == nil) {
		t.Fail()
	}

	p, d = splice.getPosition(5)
	if !(p.slice == 1 && p.offset == 2 && d == nil) {
		t.Fail()
	}

	//p, d = splice.getPosition(6)
	//if !(d != nil) {
	//	t.Fail()
	//}
}

func TestCompact(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foo")), ptr([]byte("bar"))}
	splice := Splice{
		slices: slices,
		len:    6,
	}
	compacted := splice.Compact()
	expected := []byte("foobar")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestAppend(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foo"))}
	splice := Splice{
		slices: slices,
		len:    3,
	}
	appendage := ptr([]byte("bar"))
	err := splice.Append(appendage)
	if err != nil {
		t.Fail()
	}
	compacted := splice.Compact()
	expected := []byte("foobar")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestInsertAsAppend(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foo"))}
	splice := Splice{
		slices: slices,
		len:    3,
	}
	appendage := ptr([]byte("bar"))
	err := splice.Insert(appendage, 3)
	if err != nil {
		t.Fail()
	}
	compacted := splice.Compact()
	expected := []byte("foobar")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestPrepend(t *testing.T) {
	slices := []*[]byte{ptr([]byte("bar"))}
	splice := Splice{
		slices: slices,
		len:    3,
	}
	item := ptr([]byte("foo"))
	err := splice.Prepend(item)
	if err != nil {
		t.Fail()
	}
	compacted := splice.Compact()
	expected := []byte("foobar")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestInsertAsPrepend(t *testing.T) {
	slices := []*[]byte{ptr([]byte("bar"))}
	splice := Splice{
		slices: slices,
		len:    3,
	}
	item := ptr([]byte("foo"))
	err := splice.Insert(item, 0)
	if err != nil {
		t.Fail()
	}
	compacted := splice.Compact()
	expected := []byte("foobar")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestInsertSplit(t *testing.T) {
	slices := []*[]byte{ptr([]byte("abcghi"))}
	splice := Splice{
		slices: slices,
		len:    6,
	}

	// Pre Insertion Checks
	if splice.CountSlices() != 1 {
		t.Fail()
	}
	if splice.Length() != 6 {
		t.Fail()
	}
	if !bytes.Equal(splice.Compact(), []byte("abcghi")) {
		t.Fail()
	}

	// Insertion
	insertion := ptr([]byte("def"))
	err := splice.Insert(insertion, 3)
	if err != nil {
		t.Fail()
	}

	// Post Insertion Checks
	if !bytes.Equal(splice.Compact(), []byte("abcdefghi")) {
		t.Fail()
	}
	if splice.Length() != 9 {
		t.Fail()
	}
	if splice.CountSlices() != 3 {
		t.Fail()
	}
}

func TestInsertBetween(t *testing.T) {
	slices := []*[]byte{ptr([]byte("abc")), ptr([]byte("ghi"))}
	splice := Splice{
		slices: slices,
		len:    6,
	}

	// Pre Insertion Checks
	if splice.CountSlices() != 2 {
		t.Fail()
	}
	if splice.Length() != 6 {
		t.Fail()
	}
	if !bytes.Equal(splice.Compact(), []byte("abcghi")) {
		t.Fail()
	}

	// Insertion
	insertion := ptr([]byte("def"))
	err := splice.Insert(insertion, 3)
	if err != nil {
		t.Fail()
	}

	// Post Insertion Checks
	if !bytes.Equal(splice.Compact(), []byte("abcdefghi")) {
		t.Fail()
	}
	if splice.Length() != 9 {
		t.Fail()
	}
	if splice.CountSlices() != 3 {
		t.Fail()
	}
}

func TestDeleteFirstSlice(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foo")), ptr([]byte("bar"))}
	splice := Splice{
		slices: slices,
		len:    6,
	}
	err := splice.Delete(3, 3)
	if err != nil {
		t.Fail()
		return
	}
	compacted := splice.Compact()
	expected := []byte("foo")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
	if splice.len != 3 {
		t.Fail()
	}
}

func TestDeleteMiddle(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foobarbaz"))}
	splice := Splice{
		slices: slices,
		len:    9,
	}
	err := splice.Delete(3, 3)
	if err != nil {
		t.Fail()
		return
	}
	if splice.CountSlices() != 2 {
		t.Fail()
		return
	}
	if splice.len != 6 {
		t.Fail()
	}
	compacted := splice.Compact()
	expected := []byte("foobaz")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestDeleteMiddleSlice(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foo")), ptr([]byte("bar")), ptr([]byte("baz"))}
	splice := Splice{
		slices: slices,
		len:    9,
	}
	err := splice.Delete(3, 3)
	if err != nil {
		t.Fail()
		return
	}
	compacted := splice.Compact()
	expected := []byte("foobaz")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
	if splice.len != 6 {
		t.Fail()
	}
	if splice.CountSlices() != 2 {
		t.Fail()
	}
}

func TestRegion(t *testing.T) {
	region := Region{
		start: Position{0, 0},
		end:   Position{0, 2},
	}

	if !contains(region, Position{
		slice:  0,
		offset: 0,
	}) {
		t.Fail()
	}

	if !contains(region, Position{
		slice:  0,
		offset: 1,
	}) {
		t.Fail()
	}

	if !contains(region, Position{
		slice:  0,
		offset: 2,
	}) {
		t.Fail()
	}

	if contains(region, Position{
		slice:  0,
		offset: 3,
	}) {
		t.Fail()
	}

	if contains(region, Position{
		slice:  1,
		offset: 0,
	}) {
		t.Fail()
	}

	if contains(region, Position{
		slice:  -1,
		offset: 0,
	}) {
		t.Fail()
	}
}

func TestDisjointOverlap(t *testing.T) {
	_, ok := overlap(
		Region{start: Position{0, 0}, end: Position{slice: 0, offset: 2}},
		Region{start: Position{1, 0}, end: Position{slice: 1, offset: 2}},
	)
	if ok {
		t.Fail()
	}
}

func TestRightOverlap(t *testing.T) {
	r, ok := overlap(
		Region{start: Position{0, 0}, end: Position{slice: 0, offset: 2}},
		Region{start: Position{0, 1}, end: Position{slice: 1, offset: 2}},
	)
	expected := Region{
		start: Position{0, 1},
		end:   Position{0, 2},
	}
	if !ok {
		t.Fail()
	}
	if !r.equals(expected) {
		t.Fail()
	}
}

func TestHead(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foobarbaz"))}
	splice := Splice{
		slices: slices,
		len:    9,
	}
	head, err := splice.Head(6)
	if err != nil {
		t.Fail()
		return
	}

	if splice.CountSlices() != 1 {
		t.Fail()
		return
	}
	if splice.len != 9 {
		t.Fail()
	}

	if head.CountSlices() != 1 {
		t.Fail()
		return
	}
	if head.len != 6 {
		t.Fail()
	}

	compacted := head.Compact()
	expected := []byte("foobar")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestTail(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foobarbaz"))}
	splice := Splice{
		slices: slices,
		len:    9,
	}
	tail, err := splice.Tail(6)
	if err != nil {
		t.Fail()
		return
	}

	if splice.CountSlices() != 1 {
		t.Fail()
		return
	}
	if splice.len != 9 {
		t.Fail()
	}

	if tail.CountSlices() != 1 {
		t.Fail()
		return
	}
	if tail.len != 6 {
		t.Fail()
	}

	compacted := tail.Compact()
	expected := []byte("barbaz")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestMiddle(t *testing.T) {
	slices := []*[]byte{ptr([]byte("foobarbaz"))}
	splice := Splice{
		slices: slices,
		len:    9,
	}
	middle, err := splice.Middle(2, 6)
	if err != nil {
		t.Fail()
		return
	}
	if splice.CountSlices() != 1 {
		t.Fail()
		return
	}
	if splice.len != 9 {
		t.Fail()
	}

	if middle.CountSlices() != 1 {
		t.Fail()
		return
	}
	if middle.len != 6 {
		t.Fail()
	}

	compacted := middle.Compact()
	expected := []byte("obarba")
	if !bytes.Equal(compacted, expected) {
		t.Fail()
	}
}

func TestIteration(t *testing.T) {
	expected := "foobarbaz"
	slices := []*[]byte{ptr([]byte("foo")), ptr([]byte("bar")), ptr([]byte("baz"))}
	splice := Splice{
		slices: slices,
		len:    9,
	}
	it := splice.Iterate()

	for _, c := range expected {
		if it.Next() {
			_, value, err := it.Get()
			if err != nil {
				t.Fail()
				return
			}
			if value != byte(c) {
				t.Fail()
				return
			}
		}
	}

	if it.Next() {
		t.Fail()
	}
}
