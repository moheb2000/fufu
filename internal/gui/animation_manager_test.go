package gui

import (
	"testing"
	"time"
)

func TestNewAnimationManager(t *testing.T) {
	am := NewAnimationManager()

	if am.animations == nil {
		t.Errorf("animations slice in AnimationManager should be an empty slice; got: %v", nil)
	}

	if len(am.animations) != 0 {
		t.Errorf("animations field length expected %v; got: %v", 0, len(am.animations))
	}
}

func TestAdd(t *testing.T) {
	am := NewAnimationManager()

	am.Add(func(dt time.Duration) bool { return false })

	if len(am.animations) != 1 {
		t.Errorf("Add should increase animations field length by one; expected: %v; got: %v", 1, len(am.animations))
	}
}

func TestUpdate(t *testing.T) {
	am := NewAnimationManager()

	am.Add(func(dt time.Duration) bool {
		return dt == time.Second*2
	})

	for i := 0; i < 3; i++ {
		am.Update(time.Second * time.Duration(i))

		if i == 1 {
			if len(am.animations) != 1 {
				t.Errorf("function that returns false should remain in animation slice; animation length expected: %v; got: %v", 1, len(am.animations))
			}
		} else if i == 2 {
			if len(am.animations) != 0 {
				t.Errorf("function that returns true should be removed from animation slice; animation length expected: %v; got: %v", 0, len(am.animations))
			}
		}
	}
}
