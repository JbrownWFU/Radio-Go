package player

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type StationInfo struct {
	Name        string
	Genre       string
	Description string
	Bitrate     string
	URL         string
}

type Metadata struct {
	StreamTitle string
	Artist      string
	Title       string
}

// ParseMetadata is a heuristic parser for ICECast metadata.
func ParseMetadata(rawTitle string) Metadata {
	meta := Metadata{StreamTitle: rawTitle}

	// Common delimiters in order of likelihood
	delimiters := []string{" - ", " | ", " : ", " // "}

	for _, d := range delimiters {
		if strings.Contains(rawTitle, d) {
			parts := strings.SplitN(rawTitle, d, 2)
			meta.Artist = strings.TrimSpace(parts[0])
			meta.Title = strings.TrimSpace(parts[1])

			// Further cleaning
			meta.Artist = cleanMetadata(meta.Artist)
			meta.Title = cleanMetadata(meta.Title)
			return meta
		}
	}

	// Fallback if no delimiter found
	meta.Artist = "Unknown"
	meta.Title = cleanMetadata(rawTitle)
	return meta
}

func cleanMetadata(s string) string {
	noise := []string{"(LIVE)", "[HQ]", "(Official Audio)", "(Official Video)", " - Single"}
	for _, n := range noise {
		s = strings.ReplaceAll(s, n, "")
		s = strings.ReplaceAll(s, strings.ToLower(n), "")
	}
	return strings.TrimSpace(s)
}

type Player struct {
	cmd        *exec.Cmd
	mu         sync.Mutex
	active     bool
	streamBody io.ReadCloser
	MetaChan   chan Metadata
}

func New() *Player {
	return &Player{
		MetaChan: make(chan Metadata, 10),
	}
}

// Start starts ffplay with the given URL and calls onMeta when metadata is received.
func (p *Player) Start(url string) (*StationInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.active {
		p.stopLocked()
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Icy-MetaData", "1")
	req.Header.Set("User-Agent", "VLC/3.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	info := &StationInfo{
		Name:        resp.Header.Get("icy-name"),
		Genre:       resp.Header.Get("icy-genre"),
		Description: resp.Header.Get("icy-description"),
		Bitrate:     resp.Header.Get("icy-br"),
		URL:         resp.Header.Get("icy-url"),
	}

	metaInt := 0
	metaHeader := resp.Header.Get("icy-metaint")
	if metaHeader != "" {
		metaInt, _ = strconv.Atoi(metaHeader)
	}

	p.streamBody = resp.Body

	args := []string{
		"-nodisp", "-loglevel", "quiet", "-hide_banner", "-autoexit",
		"-i", "pipe:0",
	}

	p.cmd = exec.Command("ffplay", args...)
	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("stdin pipe error: %w", err)
	}

	if err := p.cmd.Start(); err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("ffplay start error: %w", err)
	}

	p.active = true
	currCmd := p.cmd

	go func() {
		defer stdin.Close()
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)

		if metaInt == 0 {
			io.Copy(stdin, reader)
			return
		}

		buf := make([]byte, metaInt)
		for {
			if _, err := io.ReadFull(reader, buf); err != nil {
				break
			}
			if _, err := stdin.Write(buf); err != nil {
				break
			}

			byteSize, err := reader.ReadByte()
			if err != nil {
				break
			}

			metaLen := int(byteSize) * 16
			if metaLen > 0 {
				metaBytes := make([]byte, metaLen)
				if _, err := io.ReadFull(reader, metaBytes); err != nil {
					break
				}

				metaStr := strings.Trim(string(metaBytes), "\x00")
				if metaStr != "" {
					data := parseMetaMap(metaStr)
					if title, ok := data["StreamTitle"]; ok {
						select {
						case p.MetaChan <- ParseMetadata(title):
						default:
						}
					}
				}
			}
		}
	}()

	go func() {
		currCmd.Wait()
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.cmd == currCmd {
			p.active = false
			if p.streamBody != nil {
				p.streamBody.Close()
			}
		}
	}()

	return info, nil
}

func (p *Player) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopLocked()
}

func (p *Player) stopLocked() error {
	if p.streamBody != nil {
		p.streamBody.Close()
	}
	if !p.active || p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	p.cmd.Process.Kill()
	p.active = false
	return nil
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active
}

func parseMetaMap(meta string) map[string]string {
	data := make(map[string]string)
	parts := strings.Split(meta, ";")
	for _, part := range parts {
		if strings.Contains(part, "=") {
			keyValue := strings.SplitN(part, "=", 2)
			key := keyValue[0]
			val := strings.Trim(keyValue[1], "'")
			data[key] = val
		}
	}
	return data
}
