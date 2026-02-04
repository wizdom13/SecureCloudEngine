//go:build !plan9 && !js

// Package cachev2 implements a virtual provider to cache existing remotes.
package cachev2

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/rclone/rclone/backend/cache"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configmap"
)

const (
	// DefCacheChunkSize is the default value for chunk size
	DefCacheChunkSize = cache.DefCacheChunkSize
	// DefCacheTotalChunkSize is the default value for the maximum size of stored chunks
	DefCacheTotalChunkSize = cache.DefCacheTotalChunkSize
	// DefCacheChunkCleanInterval is the interval at which chunks are cleaned
	DefCacheChunkCleanInterval = cache.DefCacheChunkCleanInterval
	// DefCacheInfoAge is the default value for object info age
	DefCacheInfoAge = cache.DefCacheInfoAge
	// DefCacheReadRetries is the default value for read retries
	DefCacheReadRetries = cache.DefCacheReadRetries
	// DefCacheTotalWorkers is how many workers run in parallel to download chunks
	DefCacheTotalWorkers = cache.DefCacheTotalWorkers
	// DefCacheChunkNoMemory will enable or disable in-memory storage for chunks
	DefCacheChunkNoMemory = cache.DefCacheChunkNoMemory
	// DefCacheRps limits the number of requests per second to the source FS
	DefCacheRps = cache.DefCacheRps
	// DefCacheWrites will cache file data on writes through the cache
	DefCacheWrites = cache.DefCacheWrites
	// DefCacheTmpWaitTime says how long should files be stored in local cache before being uploaded
	DefCacheTmpWaitTime = cache.DefCacheTmpWaitTime
	// DefCacheDbWaitTime defines how long the cache backend should wait for the DB to be available
	DefCacheDbWaitTime = cache.DefCacheDbWaitTime
)

var legacyOptionMap = map[string]string{
	"chunk_size":           "data_chunk_size",
	"chunk_total_size":     "data_cache_max_size",
	"chunk_clean_interval": "data_cache_clean_interval",
	"chunk_path":           "data_cache_path",
	"chunk_no_memory":      "data_chunk_no_memory",
	"db_path":              "metadata_db_path",
	"db_purge":             "metadata_db_purge",
	"db_wait_time":         "metadata_db_wait_time",
	"info_age":             "metadata_cache_age",
	"read_retries":         "data_read_retries",
	"rps":                  "source_rps",
	"tmp_upload_path":      "temp_upload_path",
	"tmp_wait_time":        "temp_wait_time",
	"workers":              "data_workers",
	"writes":               "data_writes",
}

