package metrics

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func WritePID(path string) error {
	return os.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

func ReadPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func IsProcessRunning(pidFile string) bool {
	pid, err := ReadPID(pidFile)
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return process != nil && process.Pid > 0
}

func WriteStats(path string, quotesPerSec float64, lastQuote time.Time) error {
	data := fmt.Sprintf("quotes_per_sec=%f\nlast_quote_time=%s", quotesPerSec, lastQuote.Format(time.RFC3339))
	return os.WriteFile(path, []byte(data), 0644)
}
