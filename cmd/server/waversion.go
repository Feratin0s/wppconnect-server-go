package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"go.mau.fi/whatsmeow/store"
)

var clientRevisionRe = regexp.MustCompile(`"client_revision":(\d+),`)

func updateWhatsAppVersion(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, 12*time.Second)
	defer cancel()

	v, err := fetchWAVersion(ctx)
	if err != nil {
		log.Printf("WARNING: could not update WhatsApp Web version: %v; keeping compiled version", err)
		return
	}
	store.SetWAVersion(v)
	log.Printf("WhatsApp Web client version updated: %s", v.String())
}

func fetchWAVersion(ctx context.Context) (store.WAVersionContainer, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://web.whatsapp.com", nil)
	if err != nil {
		return store.WAVersionContainer{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return store.WAVersionContainer{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return store.WAVersionContainer{}, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return store.WAVersionContainer{}, err
	}
	m := clientRevisionRe.FindSubmatch(data)
	if m == nil {
		return store.WAVersionContainer{}, fmt.Errorf("client_revision not found")
	}
	rev, err := strconv.ParseUint(string(m[1]), 10, 32)
	if err != nil {
		return store.WAVersionContainer{}, err
	}
	return store.WAVersionContainer{2, 3000, uint32(rev)}, nil
}
