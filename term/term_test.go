package term

import "testing"

func TestPlainTextFillsRow(t *testing.T) {
	s := New(3, 1)
	s.Feed("abc", &Counter{})
	if got := s.Text()[0]; got != "abc" {
		t.Fatalf("普通字符应顺序填入，得到 %q", got)
	}
}

func TestScreenIsSpaceFilled(t *testing.T) {
	s := New(3, 1)
	if got := s.Text()[0]; got != "   " {
		t.Fatalf("初始应是空格，得到 %q", got)
	}
}

func TestRowsAndCols(t *testing.T) {
	s := New(2, 3)
	if len(s.Text()) != 3 {
		t.Fatalf("应有 3 行，得到 %d", len(s.Text()))
	}
	if len(s.Text()[0]) != 2 {
		t.Fatalf("每行应有 2 列，得到 %d", len(s.Text()[0]))
	}
}

func TestColorSequenceProducesNoText(t *testing.T) {
	s := New(3, 1)
	s.Feed("\x1b[31m", &Counter{})
	if got := s.Text()[0]; got != "   " {
		t.Fatalf("颜色序列不该写字符，得到 %q", got)
	}
}

func TestScannedIsCounted(t *testing.T) {
	s := New(2, 2)
	c := &Counter{}
	s.Feed("a", c)
	if c.Scanned == 0 {
		t.Fatal("应该统计扫描过的屏幕单元")
	}
}
