# MFOGTMIGFSWMM
Mainly Fedora only Gaming Timer made in Go for Sway WM (Mainly)

![screenshot](https://github.com/hablethedev/MFOGTMIGFSWMM/blob/main/screen.png?raw=true)

# Dependencies For Anything
socat - `sudo dnf install socat`

swaywm - you can use others that support binding commands to buttons

# How (Development)

- clone the repo `git clone github.com/habletedev/MFOGTMIGFSWMM.git`
- cd into it
- go run main.go

# Bindings

`bindsym Mod4+z exec --no-startup-id sh -c "echo start | socat - UNIX-CONNECT:/tmp/timer.sock"`

`bindsym Mod4+x exec --no-startup-id sh -c "echo stop | socat - UNIX-CONNECT:/tmp/timer.sock"`

`bindsym Mod4+c exec --no-startup-id sh -c "echo reset | socat - UNIX-CONNECT:/tmp/timer.sock"`
