# RINGIO

[![Build Status](https://drone.io/github.com/dullgiulio/ringio/status.png)](https://drone.io/github.com/dullgiulio/ringio/latest)

Ringio (pronounced *ring-yo*) is a tool for creating interactive data pipes. It is what you get when you mix "screen"/"tmux" with "tee".

### Usage

  - Start a new ringio session
```bash
$ ringio web-logs open &
```
  - Add some input agents
```bash
$ ringio web-logs input tail -f /var/log/httpd/access_log
```
  - Add some output agents or get output on the terminal
```bash
$ ringio web-logs output ./count-useragents
$ ringio web-logs output ./count-pagehits
$ ringio web-logs output # Will print to the console.
```
  - List agents for a session:
```bash
$ ringio web-logs list
```
  - See the internal log for the session:
```bash
$ ringio web-logs log
```
  - Close the session.
```bash
$ ringio web-logs close
```

### Installation

You need the Go lang development environment installed and set up:

```bash
$ go build github.org/dullgiulio/ringio
```

Downloads will be available soon.

