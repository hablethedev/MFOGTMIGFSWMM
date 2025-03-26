package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	startTime  time.Time
	running    bool
	timerValue time.Duration
	fontFace   font.Face
	socketPath string = "/tmp/timer.sock"
	textColor  color.Color = color.RGBA{255, 140, 185, 255}
)

type Game struct{}

func (g *Game) Update() error {
	if running {
		timerValue = time.Since(startTime)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Alpha{0})
	msg := fmt.Sprintf("%02d:%02d.%03d", int(timerValue.Minutes()), int(timerValue.Seconds())%60, timerValue.Milliseconds()%1000)
	text.Draw(screen, msg, fontFace, 40, 80, textColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 300, 120
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Println("received command:", command)
		switch command {
		case "start":
			if !running {
				startTime = time.Now()
				running = true
			}
		case "stop":
			running = false
		case "reset":
			timerValue = 0
			running = false
		}
	}
}

func listenUnixSocket() {
	os.Remove(socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("error creating unix socket:", err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func init() {
	fontData, err := ioutil.ReadFile("/usr/share/fonts/open-sans/OpenSans-Bold.ttf")
	if err != nil {
		panic(err)
	}
	
	tt, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}

	fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	go listenUnixSocket()

	ebiten.SetWindowSize(300, 120)
	ebiten.SetWindowTitle("i am NOT writing that all out as a timer you cant even see")
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
