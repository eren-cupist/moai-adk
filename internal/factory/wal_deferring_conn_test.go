// wal_deferring_conn_test.go — a queue-database connection that keeps every
// committed write in backlog.db-wal (no autocheckpoint), shared by the WAL
// identity tests.
package factory

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

// openDeferringConn opens a queue database on a connection that will NOT
// autocheckpoint. The caller closes it; while it is open, every committed
// write it makes stays in backlog.db-wal.
func openDeferringConn(t *testing.T, dbPath string) *sql.DB {
	t.Helper()
	v := url.Values{}
	v.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", backlogBusyTimeoutMS))
	v.Add("_pragma", "journal_mode(WAL)")
	v.Add("_pragma", "wal_autocheckpoint(0)")
	v.Add("_txlock", "immediate")
	p := filepath.ToSlash(dbPath)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	u := url.URL{Scheme: "file", Path: p, RawQuery: v.Encode()}
	db, err := sql.Open(sqliteDriverName, u.String())
	if err != nil {
		t.Fatalf("open deferring connection: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}
