package ins

import (
	"errors"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	var i = New()
	val, ok := i.Get("name")
	if ok {
		t.Errorf("val want %v, got %v; ok want %v, got %v", nil, val, false, ok)
	}
	if err := i.Set("name", "alex"); err != nil {
		t.Error(err)
		return
	}
	val, ok = i.Get("name")
	if !ok || val != "alex" {
		t.Errorf("val want %v, got %v; ok want %v, got %v", "alex", val, true, ok)
	}
}

func TestSetEx(t *testing.T) {
	var i = New()
	if err := i.SetEx("fruit", "watermelon", 1); err != nil {
		t.Error(err)
		return
	}
	val, ok := i.Get("fruit")
	if !ok || val != "watermelon" {

		t.Errorf("val want %v, got %v; ok want %v, got %v", "watermelon", val, true, ok)
		return
	}
	time.Sleep(1 * time.Second)
	val, ok = i.Get("fruit")
	if ok {
		t.Errorf("val want %v, got %v; ok want %v, got %v", nil, val, false, ok)
	}
}

func TestSetNx(t *testing.T) {
	var i = New()
	if err := i.SetNx("fruit", "watermelon"); err != nil {
		t.Error(err)
		return
	}

	for num := 0; num < 100; num++ {
		go func() {
			if err := i.SetNx("fruit", "watermelon"); err == nil {
				t.Error(errors.New("want err, got nil"))
				return
			}
		}()
	}

	val, ok := i.Get("fruit")
	if !ok || val != "watermelon" {
		t.Errorf("val want %v, got %v; ok want %v, got %v", "watermelon", val, true, ok)
		return
	}
}

func TestSetNxEx(t *testing.T) {
	var i = New()
	if err := i.SetNxEx("fruit", "watermelon", 2); err != nil {
		t.Error(err)
		return
	}

	val, ok := i.Get("fruit")
	if !ok || val != "watermelon" {
		t.Errorf("val want %v, got %v; ok want %v, got %v", "watermelon", val, true, ok)
		return
	}

	for num := 0; num < 100; num++ {
		go func() {
			if err := i.SetNxEx("fruit", "watermelon", 2); err == nil {
				t.Error(errors.New("want err, got nil"))
				return
			}
		}()
	}

	time.Sleep(2 * time.Second)

	val, ok = i.Get("fruit")
	if ok {
		t.Errorf("val want %v, got %v; ok want %v, got %v", nil, val, false, ok)
		return
	}

	if err := i.SetNxEx("fruit", "watermelon", 2); err != nil {
		t.Errorf("want nil, got %v", err)
		return
	}

	val, ok = i.Get("fruit")
	if !ok || val != "watermelon" {
		t.Errorf("val want %v, got %v; ok want %v, got %v", "watermelon", val, true, ok)
		return
	}
}

func TestDel(t *testing.T) {
	var i = New()
	if err := i.Set("sports", "badminton"); err != nil {
		t.Error(err)
		return
	}
	if val, ok := i.Get("sports"); !ok || val != "badminton" {
		t.Errorf("val want %v, got %v; ok want %v, got %v", "badminton", val, true, ok)
		return
	}
	if err := i.Del("sports"); err != nil {
		t.Error(err)
		return
	}
	if val, ok := i.Get("sports"); ok {
		t.Errorf("val want %v, got %v; ok want %v, got %v", nil, val, false, ok)
		return
	}
}

func BenchmarkSet(b *testing.B) {
	var i = New()
	for num := 0; num < b.N; num++ {
		str := time.Now().String()
		if err := i.Set("timeStr", str); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkGet(b *testing.B) {
	var i = New()
	str := time.Now().String()
	if err := i.Set("timeStr", str); err != nil {
		b.Error(err)
		return
	}
	for num := 0; num < b.N; num++ {
		if _, ok := i.Get("timeStr"); !ok {
			b.Error(errors.New("want ok, got false"))
		}
	}
}
