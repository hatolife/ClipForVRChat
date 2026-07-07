package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

const (
	singleInstanceLockFile  = "ClipForVRChat.lock"
	singleInstanceStateFile = "single-instance.json"
)

var (
	errSingleInstanceAlreadyRunning = errors.New("ClipForVRChat is already running")
	singleInstanceDirFunc           = defaultSingleInstanceDir
)

type singleInstance struct {
	dir       string
	statePath string
	lock      *flock.Flock
	server    *singleInstanceServer
}

type singleInstanceState struct {
	PID            int    `json:"pid"`
	ExecutablePath string `json:"executablePath"`
	Version        string `json:"version"`
	Revision       string `json:"revision"`
	Endpoint       string `json:"endpoint"`
	Token          string `json:"token"`
	StartedAt      string `json:"startedAt"`
}

type singleInstanceChoice int

const (
	singleInstanceChoiceCancel singleInstanceChoice = iota
	singleInstanceChoiceActivate
	singleInstanceChoiceReplace
)

func initializeSingleInstance(stderr io.Writer, args []string) (*singleInstance, error) {
	if err := rejectElevatedProcess(); err != nil {
		return nil, err
	}

	instance, err := tryStartPrimarySingleInstance()
	if err == nil {
		return instance, nil
	}
	if !errors.Is(err, errSingleInstanceAlreadyRunning) {
		return nil, err
	}

	instance, handled, err := handleExistingSingleInstance(stderr, args)
	if err != nil {
		showNativeMessage("ClipForVRChat", err.Error())
		return nil, err
	}
	if handled {
		return nil, nil
	}
	return instance, nil
}

func tryStartPrimarySingleInstance() (*singleInstance, error) {
	dir, err := singleInstanceDirFunc()
	if err != nil {
		return nil, err
	}
	lock, err := acquireInstanceLock()
	if err != nil {
		return nil, err
	}

	server, err := startSingleInstanceServer()
	if err != nil {
		_ = lock.Unlock()
		return nil, fmt.Errorf("単一起動IPCを開始できませんでした: %w", err)
	}

	state := singleInstanceState{
		PID:            os.Getpid(),
		ExecutablePath: currentExecutablePath(),
		Version:        appVersion(),
		Revision:       revision,
		Endpoint:       server.Endpoint(),
		Token:          server.Token(),
		StartedAt:      time.Now().Format(time.RFC3339Nano),
	}
	statePath := filepath.Join(dir, singleInstanceStateFile)
	if err := writeSingleInstanceState(statePath, state); err != nil {
		server.Close()
		_ = lock.Unlock()
		return nil, err
	}

	return &singleInstance{
		dir:       dir,
		statePath: statePath,
		lock:      lock,
		server:    server,
	}, nil
}

func handleExistingSingleInstance(stderr io.Writer, args []string) (*singleInstance, bool, error) {
	state, err := readExistingSingleInstanceState()
	if err != nil {
		return nil, false, err
	}
	if err := sendSingleInstanceCommand(state, "ping", 2*time.Second); err != nil {
		return nil, false, fmt.Errorf("既存のClipForVRChatが応答しません。既存プロセスを手動で閉じてから起動してください: %w", err)
	}

	choice := chooseExistingSingleInstanceAction(state, stderr)
	switch choice {
	case singleInstanceChoiceActivate:
		command := "activate"
		if len(args) > 0 {
			command = "open-paths"
		}
		if err := sendSingleInstanceCommandWithPaths(state, command, args, 3*time.Second); err != nil {
			return nil, false, fmt.Errorf("既存のClipForVRChatをアクティブ化できませんでした: %w", err)
		}
		return nil, true, nil
	case singleInstanceChoiceReplace:
		if err := sendSingleInstanceCommand(state, "shutdown", 3*time.Second); err != nil {
			return nil, false, fmt.Errorf("既存のClipForVRChatへ終了要求を送信できませんでした: %w", err)
		}
		instance, err := waitAndStartPrimarySingleInstance(15 * time.Second)
		if err != nil {
			return nil, false, err
		}
		return instance, false, nil
	default:
		return nil, true, nil
	}
}

