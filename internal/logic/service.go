package logic

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

type Store interface {
	Get(ctx context.Context, code string) (string, bool, error)
	Exists(ctx context.Context, code string) (bool, error)
	Save(ctx context.Context, code, url string) error
}

type Service struct{ Store Store }

var (
	ErrInvalidURL = errors.New("invalid url")
	ErrCodeInUse  = errors.New("code already in use")
	ErrNotFound   = errors.New("not found")
)

func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", ErrInvalidURL
	}
	return u.String(), nil
}

func (s *Service) Shorten(ctx context.Context, inURL, custom string) (string, error) {
	norm, err := normalizeURL(inURL)
	if err != nil {
		return "", err
	}
	// カスタムコード
	if custom != "" {
		ok, err := s.Store.Exists(ctx, custom)
		if err != nil {
			return "", err
		}
		if ok {
			return "", ErrCodeInUse
		}
		return custom, s.Store.Save(ctx, custom, norm)
	}
	// 自動生成（衝突したら桁数UP）
	for n := 6; n <= 9; n++ {
		code, err := GenerateCode(randomCode, n)
		if err != nil {
			return "", err
		}
		exists, err := s.Store.Exists(ctx, code)
		if err != nil {
			return "", err
		}
		if exists {
			continue
		}
		return code, s.Store.Save(ctx, code, norm)
	}
	return "", errors.New("failed to generate unique code")
}

func (s *Service) Resolve(ctx context.Context, code string) (string, bool, error) {
	u, cached, err := s.Store.Get(ctx, code)
	if err != nil {
		return "", false, err
	}
	if u == "" {
		return "", false, ErrNotFound
	}
	return u, cached, nil
}
