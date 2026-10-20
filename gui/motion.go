package main

import (
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// Responses of the springs, in the sense Apple uses: about how long the
// value takes to reach its target. They are not durations; see springStep.
const (
	responseMove  = 350 * time.Millisecond // pages and notifications
	responseMeter = 450 * time.Millisecond // the fill of a meter
)

// springStep returns how far a critically damped spring (damping ratio 1,
// no bounce) has gone from 0 to 1 after time t, for a given response:
//
//	x(t) = 1 - (1 + ωt)·e^(-ωt), with ω = 2π / response
//
// It starts at 0 with no speed, moves fast and settles softly at 1 without
// passing it.
func springStep(t, response time.Duration) float64 {
	omega := 2 * math.Pi / response.Seconds()
	wt := omega * t.Seconds()
	return 1 - (1+wt)*math.Exp(-wt)
}

// animateSpring moves a value from `from` to `to` with a critically damped
// spring, calls tick with the value on every frame and done (if not nil)
// when it arrives. Fyne runs the ticks on the UI goroutine, so tick may
// update widgets. The caller keeps the returned animation to Stop it when a
// new target arrives: the next spring then starts from the value on screen,
// so the motion never jumps.
func animateSpring(from, to float32, response time.Duration, tick func(v float32), done func()) *fyne.Animation {
	finish := func() {
		tick(to) // exactly to: from + (to-from)*1 may differ in the last bit
		if done != nil {
			done()
		}
	}
	// After 1.5 responses the spring is at more than 99.9%, so it stops
	// there and snaps the last tiny part.
	total := response * 3 / 2
	anim := fyne.NewAnimation(total, func(progress float32) {
		if progress >= 1 {
			finish()
			return
		}
		x := float32(springStep(time.Duration(float64(total)*float64(progress)), response))
		tick(from + (to-from)*x)
	})
	anim.Curve = fyne.AnimationLinear // progress is plain time; the spring gives the shape
	anim.Start()
	return anim
}

// stopAnimation stops anim if it is running. A nil anim is fine.
func stopAnimation(anim *fyne.Animation) {
	if anim != nil {
		anim.Stop()
	}
}

// slideLayout places its objects moved down by offset pixels. Animating
// offset slides the content in or out without changing its size.
type slideLayout struct {
	offset float32
}

// Compile-time check: the build fails if slideLayout stops being a layout.
var _ fyne.Layout = (*slideLayout)(nil)

// Layout gives every object the whole size, moved down by offset.
func (s *slideLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
		o.Move(fyne.NewPos(0, s.offset))
	}
}

// MinSize is the biggest minimum size of the objects, as a Stack does.
func (s *slideLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return layout.NewStackLayout().MinSize(objects)
}

// slider is a container whose content can slide vertically.
type slider struct {
	layout *slideLayout
	box    *fyne.Container
	anim   *fyne.Animation
}

// newSlider returns a slider showing content.
func newSlider(content ...fyne.CanvasObject) *slider {
	s := &slider{layout: &slideLayout{}}
	s.box = container.New(s.layout, content...)
	return s
}

// slideTo moves the content to offset with a spring, starting from where it
// is now on screen, and runs done (if not nil) when it arrives.
func (s *slider) slideTo(offset float32, done func()) {
	stopAnimation(s.anim)
	s.anim = animateSpring(s.layout.offset, offset, responseMove, func(v float32) {
		s.layout.offset = v
		s.box.Refresh() // a container Refresh runs its layout again
	}, done)
}

// jump moves the content to offset at once, without animation.
func (s *slider) jump(offset float32) {
	stopAnimation(s.anim)
	s.layout.offset = offset
	s.box.Refresh()
}

// toaster shows short notifications at the bottom of the window, like the
// banners of macOS: they rise from the bottom edge and leave the same way.
type toaster struct {
	text   *canvas.Text
	slide  *slider
	layer  fyne.CanvasObject // what the shell puts over the window content
	shown  int               // number of the last notification
	hidden float32           // offset that leaves the banner under the edge
}

// toastTime is how long a notification stays on screen.
const toastTime = 2600 * time.Millisecond

// newToaster returns a toaster, hidden under the bottom edge.
func newToaster() *toaster {
	t := &toaster{text: canvas.NewText("", ink), hidden: 120}
	t.text.TextStyle = fyne.TextStyle{Bold: true}
	banner := newPanel(container.NewHBox(newAppIcon(28), container.NewCenter(t.text)))
	t.slide = newSlider(container.NewBorder(nil, container.NewCenter(banner), nil, nil))
	t.slide.jump(t.hidden)
	t.layer = container.New(layout.NewCustomPaddedLayout(0, 24, 0, 0), t.slide.box)
	t.layer.Hide() // it is shown only while a notification is on screen
	return t
}

// show rises a notification with msg and hides it after toastTime. A new
// notification that arrives while another one is on screen replaces it and
// starts from where the banner is, without jumping.
func (t *toaster) show(msg string) {
	t.shown++
	n := t.shown
	t.text.Text = msg
	t.text.Refresh()
	t.layer.Show()
	t.slide.slideTo(0, nil)

	// time.AfterFunc runs its function on another goroutine, and widgets
	// may only be changed on the UI goroutine, so it asks Fyne to do it
	// there with fyne.Do.
	time.AfterFunc(toastTime, func() {
		fyne.Do(func() {
			if n != t.shown {
				return // a newer notification is on screen; it hides itself
			}
			t.slide.slideTo(t.hidden, func() {
				if n == t.shown {
					t.layer.Hide()
				}
			})
		})
	})
}
