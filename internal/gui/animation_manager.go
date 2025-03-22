// AnimationManager will control running animations in gui main loop. It will run animation functions that takes deltatime and returns true when aniamtion completed. It then remove completed animation from slice of animations so it will not run again.
package gui

import "time"

type AnimationManager struct {
	animations []func(dt time.Duration) bool
}

// NewAnimationManager returns a new AnimationManager pointer with an empty slice as animation argument
func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		animations: []func(dt time.Duration) bool{},
	}
}

// Add append a new animation function to animation slice
func (am *AnimationManager) Add(animation func(time.Duration) bool) {
	am.animations = append(am.animations, animation)
}

// Update runs all animation functions and removes finished animations from animation slice
func (am *AnimationManager) Update(dt time.Duration) {
	// TODO: Check if using slices.Clip will improve performance and memory usage here or not
	newAnimation := am.animations[:0]

	for _, anim := range am.animations {
		if !anim(dt) {
			newAnimation = append(newAnimation, anim)
		}
	}

	am.animations = newAnimation
}
