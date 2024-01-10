package splice

import (
	"errors"
	"fmt"
)

type Splice struct {
	slices []*[]byte
	len    int
}

func NewSplice(data []byte) *Splice {
	return &Splice{slices: []*[]byte{&data}, len: len(data)}
}

type InsertType int

const (
	Split InsertType = iota
	Between
	Illegal
)

func (s *Splice) getPosition(index int) (position Position, err error) {

	//if index < 0 || index >= s.Length() {
	//	return Position{}, errors.New(fmt.Sprintf("illegal index %d", index))
	//}

	currentPosition := 0
	currentSlice := 0

	for _, slice := range s.slices {
		nextSlice := currentPosition + len(*slice)
		if index < nextSlice {
			return Position{currentSlice, index - currentPosition}, nil
		}
		if index == nextSlice {
			return Position{currentSlice + 1, 0}, nil
		}
		currentPosition = nextSlice
		currentSlice = currentSlice + 1
	}
	return Position{}, errors.New("oops")
}

func (s *Splice) validateSliceNumber(n int) bool {
	return n >= 0 && n < s.CountSlices()
}

func (s *Splice) getInsertType(position Position) InsertType {
	if s.validateSliceNumber(position.slice) {
		if position.slice > 0 && position.offset == 0 {
			return Between
		} else if position.offset > 0 {
			return Split
		}

	}
	return Illegal
}

func (s *Splice) Prepend(data *[]byte) error {
	var result []*[]byte
	result = append(result, data)
	result = append(result, s.slices...)
	s.slices = result
	s.len = s.len + len(*data)
	return nil
}

func (s *Splice) Insert(data *[]byte, position int) error {
	if position == s.len {
		// The position indicates append, so don't calculate the index (which would be illegal), just do it
		err := s.Append(data)
		if err != nil {
			return err
		}
	} else if position == 0 {
		// The position indicates prepend, so just do it
		err := s.Prepend(data)
		if err != nil {
			return err
		}
	} else {
		position, err := s.getPosition(position)
		if err != nil {
			return err
		}
		switch s.getInsertType(position) {
		case Between:
			pre := s.slices[0:position.slice]
			post := s.slices[position.slice:len(s.slices)]
			var result []*[]byte
			result = append(result, pre...)
			result = append(result, data)
			result = append(result, post...)
			s.len = s.len + len(*data)
			s.slices = result
		case Split:
			pre := s.slices[0:position.slice]
			var splits []*[]byte
			before := (*s.slices[position.slice])[0:position.offset]
			after := (*s.slices[position.slice])[position.offset:len(*s.slices[position.slice])]
			splits = append(splits, &before)
			splits = append(splits, data)
			splits = append(splits, &after)
			post := s.slices[position.slice+1 : len(s.slices)]
			var result []*[]byte
			result = append(result, pre...)
			result = append(result, splits...)
			result = append(result, post...)
			s.len = s.len + len(*data)
			s.slices = result
		}
	}
	return nil
}

func (s *Splice) Append(data *[]byte) error {
	var result []*[]byte
	result = append(result, s.slices...)
	result = append(result, data)
	s.slices = result
	s.len = s.len + len(*data)
	return nil
}

//func (s *Splice) Replace(data []byte, position int) error {
//	return nil
//}

type Position struct {
	slice  int
	offset int
}

func (p *Position) equals(a Position) bool {
	return p.slice == a.slice && p.offset == a.offset
}

func compare(a Position, b Position) int {
	if a.slice < b.slice {
		return -1
	} else if a.slice > b.slice {
		return 1
	} else {
		if a.offset < b.offset {
			return -1
		} else if a.offset > b.offset {
			return 1
		} else {
			return 0
		}
	}
}

func max(a Position, b Position) Position {
	if a.slice < b.slice {
		return b
	} else if a.slice > b.slice {
		return a
	} else {
		if a.offset < b.offset {
			return b
		} else if a.offset > b.offset {
			return a
		} else {
			return a
		}
	}
}

func min(a Position, b Position) Position {
	m := max(a, b)
	if a.equals(m) {
		return b
	} else {
		return a
	}
}

type Region struct {
	start Position
	end   Position
}

func valid(region Region) bool {
	return compare(region.start, region.end) < 0
}

func (r *Region) equals(a Region) bool {
	return r.start.equals(a.start) && r.end.equals(a.end)
}

func contains(region Region, position Position) bool {
	return compare(position, region.start) >= 0 && compare(position, region.end) <= 0
}

func overlap(a Region, b Region) (Region, bool) {
	r := Region{start: max(a.start, b.start), end: min(a.end, b.end)}
	return r, valid(r)
}

func (s *Splice) getSliceRegion(sliceIndex int) (Region, error) {
	maxSliceIndex := len(s.slices) - 1
	if sliceIndex < 0 || sliceIndex > maxSliceIndex {
		return Region{}, errors.New("oops")
	}
	return Region{start: Position{slice: sliceIndex, offset: 0}, end: Position{slice: sliceIndex + 1, offset: 0}}, nil
}

type SliceAction int

const (
	Keep SliceAction = iota
	KeepHead
	DropMiddle
	Drop
	KeepTail
	Unknown
)

