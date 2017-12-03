package burl

import "strconv"
import "bytes"
import "fmt"

//Stat struct holds the value of a modifiable statistic, always an int.
//Enforces a max and min value, other nice things. Eventually will support
//temporary modifiers and the like (I think????).
type Stat struct {
	val int
	max int
	min int
}

//make a new stat with value at max, and min at 0.
func NewStat(v int) Stat {
	return Stat{v, v, 0}
}

func (s Stat) Get() int {
	return s.val
}

//Manually set a value. If v > s.Max, sets to max. if v < s.min, sets to min.
func (s *Stat) Set(v int) {
	s.val = Clamp(v, s.min, s.max)
}

//Modifies the stat value. Takes a delta (which can of course be negative).
//Calling this with d = 0 effectively re-ensures min <= val <= max
func (s *Stat) Mod(d int) {
	s.Set(s.val + d)
}

func (s Stat) Max() int {
	return s.max
}

func (s Stat) Min() int {
	return s.min
}

//Sets a new minimum. If this would make min > max, does nothing.
func (s *Stat) SetMin(m int) {
	if m <= s.max {
		s.min = m
		s.val = Clamp(s.val, s.min, s.max)
	}
}

//Sets a new maximum. If this would make max < min, does nothing.
func (s *Stat) SetMax(m int) {
	if m >= s.min {
		s.max = m
		s.val = Clamp(s.val, s.min, s.max)
	}
}

//Modifies the minimum. Takes a delta. Follows same rules as SetMin().
func (s *Stat) ModMin(d int) {
	s.SetMin(s.min + d)
}

//Modifies the maximum. Takes a delta. Follows same rules as SetMax().
func (s *Stat) ModMax(d int) {
	s.SetMax(s.max + d)
}

func (s Stat) IsMax() bool {
	return s.val == s.max
}

func (s Stat) IsMin() bool {
	return s.val == s.min
}

//returns a % (0-100) for the stat. If min == val == max, returns 0.
func (s Stat) GetPct() int {
	if s.min == s.max {
		return 0
	} else {
		return int(100 * (float32(s.val-s.min) / float32(s.max-s.min)))
	}
}

func (s Stat) String() string {
	return strconv.Itoa(s.val) + "/" + strconv.Itoa(s.max)
}

func (s Stat) GobEncode() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, s.min, s.max, s.val)
	return b.Bytes(), nil
}

func (s *Stat) GobDecode(data []byte) (err error) {
	b := bytes.NewBuffer(data)
	_, err = fmt.Fscanln(b, &s.min, &s.max, &s.val)
	return err
}
