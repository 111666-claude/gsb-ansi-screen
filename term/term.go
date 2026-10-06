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
	s.flat[s.row*s.w+s.col] = ch
	c.Scanned++
	s.col++
	if s.col >= s.w {
		s.col = s.w - 1
	}
}

func (s *Screen) moveTo(row, col int) {
	s.row = min(max(row-1, 0), s.h-1)
	s.col = min(max(col-1, 0), s.w-1)
}

func (s *Screen) clear() {
	for i := range s.flat {
		s.flat[i] = ' '
	}
	s.row = 0
	s.col = 0
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
			if j >= len(data) {
				break
			}
			p := params(data[i+2 : j])
			switch data[j] {
			case 'H', 'f':
				row, col := 1, 1
				if len(p) > 0 {
					row = p[0]
				}
				if len(p) > 1 {
					col = p[1]
				}
				s.moveTo(row, col)
			case 'J':
				if len(p) > 0 && p[0] == 2 {
					s.clear()
				}
			}
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