func (s *Splice) getAction(slice int, deletion Region) (SliceAction, int, int) {
	region, err := s.getSliceRegion(slice)
	if err != nil {
		return Unknown, -1, -1
	}
	overlappingRegion, ok := overlap(region, deletion)
	if !ok {
		return Keep, -1, -1
	}
	if overlappingRegion.start.equals(region.start) && overlappingRegion.end.equals(region.end) {
		return Drop, -1, -1
	}
	if !overlappingRegion.start.equals(region.start) && overlappingRegion.end.equals(region.end) {
		return KeepHead, overlappingRegion.start.offset, -1
	}
	if overlappingRegion.start.equals(region.start) && !overlappingRegion.end.equals(region.end) {
		return KeepTail, -1, overlappingRegion.end.offset
	}
	if !overlappingRegion.start.equals(region.start) && !overlappingRegion.end.equals(region.end) {
		return DropMiddle, overlappingRegion.start.offset, overlappingRegion.end.offset
	}

	return Unknown, -1, -1
}

func (s *Splice) Delete(index int, length int) error {
	if index < 0 || index+length > s.len {
		return errors.New(fmt.Sprintf("illegal index and/or length: %d", index))
	}
	deletionStart, err := s.getPosition(index)
	if err != nil {
		return err
	}
	deletionEnd, err := s.getPosition(index + length)
	if err != nil {
		return err
	}

	deletionRegion := Region{deletionStart, deletionEnd}

	bytesDeleted := 0
	var result []*[]byte
	for i := range s.slices {
		action, lower, upper := s.getAction(i, deletionRegion)
		if action == Keep {
			result = append(result, s.slices[i])
		}
		if action == KeepHead {
			slice := *s.slices[i]
			var head []byte
			head = append(head, slice[0:lower]...)
			result = append(result, &head)
			bytesDeleted = bytesDeleted + (len(slice) - lower)
		}
		if action == KeepTail {
			slice := *s.slices[i]
			var tail []byte
			tail = append(slice[upper:], tail...)
			result = append(result, &tail)
			bytesDeleted = bytesDeleted + upper
		}
		if action == DropMiddle {
			slice := *s.slices[i]
			var head []byte
			head = append(head, slice[0:lower]...)
			var tail []byte
			tail = append(tail, slice[upper:]...)
			result = append(result, &head)
			result = append(result, &tail)
			bytesDeleted = bytesDeleted + (upper - lower)
		}
		if action == Drop {
			slice := *s.slices[i]
			bytesDeleted = bytesDeleted + len(slice)
		}
	}
	s.slices = result
	s.len = s.len - bytesDeleted
	return nil
}

func (s *Splice) clone() *Splice {
	result := Splice{
		slices: make([]*[]byte, len(s.slices)),
		len:    s.len,
	}
	for i, slice := range s.slices {
		clonedSlice := make([]byte, len(*slice))
		copy(clonedSlice, *slice)
		result.slices[i] = &clonedSlice
	}
	return &result
}

func (s *Splice) Head(index int) (*Splice, error) {
	result := s.clone()
	err := result.Delete(index, s.len-index)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Splice) HeadUnsafe(index int) *Splice {
	result, err := s.Head(index)
	if err != nil {
		panic(1)
	}
	return result
}

func (s *Splice) Tail(index int) (*Splice, error) {
	result := s.clone()
	err := result.Delete(0, s.len-index)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Splice) TailUnsafe(index int) *Splice {
	result, err := s.Tail(index)
	if err != nil {
		panic(1)
	}
	return result
}

func (s *Splice) Middle(index int, length int) (*Splice, error) {
	result := s.clone()
	err := result.Delete(0, index)
	if err != nil {
		return nil, err
	}
	err = result.Delete(length, result.Length()-length)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Splice) MiddleUnsafe(index int, length int) *Splice {
	result, err := s.Middle(index, length)
	if err != nil {
		panic(1)
	}
	return result
}

func (s *Splice) Compact() []byte {
	value := make([]byte, s.len)
	p := 0
	for _, slice := range s.slices {
		copy(value[p:p+len(*slice)], *slice)
		p = p + len(*slice)
	}
	return value
}

func (s *Splice) Length() int {
	return s.len
}

func (s *Splice) CountSlices() int {
	return len(s.slices)
}

func (s *Splice) Get(index int) (byte, error) {
	position, err := s.getPosition(index)
	if err != nil {
		return 0, err
	}
	return (*s.slices[position.slice])[position.offset], nil
}

func (s *Splice) GetUnsafe(index int) byte {
	position, err := s.getPosition(index)
	if err != nil {
		panic(1)
	}
	return (*s.slices[position.slice])[position.offset]
}

func (s *Splice) Iterate() *Iterator {
	return &Iterator{
		splice:    s,
		nextIndex: 0,
	}
}

type Iterator struct {
	splice    *Splice
	nextIndex int
}

func (s *Iterator) Next() bool {
	return s.nextIndex < s.splice.Length()
}

func (s *Iterator) Get() (int, byte, error) {
	resultIndex := s.nextIndex
	result, err := s.splice.Get(s.nextIndex)
	if err != nil {
		return 0, 0, err
	}
	s.nextIndex = s.nextIndex + 1
	return resultIndex, result, nil
}

func (s *Iterator) GetUnsafeWithIndex() (int, byte) {
	resultIndex := s.nextIndex
	result := s.splice.GetUnsafe(s.nextIndex)
	s.nextIndex = s.nextIndex + 1
	return resultIndex, result
}

func (s *Iterator) GetUnsafe() byte {
	result := s.splice.GetUnsafe(s.nextIndex)
	s.nextIndex = s.nextIndex + 1
	return result
}
