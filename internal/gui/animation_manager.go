package gui

import "time"

type AnimationManager struct {
	animations []func(dt time.Duration) bool
}

func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		animations: []func(dt time.Duration) bool{},
	}
}

func (am *AnimationManager) Add(animation func(time.Duration) bool) {
	am.animations = append(am.animations, animation)
}

func (am *AnimationManager) Update(dt time.Duration) {
	newAnimation := am.animations[:0]

	for _, anim := range am.animations {
		if !anim(dt) {
			newAnimation = append(newAnimation, anim)
		}
	}

	am.animations = newAnimation
}
