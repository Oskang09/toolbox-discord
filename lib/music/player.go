package music

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

var wd, _ = os.Getwd()

// Music :
type Music struct {
	instance *youtube.Video
	stream   []byte

	Title       string
	Author      string
	Description string
	Display     struct {
		URL    string
		Width  int
		Height int
	}
}

// Player :
type Player struct {
	guildID    string
	channelID  string
	client     youtube.Client
	discord    *discordgo.Session
	channel    *discordgo.VoiceConnection
	current    *Music
	queue      []Music
	pcmChannel chan []int16
	playing    bool
	skip       bool
	stop       bool
	repeat     bool
}

func (player *Player) Start() {
	go player.play()
}

func (player *Player) Stop() {
	player.stop = true
}

func (player *Player) Skip() {
	player.skip = true
}

func (player *Player) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(player.queue), func(i, j int) {
		player.queue[i], player.queue[j] = player.queue[j], player.queue[i]
	})
}

func (player *Player) Repeat() bool {
	player.repeat = !player.repeat
	return player.repeat
}

func (player *Player) AddSong(link string) (*youtube.Video, error) {
	video, err := player.client.GetVideo(link)
	if err != nil {
		return nil, err
	}

	player.push(video)
	go player.play()
	return video, nil
}

func (player *Player) AddPlayList(link string) (int64, error) {
	playlist, err := player.client.GetPlaylist(link)
	if err != nil {
		return 0, err
	}

	for _, v := range playlist.Videos {
		go func(vid *youtube.PlaylistEntry) {
			video, err := player.client.VideoFromPlaylistEntry(vid)
			if err != nil {
				return
			}

			player.push(video)
		}(v)
	}
	go player.play()
	return int64(len(playlist.Videos)), nil
}

func (player *Player) ClearQueue() {
	player.queue = make([]Music, 0)
}

func (player *Player) Current() *Music {
	return player.current
}

func (player *Player) ListQueueAndCurrent() (*Music, []Music) {
	return player.current, player.queue
}

func (player *Player) play() {
	if player.playing {
		return
	}

	player.playing = true
	if player.channel == nil {
		player.channel, _ = player.discord.ChannelVoiceJoin(player.guildID, player.channelID, false, true)
	}

	if player.current == nil {
		if len(player.queue) == 0 {
			player.playing = false
			player.channel.Disconnect()
			player.channel = nil
			return
		}

		m := player.pop()
		player.current = &m
	}

	if player.current.stream == nil {
		video := player.current.instance
		stream, _, _ := player.client.GetStream(video, video.Formats.WithAudioChannels().FindByQuality("tiny"))
		defer stream.Close()

		player.current.stream, _ = ioutil.ReadAll(stream)
	}

	go sendPCM(player.channel, player.pcmChannel)
	run := exec.Command("ffmpeg.exe", "-i", "-", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	run.Dir = wd + "/cli"

	run.Stdin = bytes.NewReader(player.current.stream)
	stdout, err := run.StdoutPipe()
	if err != nil {
		return
	}

	err = run.Start()
	if err != nil {
		return
	}
	defer run.Process.Kill()

	player.channel.Speaking(true)
	defer func() {
		if player.channel != nil {
			player.channel.Speaking(false)
		}
	}()

	audiobuf := make([]int16, frameSize*channels)
	for {
		err = binary.Read(stdout, binary.LittleEndian, &audiobuf)
		if err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			fmt.Println("binary.Read: ", err)
			return
		}

		if player.skip {
			player.skip = false
			player.current = nil
			player.playing = false
			go player.play()
			return
		}

		if player.stop {
			player.stop = false
			player.playing = false
			player.channel.Disconnect()
			player.channel = nil
			return
		}

		player.pcmChannel <- audiobuf
	}

	player.playing = false
	if !player.repeat {
		player.current = nil
	}
	go player.play()
}

func (player *Player) pop() (m Music) {
	m, player.queue = player.queue[0], player.queue[1:]
	return m
}

func (player *Player) push(video *youtube.Video) {
	music := Music{
		instance:    video,
		stream:      nil,
		Title:       video.Title,
		Author:      video.Author,
		Description: video.Description,
	}

	if len(video.Thumbnails) > 0 {
		music.Display.URL = video.Thumbnails[0].URL
		music.Display.Width = int(video.Thumbnails[0].Width)
		music.Display.Height = int(video.Thumbnails[0].Height)
	}

	player.queue = append(player.queue, music)
}

func (player *Player) CanContinue(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	state, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		return false
	}

	player.channelID = state.ChannelID
	player.guildID = i.GuildID
	player.discord = s
	return true
}

func NewPlayer() *Player {
	player := new(Player)
	player.guildID = ""
	player.channelID = ""
	player.client = youtube.Client{}
	player.discord = nil
	player.channel = nil
	player.current = nil
	player.queue = make([]Music, 0)
	player.pcmChannel = make(chan []int16, 2)
	player.skip = false
	player.stop = false
	player.repeat = false
	player.playing = false
	return player
}
