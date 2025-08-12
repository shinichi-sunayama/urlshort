package logic

import (
	"context"
	"errors"
	"testing"
)

// メモリStore（モック）
type memStore struct {
	m   map[string]string
	err error
}

func (m *memStore) Get(ctx context.Context, code string) (string, bool, error) {
	if m.err != nil {
		return "", false, m.err
	}
	v := m.m[code]
	if v == "" {
		return "", false, nil
	}
	return v, false, nil
}
func (m *memStore) Exists(ctx context.Context, code string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	_, ok := m.m[code]
	return ok, nil
}
func (m *memStore) Save(ctx context.Context, code, url string) error {
	if m.err != nil {
		return m.err
	}
	if m.m == nil {
		m.m = map[string]string{}
	}
	m.m[code] = url
	return nil
}

func TestShorten_ValidURL(t *testing.T) {
	svc := &Service{Store: &memStore{}}
	code, err := svc.Shorten(context.Background(), "https://example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) == 0 {
		t.Fatalf("expected code, got empty")
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	svc := &Service{Store: &memStore{}}
	_, err := svc.Shorten(context.Background(), "not-a-url", "")
	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestResolve_NotFound(t *testing.T) {
	svc := &Service{Store: &memStore{m: map[string]string{}}}
	_, _, err := svc.Resolve(context.Background(), "nope")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestResolve_Found(t *testing.T) {
	svc := &Service{Store: &memStore{m: map[string]string{"abc": "https://example.com"}}}
	u, cached, err := svc.Resolve(context.Background(), "abc")
	if err != nil || u == "" || cached {
		t.Fatalf("unexpected: u=%q cached=%v err=%v", u, cached, err)
	}
}
