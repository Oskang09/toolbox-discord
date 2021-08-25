package music

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	Instance *youtube.Video
	Title    string
	Display  string
	Duration time.Duration
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
}

func (player *Player) Start() {
	player.play()
}

func (player *Player) Stop() {
	player.stop = true
}

func (player *Player) Skip() {
	player.skip = true
}

func (player *Player) AddSong(link string) {
	video, _ := player.client.GetVideo(link)
	player.push(video)
	player.play()
}

func (player *Player) AddPlayList(link string) {
	playlist, err := player.client.GetPlaylist(link)
	log.Println(err)
	for _, vid := range playlist.Videos {
		video, _ := player.client.VideoFromPlaylistEntry(vid)
		player.push(video)
	}
	player.play()
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

		log.Println(player.queue)
		if len(player.queue) == 0 {
			player.playing = false
			player.channel.Disconnect()
			player.channel = nil
			return
		}

		// TODO: shuffle handle
		m := player.pop()
		player.current = &m
	}

	go sendPCM(player.channel, player.pcmChannel)
	video := player.current.Instance
	stream, _, _ := player.client.GetStream(video, video.Formats.WithAudioChannels().FindByQuality("tiny"))
	defer stream.Close()

	streamBytes, _ := ioutil.ReadAll(stream)

	run := exec.Command("ffmpeg.exe", "-i", "-", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	run.Dir = wd + "/cli"

	run.Stdin = bytes.NewReader(streamBytes)
	stdout, err := run.StdoutPipe()
	if err != nil {
		return
	}

	io.Copy(run.Stdout, stream)

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
		if err != nil {
			log.Println(len(audiobuf))
			fmt.Println("binary.Read: ", err)
			break
		}

		if player.skip {
			player.skip = false
			player.current = nil
			player.playing = false
			player.play()
			return
		}

		if player.stop {
			player.stop = false
			player.playing = false
			return
		}

		player.pcmChannel <- audiobuf
	}

	player.playing = false
	// TODO: repeat handle
	player.current = nil
	player.play()
}

func (player *Player) pop() (m Music) {
	m, player.queue = player.queue[0], player.queue[1:]
	return m
}

func (player *Player) push(video *youtube.Video) {
	display := ""
	if len(video.Thumbnails) > 0 {
		display = video.Thumbnails[0].URL
	}

	player.queue = append(player.queue, Music{
		Instance: video,
		Display:  display,
		Title:    video.Title,
		Duration: video.Duration,
	})
}

func NewPlayer(discord *discordgo.Session, guildID string, channelID string) *Player {
	player := new(Player)
	player.guildID = guildID
	player.channelID = channelID
	player.client = youtube.Client{}
	player.discord = discord
	player.channel = nil
	player.current = nil
	player.queue = make([]Music, 0)
	player.pcmChannel = make(chan []int16, 2)
	player.skip = false
	player.stop = false
	player.playing = false
	return player
}
