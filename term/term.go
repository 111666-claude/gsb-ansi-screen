// Package term 解析 ANSI 转义并维护屏幕缓冲。
package term

import "strconv"

// Screen 是一块固定大小的屏幕。
type Screen struct {
	w    int
	h    int
	flat []rune
	row  int
	col  int
}

// Counter 记录本轮扫描过的屏幕单元数。
type Counter struct {
	Scanned int
}

// New 建一块 w×h 的屏幕。
func New(w, h int) *Screen {
	s := &Screen{w: w, h: h, flat: make([]rune, w*h)}
	for i := range s.flat {
		s.flat[i] = ' '
	}
	return s
}

func (s *Screen) put(ch rune, c *Counter) {
	target := s.row*s.w + s.col
	idx := -1
	for i := range s.flat {
		c.Scanned++
		if i == target {
			idx = i
			break
		}
	}
	if idx >= 0 {
		s.flat[idx] = ch
	}
	s.col++
}

func params(raw string) []int {
	out := []int{}
	cur := ""
	for _, r := range raw {
		if r == ';' {
			n, _ := strconv.Atoi(cur)
			out = append(out, n)
			cur = ""
			continue
		}
		cur += string(r)
	}
	n, _ := strconv.Atoi(cur)
	return append(out, n)
}

// Feed 吃进一段带转义序列的文本。
func (s *Screen) Feed(data string, c *Counter) {
	i := 0
	for i < len(data) {
		if data[i] == 0x1b && i+1 < len(data) && data[i+1] == '[' {
			j := i + 2
			for j < len(data) && (data[j] == ';' || (data[j] >= '0' && data[j] <= '9')) {
				j++
			}
			_ = params(data[i+2 : j])
			i = j + 1
			continue
		}
		s.put(rune(data[i]), c)
		i++
	}
}

// Text 返回屏幕当前内容。
func (s *Screen) Text() []string {
	out := make([]string, 0, s.h)
	for r := 0; r < s.h; r++ {
		out = append(out, string(s.flat[r*s.w:(r+1)*s.w]))
	}
	return out
}