func waitAndStartPrimarySingleInstance(timeout time.Duration) (*singleInstance, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		instance, err := tryStartPrimarySingleInstance()
		if err == nil {
			return instance, nil
		}
		lastErr = err
		if !errors.Is(err, errSingleInstanceAlreadyRunning) {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil, fmt.Errorf("既存のClipForVRChatの終了を待ちましたが、単一起動ロックを取得できませんでした: %w", lastErr)
}

func acquireInstanceLock() (*flock.Flock, error) {
	lock, locked, err := tryAcquireInstanceLock()
	if err != nil {
		return nil, fmt.Errorf("起動ロックを取得できませんでした: %w", err)
	}
	if !locked {
		return nil, errSingleInstanceAlreadyRunning
	}
	return lock, nil
}

func tryAcquireInstanceLock() (*flock.Flock, bool, error) {
	dir, err := singleInstanceDirFunc()
	if err != nil {
		return nil, false, err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, false, fmt.Errorf("単一起動ディレクトリを作成できませんでした: %w", err)
	}
	lockPath := filepath.Join(dir, singleInstanceLockFile)
	fileLock := flock.New(lockPath)
	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, false, err
	}
	if !locked {
		return nil, false, nil
	}
	return fileLock, true, nil
}

func defaultSingleInstanceDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("ユーザーキャッシュディレクトリを取得できませんでした: %w", err)
	}
	return filepath.Join(base, "ClipForVRChat", "single-instance"), nil
}

func currentExecutablePath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return exe
}

func readExistingSingleInstanceState() (singleInstanceState, error) {
	dir, err := singleInstanceDirFunc()
	if err != nil {
		return singleInstanceState{}, err
	}
	return readSingleInstanceState(filepath.Join(dir, singleInstanceStateFile))
}

func readSingleInstanceState(path string) (singleInstanceState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return singleInstanceState{}, fmt.Errorf("既存プロセス情報を読み込めませんでした: %w", err)
	}
	var state singleInstanceState
	if err := json.Unmarshal(data, &state); err != nil {
		return singleInstanceState{}, fmt.Errorf("既存プロセス情報が壊れています: %w", err)
	}
	if strings.TrimSpace(state.Endpoint) == "" || strings.TrimSpace(state.Token) == "" {
		return singleInstanceState{}, errors.New("既存プロセス情報にIPC接続情報がありません")
	}
	return state, nil
}

func writeSingleInstanceState(path string, state singleInstanceState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("単一起動ディレクトリを作成できませんでした: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("既存プロセス情報を書き込めませんでした: %w", err)
	}
	return nil
}

func (s *singleInstance) BindApp(app *App) {
	if s == nil || s.server == nil || app == nil {
		return
	}
	s.server.SetHandlers(app.activateFromSingleInstance, app.shutdownFromSingleInstance, app.openPathsFromSingleInstance)
}

func (s *singleInstance) Close() {
	if s == nil {
		return
	}
	if s.server != nil {
		s.server.Close()
	}
	if s.statePath != "" {
		_ = os.Remove(s.statePath)
	}
	if s.lock != nil {
		_ = s.lock.Unlock()
	}
}

type singleInstanceServer struct {
	listener net.Listener
	token    string
	done     chan struct{}
	once     sync.Once

	mu        sync.RWMutex
	activate  func() error
	shutdown  func() error
	openPaths func([]string) error
}

type singleInstanceRequest struct {
	Token   string   `json:"token"`
	Command string   `json:"command"`
	Paths   []string `json:"paths,omitempty"`
}

type singleInstanceResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

