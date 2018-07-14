package gosql

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	fmtLogQuery     = `Query: %s`
	fmtLogArgs      = `Args:  %#v`
	fmtLogError     = `Error: %v`
	fmtLogTimeTaken = `Time:  %0.5fs`
)

var (
	reInvisibleChars = regexp.MustCompile(`[\s\r\n\t]+`)
)

// QueryStatus represents the status of a query after being executed.
type QueryStatus struct {
	Query string
	Args  interface{}

	Start time.Time
	End   time.Time

	Err error
}

// String returns a formatted log message.
func (q *QueryStatus) String() string {
	lines := make([]string, 0, 8)

	if query := q.Query; query != "" {
		query = reInvisibleChars.ReplaceAllString(query, ` `)
		query = strings.TrimSpace(query)
		lines = append(lines, fmt.Sprintf(fmtLogQuery, query))
	}

	if q.Args != nil {
		lines = append(lines, fmt.Sprintf(fmtLogArgs, q.Args))
	}

	if q.Err != nil {
		lines = append(lines, fmt.Sprintf(fmtLogError, q.Err))
	}

	lines = append(lines, fmt.Sprintf(fmtLogTimeTaken, float64(q.End.UnixNano()-q.Start.UnixNano())/float64(1e9)))

	return strings.Join(lines, "\n")
}

// Logger represents a logging collector. You can pass a logging collector to
// db.DefaultSettings.SetLogger(myCollector) to make it collect db.QueryStatus messages
// after executing a query.
type Logger interface {
	Log(*QueryStatus)
	SetLogging(logging bool)
}

type defaultLogger struct {
	logging bool
}

func (lg *defaultLogger) Log(m *QueryStatus) {
	if lg.logging {
		log.Printf("\n\t%s\n\n", strings.Replace(m.String(), "\n", "\n\t", -1))
	}
}

func (lg *defaultLogger) SetLogging(logging bool) {
	lg.logging = logging
}

var _ Logger = (*defaultLogger)(nil)