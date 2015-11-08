# noti
Display a notification in OS X after a program finishes running.

## Installation
If you have Go installed, then you can do this.

```
go get github.com/variadico/noti
```

Otherwise, you can download the standalone binary on the
[releases page](https://github.com/variadico/noti/releases/latest). Then give
it execute permissions.

```
chmod u+x noti
```

## Usage
Just put `noti` at the beginning or end of your regular commands.

```
noti [options] [utility [args...]]
```

### Options
```
-t, -title
Set the notification title. If no arguments passed, default is "noti",
otherwise default is utility name.

-m, -message
Set notification message. Default is "Done!"

-s, -sound
Set notification sound. Default is Ping. Possible options are Basso,
Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
Submarine, Tink. Check /System/Library/Sounds for available sounds.

-f, -foreground
Bring the terminal to the foreground.

-p, -pushbullet
Send a Pushbullet notification. Access token must be set in NOTI_PB
environment variable.

-v, -version
Print noti version and exit.

-h, -help
Display usage information and exit.
```

## Examples
Display a notification when `tar` finishes compressing files.

```
noti tar -cjf music.tar.bz2 Music/
```

Display a notification when `apt-get` finishes updating on a remote server.

```
noti ssh you@server.com apt-get update
```

Set the notification title to "homebrew" and message to "Up to date" and
display it after `brew` finishes updating.

```
noti -t "homebrew" -m "up to date" brew upgrade --all
```

You can also add `noti` after a command, in case you forgot at the beginning.

```
clang foo.c -Wall -lm -L/usr/X11R6/lib -lX11 -o bizz; noti
```

Create a reminder to get back to a friend.

```
noti -t "Reply to Pedro" gsleep 5m &
```

Send a Pushbullet notification after tests finish.

```
noti -p go test ./...
```