func startSingleInstanceServer() (*singleInstanceServer, error) {
	token, err := randomSingleInstanceToken()
	if err != nil {
		return nil, err
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	server := &singleInstanceServer{
		listener: listener,
		token:    token,
		done:     make(chan struct{}),
		activate: func() error {
			return errors.New("既存ウィンドウはまだ起動準備中です")
		},
		shutdown: func() error {
			go func() {
				time.Sleep(100 * time.Millisecond)
				os.Exit(0)
			}()
			return nil
		},
		openPaths: func([]string) error {
			return errors.New("既存ウィンドウはまだ起動準備中です")
		},
	}
	go server.serve()
	return server, nil
}

func randomSingleInstanceToken() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("単一起動IPC tokenを生成できませんでした: %w", err)
	}
	return hex.EncodeToString(buf[:]), nil
}

func (s *singleInstanceServer) Endpoint() string {
	if s == nil || s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

func (s *singleInstanceServer) Token() string {
	if s == nil {
		return ""
	}
	return s.token
}

func (s *singleInstanceServer) SetHandlers(activate func() error, shutdown func() error, openPaths func([]string) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if activate != nil {
		s.activate = activate
	}
	if shutdown != nil {
		s.shutdown = shutdown
	}
	if openPaths != nil {
		s.openPaths = openPaths
	}
}

func (s *singleInstanceServer) Close() {
	if s == nil {
		return
	}
	s.once.Do(func() {
		_ = s.listener.Close()
		close(s.done)
	})
}

func (s *singleInstanceServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				continue
			}
		}
		go s.handle(conn)
	}
}

func (s *singleInstanceServer) handle(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))

	var req singleInstanceRequest
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&req); err != nil {
		writeSingleInstanceResponse(conn, false, "invalid request")
		return
	}
	if req.Token != s.token {
		writeSingleInstanceResponse(conn, false, "invalid token")
		return
	}

	switch strings.ToLower(strings.TrimSpace(req.Command)) {
	case "ping":
		writeSingleInstanceResponse(conn, true, "")
	case "activate":
		if err := s.callActivate(); err != nil {
			writeSingleInstanceResponse(conn, false, err.Error())
			return
		}
		writeSingleInstanceResponse(conn, true, "")
	case "open-paths":
		if err := s.callOpenPaths(req.Paths); err != nil {
			writeSingleInstanceResponse(conn, false, err.Error())
			return
		}
		writeSingleInstanceResponse(conn, true, "")
	case "shutdown":
		if err := s.callShutdown(); err != nil {
			writeSingleInstanceResponse(conn, false, err.Error())
			return
		}
		writeSingleInstanceResponse(conn, true, "")
	default:
		writeSingleInstanceResponse(conn, false, "unknown command")
	}
}

func (s *singleInstanceServer) callActivate() error {
	s.mu.RLock()
	fn := s.activate
	s.mu.RUnlock()
	return fn()
}

func (s *singleInstanceServer) callShutdown() error {
	s.mu.RLock()
	fn := s.shutdown
	s.mu.RUnlock()
	return fn()
}

func (s *singleInstanceServer) callOpenPaths(paths []string) error {
	s.mu.RLock()
	fn := s.openPaths
	s.mu.RUnlock()
	return fn(paths)
}

func writeSingleInstanceResponse(w io.Writer, ok bool, message string) {
	_ = json.NewEncoder(w).Encode(singleInstanceResponse{OK: ok, Message: message})
}

func sendSingleInstanceCommand(state singleInstanceState, command string, timeout time.Duration) error {
	return sendSingleInstanceCommandWithPaths(state, command, nil, timeout)
}

func sendSingleInstanceCommandWithPaths(state singleInstanceState, command string, paths []string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", state.Endpoint, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	if err := json.NewEncoder(conn).Encode(singleInstanceRequest{Token: state.Token, Command: command, Paths: paths}); err != nil {
		return err
	}
	var res singleInstanceResponse
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&res); err != nil {
		return err
	}
	if !res.OK {
		if res.Message == "" {
			res.Message = "既存プロセスが要求を拒否しました"
		}
		return errors.New(res.Message)
	}
	return nil
}
