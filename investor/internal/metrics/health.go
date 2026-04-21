package metrics

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

// WritePID writes the current process PID to a file.
// If logger is provided, logs the operation result.
func WritePID(path string, log *logger.Logger) error {
	pid := os.Getpid()
	err := os.WriteFile(path, []byte(fmt.Sprintf("%d", pid)), 0644)
	
	if log != nil {
		ctx := context.Background()
		if err != nil {
			log.Warn(ctx, "failed to write PID file", logger.Fields{
				"component": "metrics",
				"error":     err.Error(),
				"pid_file":  path,
			})
		} else {
			log.Info(ctx, "PID file written", logger.Fields{
				"component": "metrics",
				"pid":       pid,
				"pid_file":  path,
			})
		}
	}
	
	return err
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