// Register with Fs
func init() {
	fs.Register(&fs.RegInfo{
		Name:        "cachev2",
		Description: "Cache a remote (v2)",
		NewFs:       NewFs,
		CommandHelp: commandHelp,
		Options: []fs.Option{{
			Name:     "remote",
			Help:     "Remote to cache. Should contain a ':' and a path, e.g. \"myremote:path/to/dir\".",
			Required: true,
		}, {
			Name: "plex_url",
			Help: "The URL of the Plex server.",
		}, {
			Name:      "plex_username",
			Help:      "The username of the Plex user.",
			Sensitive: true,
		}, {
			Name:       "plex_password",
			Help:       "The password of the Plex user.",
			IsPassword: true,
		}, {
			Name:      "plex_token",
			Help:      "The plex token for authentication - auto set normally.",
			Hide:      fs.OptionHideBoth,
			Advanced:  true,
			Sensitive: true,
		}, {
			Name:     "plex_insecure",
			Help:     "Skip all certificate verification when connecting to the Plex server.",
			Advanced: true,
		}, {
			Name:    "data_chunk_size",
			Help:    "The size of a chunk (partial file data).",
			Default: DefCacheChunkSize,
		}, {
			Name:    "metadata_cache_age",
			Help:    "How long to cache file structure information (directory listings, file size, times, etc.).",
			Default: DefCacheInfoAge,
		}, {
			Name:    "data_cache_max_size",
			Help:    "The total size that the cached chunks can take up on the local disk.",
			Default: DefCacheTotalChunkSize,
		}, {
			Name:     "metadata_db_path",
			Default:  filepath.Join(config.GetCacheDir(), "cache-backend"),
			Help:     "Directory to store file structure metadata DB. Defaults to the rclone cache dir.",
			Advanced: true,
		}, {
			Name:     "data_cache_path",
			Default:  filepath.Join(config.GetCacheDir(), "cache-backend"),
			Help:     "Directory to cache chunk files. Defaults to the rclone cache dir.",
			Advanced: true,
		}, {
			Name:     "metadata_db_purge",
			Default:  false,
			Help:     "Clear all the cached data for this remote on start.",
			Hide:     fs.OptionHideConfigurator,
			Advanced: true,
		}, {
			Name:     "data_cache_clean_interval",
			Default:  DefCacheChunkCleanInterval,
			Help:     "How often should the cache perform cleanups of the chunk storage.",
			Advanced: true,
		}, {
			Name:     "data_read_retries",
			Default:  DefCacheReadRetries,
			Help:     "How many times to retry a read from a cache storage.",
			Advanced: true,
		}, {
			Name:     "data_workers",
			Default:  DefCacheTotalWorkers,
			Help:     "How many workers should run in parallel to download chunks.",
			Advanced: true,
		}, {
			Name:     "data_chunk_no_memory",
			Default:  DefCacheChunkNoMemory,
			Help:     "Disable the in-memory cache for storing chunks during streaming.",
			Advanced: true,
		}, {
			Name:     "source_rps",
			Default:  int(DefCacheRps),
			Help:     "Limits the number of requests per second to the source FS (-1 to disable).",
			Advanced: true,
		}, {
			Name:     "data_writes",
			Default:  DefCacheWrites,
			Help:     "Cache file data on writes through the FS.",
			Advanced: true,
		}, {
			Name:     "temp_upload_path",
			Default:  "",
			Help:     "Directory to keep temporary files until they are uploaded.",
			Advanced: true,
		}, {
			Name:     "temp_wait_time",
			Default:  DefCacheTmpWaitTime,
			Help:     "How long should files be stored in local cache before being uploaded.",
			Advanced: true,
		}, {
			Name:     "metadata_db_wait_time",
			Default:  DefCacheDbWaitTime,
			Help:     "How long to wait for the DB to be available - 0 is unlimited.",
			Advanced: true,
		}},
	})
}

var warnLegacyKeysOnce sync.Once

// NewFs constructs an Fs from the path, container:path
func NewFs(ctx context.Context, name, rootPath string, m configmap.Mapper) (fs.Fs, error) {
	warnLegacyKeysOnce.Do(func() {
		legacyKeys := legacyOptionsUsed(m)
		if len(legacyKeys) == 0 {
			return
		}
		fs.Logf(name, "Detected legacy cache option keys for cachev2 (%s). Use cachev2 option names instead.", strings.Join(legacyKeys, ", "))
	})

	mapped := legacyMapper{mapper: m}
	if !hasOption(m, "metadata_db_path", "db_path") {
		mapped.Set("db_path", filepath.Join(config.GetCacheDir(), "cache-backend"))
	}
	if !hasOption(m, "data_cache_path", "chunk_path") {
		mapped.Set("chunk_path", filepath.Join(config.GetCacheDir(), "cache-backend"))
	}
	return cache.NewFs(ctx, name, rootPath, mapped)
}

func legacyOptionsUsed(m configmap.Mapper) []string {
	var legacy []string
	for legacyKey, newKey := range legacyOptionMap {
		if _, ok := m.Get(legacyKey); !ok {
			continue
		}
		if _, ok := m.Get(newKey); ok {
			continue
		}
		legacy = append(legacy, legacyKey)
	}
	sort.Strings(legacy)
	return legacy
}

func hasOption(m configmap.Mapper, newKey, legacyKey string) bool {
	if _, ok := m.Get(newKey); ok {
		return true
	}
	if _, ok := m.Get(legacyKey); ok {
		return true
	}
	return false
}

type legacyMapper struct {
	mapper configmap.Mapper
}

func (m legacyMapper) Get(key string) (value string, ok bool) {
	if newKey, found := legacyOptionMap[key]; found {
		if value, ok = m.mapper.Get(newKey); ok {
			return value, ok
		}
	}
	return m.mapper.Get(key)
}

func (m legacyMapper) Set(key, value string) {
	if newKey, found := legacyOptionMap[key]; found {
		m.mapper.Set(newKey, value)
		return
	}
	m.mapper.Set(key, value)
}

var commandHelp = []fs.CommandHelp{
	{
		Name:  "stats",
		Short: "Print stats on the cache backend in JSON format.",
	},
}
