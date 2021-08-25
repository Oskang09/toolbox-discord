package music

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

// Borrowed from https://github.com/bwmarrin/dgvoice

var (
	sendpcm bool
	mu      sync.Mutex
)

const (
	channels  int = 2
	frameRate int = 48000
	frameSize int = 960
	maxBytes  int = (frameSize * 2) * 2
)

func sendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	mu.Lock()
	if sendpcm || pcm == nil {
		mu.Unlock()
		return
	}
	sendpcm = true
	mu.Unlock()
	defer func() { sendpcm = false }()

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		fmt.Println("gopus.NewEncoder: ", err)
		return
	}

	for {
		recv, ok := <-pcm
		if !ok {
			return
		}

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			fmt.Println("opusEncoder.Encode: ", err)
			return
		}

		if !v.Ready || v.OpusSend == nil {
			return
		}
		v.OpusSend <- opus
	}
}
