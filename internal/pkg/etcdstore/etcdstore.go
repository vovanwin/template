package etcdstore

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Config настройки подключения к etcd
type Config struct {
	Endpoints   []string
	Prefix      string
	DialTimeout time.Duration
}

// EtcdStore реализует config.FlagStore с бэкендом etcd
type EtcdStore struct {
	client   *clientv3.Client
	prefix   string
	log      *slog.Logger
	defaults map[string]any

	mu     sync.RWMutex
	cache  map[string]string
	cancel context.CancelFunc
}

// New создает EtcdStore и загружает текущие значения из etcd
func New(cfg Config, log *slog.Logger) (*EtcdStore, error) {
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "/flags/"
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd connect: %w", err)
	}

	s := &EtcdStore{
		client: client,
		prefix: cfg.Prefix,
		log:    log,
		cache:  make(map[string]string),
	}

	// Загружаем текущие значения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	resp, err := client.Get(ctx, cfg.Prefix, clientv3.WithPrefix())
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("etcd initial load: %w", err)
	}

	for _, kv := range resp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), cfg.Prefix)
		s.cache[key] = string(kv.Value)
	}

	return s, nil
}

// SetDefaults записывает дефолтные значения в etcd для ключей, которых ещё нет
func (s *EtcdStore) SetDefaults(defaults map[string]any) {
	s.mu.Lock()
	s.defaults = defaults
	s.mu.Unlock()

	for key, val := range defaults {
		s.mu.RLock()
		_, exists := s.cache[key]
		s.mu.RUnlock()

		if exists {
			continue
		}

		strVal := fmt.Sprintf("%v", val)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, err := s.client.Put(ctx, s.prefix+key, strVal)
		cancel()

		if err != nil {
			s.log.Warn("etcd: failed to set default", slog.String("key", key), slog.Any("error", err))
			continue
		}

		s.mu.Lock()
		s.cache[key] = strVal
		s.mu.Unlock()
	}
}

func (s *EtcdStore) GetBool(key string, def bool) bool {
	s.mu.RLock()
	v, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func (s *EtcdStore) GetInt(key string, def int) int {
	s.mu.RLock()
	v, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func (s *EtcdStore) GetFloat(key string, def float64) float64 {
	s.mu.RLock()
	v, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func (s *EtcdStore) GetString(key string, def string) string {
	s.mu.RLock()
	v, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok {
		return def
	}
	return v
}

// Watch запускает etcd watcher на prefix и вызывает onChange при изменении
func (s *EtcdStore) Watch(ctx context.Context, onChange func(key string)) error {
	ctx, s.cancel = context.WithCancel(ctx)

	go func() {
		wch := s.client.Watch(ctx, s.prefix, clientv3.WithPrefix())
		for resp := range wch {
			for _, ev := range resp.Events {
				key := strings.TrimPrefix(string(ev.Kv.Key), s.prefix)

				s.mu.Lock()
				if ev.Type == clientv3.EventTypeDelete {
					delete(s.cache, key)
				} else {
					s.cache[key] = string(ev.Kv.Value)
				}
				s.mu.Unlock()

				if onChange != nil {
					onChange(key)
				}

				s.log.Debug("etcd flag updated",
					slog.String("key", key),
					slog.String("value", string(ev.Kv.Value)),
				)
			}
		}
	}()

	return nil
}

// Close закрывает etcd клиент и watcher
func (s *EtcdStore) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return s.client.Close()
}
