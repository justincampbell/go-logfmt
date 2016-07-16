package logfmt_test

import (
	"testing"
	"time"

	"github.com/justincampbell/go-logfmt"
)

type Person struct {
	Name                    string        `logfmt:"name"`
	Age                     int           `logfmt:"age"`
	Alive                   bool          `logfmt:"is_alive"`
	BirthDate               time.Time     `logfmt:"birth_date,Monday, 02-Jan-06 15:04:05 MST"`
	FavoriteColors          []string      `logfmt:"fav_colors"`
	FavoriteNumbers         []int         `logfmt:"fav_numbers"`
	CanHoldBreathFor        time.Duration `logfmt:"hold_breath"`
	CanHoldBreathForSeconds time.Duration `logfmt:"hold_breath_seconds,s"`
	Empty                   string        `logfmt:"empty"`
}

var person = &Person{
	Name:                    "Jane Doe",
	Age:                     32,
	Alive:                   true,
	BirthDate:               time.Date(2016, time.July, 15, 21, 4, 45, 0, time.FixedZone("MST", -8*3600)),
	FavoriteColors:          []string{"orange", "sky blue"},
	FavoriteNumbers:         []int{1, 2},
	CanHoldBreathFor:        45 * time.Second,
	CanHoldBreathForSeconds: 45 * time.Second,
	Empty: "",
}

func Test_structtags_Encode(t *testing.T) {
	b, err := logfmt.Encode(*person)
	if err != nil {
		t.Fatal(err)
	}
	expected := `name="Jane Doe" age=32 is_alive=true birth_date="Friday, 15-Jul-16 21:04:45 MST" fav_colors="orange,sky blue" fav_numbers=1,2 hold_breath=45s hold_breath_seconds=45`
	if string(b) != expected {
		t.Fatalf("bad: %#s", b)
	}
}

func Test_structtags_Unmarshal_empty(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(``), p)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_structtags_Unmarshal_string(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`name="Jane Doe"`), p)
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "Jane Doe" {
		t.Fatalf("bad: %#v", p.Name)
	}
}

func Test_structtags_Unmarshal_int(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`age=32`), p)
	if err != nil {
		t.Fatal(err)
	}

	if p.Age != 32 {
		t.Fatalf("bad: %#v", p.Age)
	}
}

func Test_structtags_Unmarshal_bool(t *testing.T) {
	p := &Person{}

	// true
	if err := logfmt.Unmarshal([]byte(`is_alive=true`), p); err != nil {
		t.Fatal(err)
	}
	if !p.Alive {
		t.Fatalf("bad: %#v", p.Alive)
	}

	// false
	if err := logfmt.Unmarshal([]byte(`is_alive=false`), p); err != nil {
		t.Fatal(err)
	}
	if p.Alive {
		t.Fatalf("bad: %#v", p.Alive)
	}

	// 1
	if err := logfmt.Unmarshal([]byte(`is_alive=1`), p); err != nil {
		t.Fatal(err)
	}
	if !p.Alive {
		t.Fatalf("bad: %#v", p.Alive)
	}

	// 0
	if err := logfmt.Unmarshal([]byte(`is_alive=0`), p); err != nil {
		t.Fatal(err)
	}
	if p.Alive {
		t.Fatalf("bad: %#v", p.Alive)
	}

	// anything else
	if err := logfmt.Unmarshal([]byte(`is_alive=what`), p); err == nil {
		t.Fatalf("bad: expected error")
	}
}

func Test_structtags_Unmarshal_intSlice(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`fav_numbers=1,2`), p)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.FavoriteNumbers) != 2 {
		t.Fatalf("bad: %#v", p.FavoriteNumbers)
	}
	if p.FavoriteNumbers[0] != 1 {
		t.Fatalf("bad: %#v", p.FavoriteNumbers[0])
	}
	if p.FavoriteNumbers[1] != 2 {
		t.Fatalf("bad: %#v", p.FavoriteNumbers[1])
	}
}

func Test_structtags_Unmarshal_stringSlice(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`fav_colors="orange,sky blue"`), p)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.FavoriteColors) != 2 {
		t.Fatalf("bad: %#v", p.FavoriteColors)
	}
	if p.FavoriteColors[0] != "orange" {
		t.Fatalf("bad: %#v", p.FavoriteColors[0])
	}
	if p.FavoriteColors[1] != "sky blue" {
		t.Fatalf("bad: %#v", p.FavoriteColors[1])
	}
}

func Test_structtags_Unmarshal_duration(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`hold_breath=45s`), p)
	if err != nil {
		t.Fatal(err)
	}

	if p.CanHoldBreathFor != 45*time.Second {
		t.Fatalf("bad: %#v", p.CanHoldBreathFor)
	}
}

func Test_structtags_Unmarshal_durationWithSuffix(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`hold_breath_seconds=45`), p)
	if err != nil {
		t.Fatal(err)
	}

	if p.CanHoldBreathForSeconds != 45*time.Second {
		t.Fatalf("bad: %#v", p.CanHoldBreathForSeconds)
	}
}

func Test_structtags_Unmarshal_time(t *testing.T) {
	p := &Person{}
	err := logfmt.Unmarshal([]byte(`birth_date="Friday, 15-Jul-16 09:04:45 MST"`), p)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := time.Parse(
		"Monday, 02-Jan-06 15:04:05 MST",
		"Friday, 15-Jul-16 09:04:45 MST",
	)

	if p.BirthDate != expected {
		t.Fatalf("bad: %#v", p.BirthDate)
	}
}
