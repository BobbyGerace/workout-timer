package fifo

import (
	"bufio"
	"os"
	"strings"
	"syscall"
)

// Listen creates the named pipe at path if it doesn't exist, then starts a
// goroutine that reads commands line-by-line and passes each to send.
// The goroutine runs for the lifetime of the process.
//
// Opening with O_RDWR keeps a write end permanently held by this process,
// so the open never blocks and external writers never block while we're running.
// EOF is never signalled (we hold the write end), so one open suffices for
// the lifetime of the process.
func Listen(path string, send func(string)) {
	syscall.Mkfifo(path, 0644) // no-op if already exists

	go func() {
		f, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)
		if err != nil {
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" {
				send(line)
			}
		}
	}()
}
