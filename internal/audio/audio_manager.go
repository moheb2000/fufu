package audio

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

type AudioManager struct {
	sampleRate beep.SampleRate
	music      *Music
	sounds     map[string]*beep.Buffer
}

type Music struct {
	id       int
	streamer beep.StreamSeekCloser
	ctrl     *beep.Ctrl
}

func NewAudioManager(fps int) *AudioManager {
	aum := &AudioManager{
		sampleRate: beep.SampleRate(44100),
		sounds:     make(map[string]*beep.Buffer),
	}

	speaker.Init(aum.sampleRate, aum.sampleRate.N(time.Second/time.Duration(fps)))

	return aum
}

// PlayMusic plays the music specified in the path parameter
func (aum *AudioManager) PlayMusic(path string, loop bool) error {
	m := Music{}
	if aum.music != nil {
		m.id = aum.music.id + 1
	}

	// Stop the old running music (if it exists) before starting a new one
	// TODO: Check if the VN engine needs the ability to play two musics simulantly or not
	aum.StopMusic()

	streamer, format, err := decodeAudioFile(path)
	if err != nil {
		return err
	}
	m.streamer = streamer

	aum.music = &m

	// Loop over music based on the value of loop bool parameter
	var resampled *beep.Resampler
	if loop {
		lp, err := beep.Loop2(aum.music.streamer)
		if err != nil {
			return err
		}

		// Fix the sample rate to a consistant sample rate
		resampled = beep.Resample(4, format.SampleRate, aum.sampleRate, lp)
	} else {
		resampled = beep.Resample(4, format.SampleRate, aum.sampleRate, aum.music.streamer)
	}

	ctrl := &beep.Ctrl{
		Streamer: resampled,
		Paused:   false,
	}

	m.ctrl = ctrl

	// Close the music after playing completed
	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		// When old music is closed manually the callback called here. If immediately another music starts playing the aum music will not be nil and Close again.
		if aum.music != nil && m.id == aum.music.id {
			aum.music.streamer.Close()
			aum.music = nil
		}
	})))

	return nil
}

// Stop the music
func (aum *AudioManager) StopMusic() {
	if aum.music != nil {
		speaker.Lock()
		aum.music.streamer.Close()
		speaker.Unlock()

		aum.music = nil
	}
}

func (aum *AudioManager) PauseMusic() {
	if aum.music != nil {
		aum.music.ctrl.Paused = true
	}
}

func (aum *AudioManager) ResumeMusic() {
	if aum.music != nil {
		aum.music.ctrl.Paused = false
	}
}

func (aum *AudioManager) PlaySound(path string) error {
	// Check if the sound is already loaded on memory or not and if not, load it first
	if _, exists := aum.sounds[path]; !exists {
		streamer, format, err := decodeAudioFile(path)
		if err != nil {
			return err
		}

		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		aum.sounds[path] = buffer
		streamer.Close()
	}

	sound := aum.sounds[path].Streamer(0, aum.sounds[path].Len())
	speaker.Play(sound)

	return nil
}

// decodeAudioFile gets a path file to music, decode it based on the format of the file and returns like other standard beep Decode functions
func decodeAudioFile(path string) (beep.StreamSeekCloser, beep.Format, error) {
	// Initialize variables
	var stream beep.StreamSeekCloser
	format := beep.Format{}

	f, err := os.Open(path)
	if err != nil {
		return stream, format, err
	}

	// Get the format of the file
	l := strings.Split(path, ".")
	suffix := ""
	if len(l) > 0 {
		suffix = strings.ToLower(l[len(l)-1])
	}

	// Decode file based on the format
	switch suffix {
	case "mp3":
		stream, format, err = mp3.Decode(f)
	case "wav":
		stream, format, err = wav.Decode(f)
	case "ogg":
		stream, format, err = vorbis.Decode(f)
	case "flac":
		stream, format, err = flac.Decode(f)
	default:
		return stream, format, fmt.Errorf("%s format is not supported by audio package", suffix)
	}

	// Return the final result. If the process be successful, the err is nil
	return stream, format, err
}
