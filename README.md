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
Added agent %1
```
  - Add some output agents or get output on the terminal
```bash
$ ringio web-logs output ./count-useragents
Added agent %2
$ ringio web-logs output wc -l
Added agent %3
$ ringio web-logs output # Will print to the console.
```
  - List agents for a session:
```bash
$ ringio web-logs list
1 R <- tail -f /var/log/httpd/access_log
2 R -> ./count-useragents
3 R -> wc -l
```
  - See the internal log for the session:
```bash
$ ringio web-logs log
```
  - Stop the agent counting the lines:
```bash
$ ringio web-logs stop %3
```
  - Retrieve 'wc -l' output (you can see any output by filtering it explicitly):
```bash
$ ringio web-logs output %3
4526
```
  - Close the session.
```bash
$ ringio web-logs close
```

### Filtering

 - Input and output agents can be filtered by writing their ID:
```bash
$ ringio my-session output 3
```
   Will display the output of agent %3. To filter out, negate the ID of the agent (ex: -3).

### Installation

You need the Go lang development environment installed and set up:

```bash
$ go build github.org/dullgiulio/ringio
```

Please see [the releases page](https://github.com/dullgiulio/ringio/releases) for further information.
