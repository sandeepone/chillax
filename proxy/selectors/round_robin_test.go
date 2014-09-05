package selectors

import (
	"testing"
)

var (
	a = "1.2.3.4"
	b = "2.3.4.5"
)

func TestLen(t *testing.T) {
	r := NewRoundRobin([]string{a, b})
	if r.Len() != 2 {
		t.Error("Round Robin Len should be 2.")
	}
}

func TestChoose(t *testing.T) {
	r := NewRoundRobin([]string{a, b})

	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
	if r.Choose() != b {
		t.Errorf("r.Choose() should == %v", b)
	}
	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
	if r.Choose() != b {
		t.Errorf("r.Choose() should == %v", b)
	}
}

func TestChooseEmpty(t *testing.T) {
	r := NewRoundRobin([]string{})

	if r.Choose() != "" {
		t.Errorf("r.Choose() should == empty")
	}
}

func TestAdd(t *testing.T) {
	r := NewRoundRobin([]string{a})

	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}

	r.Add(b)

	if r.Choose() != b {
		t.Errorf("r.Choose() should == %v", b)
	}
	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
}

func TestAddEmpty(t *testing.T) {
	r := NewRoundRobin([]string{})
	r.Add(a)
	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
}

func TestRemove(t *testing.T) {
	r := NewRoundRobin([]string{a, b})

	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
	if r.Choose() != b {
		t.Errorf("r.Choose() should == %v", b)
	}

	r.Remove(b)

	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}
	if r.Choose() != a {
		t.Errorf("r.Choose() should == %v", a)
	}

	r.Remove(a)

	if r.Len() != 0 {
		t.Error("Round Robin Len should be 0.")
	}
}

func TestRemoveEmpty(t *testing.T) {
	r := NewRoundRobin([]string{})
	r.Remove(a)

	if r.Len() != 0 {
		t.Error("Round Robin Len should be 0.")
	}

	r.Add(a)

	if r.Len() != 1 {
		t.Error("Round Robin Len should be 1.")
	}

	r.Remove(a)

	if r.Len() != 0 {
		t.Error("Round Robin Len should be 0.")
	}

	r.Remove(a)

	if r.Len() != 0 {
		t.Error("Round Robin Len should be 0.")
	}
}
