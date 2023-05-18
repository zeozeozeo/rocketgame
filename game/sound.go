package game

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const SAMPLE_RATE = 44100

var audioContext = audio.NewContext(SAMPLE_RATE)

type Sound struct {
	audioPlayer *audio.Player
	total       time.Duration
	seCh        chan []byte
	success     bool
}

func NewSound(data []byte) *Sound {
	sound := &Sound{}

	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4
	var s audioStream
	{
		var err error
		s, err = wav.DecodeWithoutResampling(bytes.NewReader(data))
		if err != nil {
			fmt.Println(err)
			return sound
		}
	}

	p, err := audioContext.NewPlayer(s)
	if err != nil {
		fmt.Println(err)
		return sound
	}
	sound.audioPlayer = p
	sound.total = time.Second * time.Duration(s.Length()) / bytesPerSample / SAMPLE_RATE
	if sound.total == 0 {
		sound.total = 1
	}
	sound.seCh = make(chan []byte)
	sound.success = true

	return sound
}

func (sound *Sound) SetVolume(vol float64) *Sound {
	if !sound.success {
		return sound
	}
	sound.audioPlayer.SetVolume(vol)
	return sound
}

func (sound *Sound) Play() *Sound {
	if !sound.success {
		return sound
	}
	sound.audioPlayer.Play()
	return sound
}
