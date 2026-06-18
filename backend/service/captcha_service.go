package service

import (
	"fmt"
	mrand "math/rand"
	"sync"
	"time"
)

type captchaEntry struct {
	answer    string
	expiresAt time.Time
}

type CaptchaService struct {
	store sync.Map
}

func NewCaptchaService() *CaptchaService {
	s := &CaptchaService{}
	go s.cleanup()
	return s
}

func (s *CaptchaService) cleanup() {
	for {
		time.Sleep(time.Minute)
		now := time.Now()
		s.store.Range(func(key, value interface{}) bool {
			if entry, ok := value.(captchaEntry); ok && now.After(entry.expiresAt) {
				s.store.Delete(key)
			}
			return true
		})
	}
}

type captchaQuestion struct {
	text   string
	answer func(a, b int) string
}

var questions = []captchaQuestion{
	{text: "%d + %d = ?", answer: func(a, b int) string { return fmt.Sprintf("%d", a+b) }},
	{text: "%d × %d = ?", answer: func(a, b int) string { return fmt.Sprintf("%d", a*b) }},
}

func randomColor() string {
	colors := []string{"#e74c3c", "#2ecc71", "#3498db", "#f39c12", "#9b59b6", "#1abc9c", "#e67e22"}
	return colors[mrand.Intn(len(colors))]
}

func (s *CaptchaService) Generate() (id, svg string) {
	q := questions[mrand.Intn(len(questions))]
	a := mrand.Intn(9) + 1
	b := mrand.Intn(9) + 1
	answer := q.answer(a, b)
	questionText := fmt.Sprintf(q.text, a, b)

	id = fmt.Sprintf("captcha_%d_%d", time.Now().UnixNano(), mrand.Intn(10000))

	s.store.Store(id, captchaEntry{
		answer:    answer,
		expiresAt: time.Now().Add(5 * time.Minute),
	})

	noiseLines := randomNoiseLines(4)
	noiseDots := randomNoiseDots(30)
	distort := distortFilter()

	textLength := len([]rune(questionText))
	charSpacing := 180 / (textLength + 1)
	textElements := ""
	for i, ch := range questionText {
		x := 30 + i*charSpacing
		y := 25 + mrand.Intn(15)
		rotate := mrand.Intn(30) - 15
		fontSize := 22 + mrand.Intn(8)
		textElements += fmt.Sprintf(`<text x="%d" y="%d" font-size="%d" font-family="Arial,Helvetica,sans-serif" fill="%s" transform="rotate(%d,%d,%d)" font-weight="bold" opacity="0.85">%s</text>`,
			x, y, fontSize, randomColor(), rotate, x, y, string(ch))
	}

	svg = fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="60" viewBox="0 0 200 60">
  <defs>%s</defs>
  <rect width="200" height="60" fill="#fafafa" rx="6" stroke="#e0e0e0" stroke-width="1"/>
  <rect x="3" y="3" width="194" height="54" fill="none" stroke="#f0f0f0" stroke-width="1" rx="4"/>
  %s
  %s
  <g filter="url(#d)">%s</g>
</svg>`, distort, noiseLines, noiseDots, textElements)
	return
}

func randomNoiseLines(count int) string {
	lines := ""
	for i := 0; i < count; i++ {
		x1 := mrand.Intn(200)
		y1 := mrand.Intn(60)
		x2 := mrand.Intn(200)
		y2 := mrand.Intn(60)
		opacity := 0.15 + mrand.Float64()*0.2
		lines += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="%.1f" opacity="%.2f"/>`,
			x1, y1, x2, y2, randomColor(), 0.5+mrand.Float64()*1.0, opacity)
	}
	return lines
}

func randomNoiseDots(count int) string {
	dots := ""
	for i := 0; i < count; i++ {
		x := mrand.Intn(200)
		y := mrand.Intn(60)
		r := 0.5 + mrand.Float64()*1.5
		dots += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%.1f" fill="%s" opacity="0.3"/>`,
			x, y, r, randomColor())
	}
	return dots
}

func distortFilter() string {
	scale := 3 + mrand.Intn(5)
	return fmt.Sprintf(`<filter id="d"><feTurbulence type="turbulence" baseFrequency="0.02 0.03" numOctaves="3" result="t"/><feDisplacementMap in="SourceGraphic" in2="t" scale="%d" xChannelSelector="R" yChannelSelector="G"/></filter>`, scale)
}

func (s *CaptchaService) Validate(captchaID, answer string) bool {
	if captchaID == "" || answer == "" {
		return false
	}
	val, ok := s.store.Load(captchaID)
	if !ok {
		return false
	}
	entry := val.(captchaEntry)
	if time.Now().After(entry.expiresAt) {
		s.store.Delete(captchaID)
		return false
	}
	s.store.Delete(captchaID)
	return entry.answer == answer
}
