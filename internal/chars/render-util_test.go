package chars_test

import (
	"bufio"
	"errors"
	"strings"
	"testing"

	"github.com/Foxcapades/Argonaut/internal/chars"
	"github.com/Foxcapades/Argonaut/internal/util"
)

func TestPad(t *testing.T) {
	sb := new(strings.Builder)
	wri := bufio.NewWriter(sb)

	chars.Pad(10, wri)

	wri.Flush()

	if sb.String() != "          " {
		t.Fail()
	}
}

func TestIsBreakChar(t *testing.T) {
	if !chars.IsBreakChar(' ') {
		t.Error("expected true but got false")
	}
	if !chars.IsBreakChar('\t') {
		t.Error("expected true but got false")
	}
	if chars.IsBreakChar('a') {
		t.Error("expected false but got true")
	}
}

func TestDescriptionFormatter_Format01(t *testing.T) {
	expect := `    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
    tempor incididunt ut labore et dolore magna aliqua. Adipiscing enim eu
    turpis egestas pretium aenean. Morbi tincidunt ornare massa eget egestas
    purus viverra accumsan. Maecenas accumsan lacus vel facilisis volutpat est.`
	input := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do" +
		" eiusmod tempor incididunt ut labore et dolore magna aliqua. Adipiscing" +
		" enim eu turpis egestas pretium aenean. Morbi tincidunt ornare massa" +
		" eget egestas purus viverra accumsan. Maecenas accumsan lacus vel" +
		" facilisis volutpat est."

	sb := new(strings.Builder)
	sb.Grow(1024)

	writer := bufio.NewWriter(sb)

	formatter := chars.NewDescriptionFormatter("    ", 80, writer)
	util.Must(formatter.Format(input))
	util.Must(writer.Flush())

	if sb.String() != expect {
		t.Errorf("expected `%s`\ngot `%s`", expect, sb.String())
	}
}

func TestDescriptionFormatter_Format02(t *testing.T) {
	expect := `    Loremipsumdolorsitametconsecteturadipiscingelitseddoeiusmodtemporincididunt-
    utlaboreetdoloremagnaaliquaAdipiscingenimeuturpisegestaspretiumaeneanMorbit-
    inciduntornaremassaegetegestaspurusviverraaccumsanMaecenasaccumsanlacusvelf-
    acilisisvolutpatest.`
	input := "LoremipsumdolorsitametconsecteturadipiscingelitseddoeiusmodtemporincididuntutlaboreetdoloremagnaaliquaAdipiscingenimeuturpisegestaspretiumaeneanMorbitinciduntornaremassaegetegestaspurusviverraaccumsanMaecenasaccumsanlacusvelfacilisisvolutpatest."

	sb := new(strings.Builder)
	sb.Grow(1024)

	writer := bufio.NewWriter(sb)

	formatter := chars.NewDescriptionFormatter("    ", 80, writer)
	util.Must(formatter.Format(input))
	util.Must(writer.Flush())

	if sb.String() != expect {
		t.Errorf("expected `%s`\ngot `%s`", expect, sb.String())
	}
}

func TestDescriptionFormatter_Format03(t *testing.T) {
	defer func() { recover() }()

	chars.NewDescriptionFormatter("    ", 4, nil)

	t.Error("expected function to panic but it didn't")
}

func TestDescriptionFormatter_Format04(t *testing.T) {
	expect := `  A
  p
  p
  l
  e`

	input := "Apple"

	sb := new(strings.Builder)
	sb.Grow(64)

	writer := bufio.NewWriter(sb)
	formatter := chars.NewDescriptionFormatter("  ", 3, writer)

	util.Must(formatter.Format(input))
	util.Must(writer.Flush())

	if expect != sb.String() {
		t.Error("failed formatting expectation")
	}
}

func TestDescriptionFormatter_Format05(t *testing.T) {
	expect := `  A
  p
  p
  l
  e`

	input := "A\np\np\nl\ne"

	sb := new(strings.Builder)
	sb.Grow(64)

	writer := bufio.NewWriter(sb)
	formatter := chars.NewDescriptionFormatter("  ", 3, writer)

	util.Must(formatter.Format(input))
	util.Must(writer.Flush())

	if expect != sb.String() {
		t.Log(expect)
		t.Log(sb.String())
		t.Error("failed formatting expectation")
	}
}

func TestDescriptionFormatter_Format06(t *testing.T) {
	expect := `    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
    tempor incididunt ut labore et dolore magna aliqua. Adipiscing enim eu
    turpis egestas pretium aenean. Morbi tincidunt ornare massa eget egestas
    purus viverra accumsan. Maecenas accumsan lacus vel facilisis volutpat est.`
	input := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do" +
		" eiusmod tempor incididunt ut labore et dolore magna aliqua. Adipiscing" +
		" enim eu turpis egestas pretium aenean. Morbi tincidunt ornare massa" +
		" eget egestas purus viverra accumsan. Maecenas accumsan lacus vel" +
		" facilisis volutpat est."

	for i := 0; i < len(expect); i++ {
		buffer := FailingWriter{FailAfter: i}
		writer := bufio.NewWriterSize(&buffer, 1)

		formatter := chars.NewDescriptionFormatter("    ", 80, writer)
		if err := formatter.Format(input); err == nil {
			if err = writer.Flush(); err == nil {
				t.Error("expected error not to be nil but it was")
			}
		}
	}
}

func TestDescriptionFormatter_Format07(t *testing.T) {
	expect := `  A
  p
  p
  l
  e`

	input := "Apple"

	for i := 0; i < len(expect); i++ {
		buffer := FailingWriter{FailAfter: i}
		writer := bufio.NewWriterSize(&buffer, 1)

		formatter := chars.NewDescriptionFormatter("  ", 3, writer)
		if err := formatter.Format(input); err == nil {
			if err = writer.Flush(); err == nil {
				t.Error("expected error not to be nil but it was")
			}
		}
	}
}

func TestDescriptionFormatter_Format08(t *testing.T) {
	expect := `  A
  p
  p
  l
  e`

	input := "A\np\np\nl\ne"

	for i := 0; i < len(expect); i++ {
		buffer := FailingWriter{FailAfter: i}
		writer := bufio.NewWriterSize(&buffer, 1)

		formatter := chars.NewDescriptionFormatter("  ", 3, writer)
		if err := formatter.Format(input); err == nil {
			if err = writer.Flush(); err == nil {
				t.Error("expected error not to be nil but it was")
			}
		}
	}
}

type FailingWriter struct {
	FailAfter int
	current   int
}

func (f *FailingWriter) Write(p []byte) (n int, err error) {
	if f.current < f.FailAfter {
		f.current++
		return len(p), nil
	} else {
		return 0, errors.New("fake error")
	}
}
