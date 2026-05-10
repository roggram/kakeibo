package srv

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"

	"srv.exe.dev/db"
)

type Server struct {
	DB        *sql.DB
	Hostname  string
	StaticDir string
}

func New(dbPath, hostname string) (*Server, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(thisFile)
	s := &Server{
		Hostname:  hostname,
		StaticDir: filepath.Join(baseDir, "static"),
	}
	if err := s.setUpDatabase(dbPath); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) setUpDatabase(dbPath string) error {
	wdb, err := db.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	s.DB = wdb
	if err := db.RunMigrations(wdb); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

func (s *Server) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/state", s.handleGetState)
	mux.HandleFunc("PUT /api/state", s.handlePutState)
	mux.Handle("/", http.FileServer(http.Dir(s.StaticDir)))
	slog.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	var data string
	err := s.DB.QueryRowContext(r.Context(), `SELECT data FROM app_state WHERE id = 1`).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		data = "null"
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	io.WriteString(w, data)
}

func (s *Server) handlePutState(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 16*1024*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	_, err = s.DB.ExecContext(r.Context(),
		`INSERT INTO app_state (id, data, updated_at) VALUES (1, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(id) DO UPDATE SET data = excluded.data, updated_at = CURRENT_TIMESTAMP`,
		string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
