package audio

import (
	"bytes"
	_ "embed"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	oto "github.com/ebitengine/oto/v3"
)

//go:embed beep.wav
var beepWav []byte

var (
	otoCtx  *oto.Context
	otoOnce sync.Once
)

func initOto() {
	ctx, readyChan, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return
	}
	<-readyChan
	otoCtx = ctx
}

// Beep plays a short alert sound. The command runs in a goroutine so it
// never blocks the tick loop.
func Beep() {
	go func() {
		if err := playBeep(); err != nil {
			fmt.Fprint(os.Stderr, "\a")
		}
	}()
}

func playBeep() error {
	otoOnce.Do(initOto)
	if otoCtx == nil {
		return fmt.Errorf("oto context unavailable")
	}
	pcm := generateSine(880, 0.5, 44100)
	player := otoCtx.NewPlayer(bytes.NewReader(pcm))
	player.Play()
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	player.Close()
	return nil
}

// generateSine returns raw 16-bit signed little-endian PCM for a sine wave
// at the given frequency (Hz) and duration (seconds).
func generateSine(freq, duration float64, sampleRate int) []byte {
	numSamples := int(float64(sampleRate) * duration)
	buf := make([]byte, numSamples*2)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := math.Sin(2 * math.Pi * freq * t)
		val := int16(sample * 16383) // half amplitude
		buf[i*2] = byte(val)
		buf[i*2+1] = byte(val >> 8)
	}
	return buf
}
