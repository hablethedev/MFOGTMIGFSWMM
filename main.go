package main

import (
	"bufio"
	"fmt"
	"image/color"
	"net"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

var (
	startTime  time.Time
	running    bool
	timerValue time.Duration
	fontFace   font.Face = inconsolata.Bold8x16
	socketPath string    = "/tmp/timer.sock"
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
	text.Draw(screen, msg, fontFace, 20, 40, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 200, 60
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

func main() {
	go listenUnixSocket()

	ebiten.SetWindowSize(200, 60)
	ebiten.SetWindowTitle("i am NOT writing that all out as a timer you cant even see")
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
