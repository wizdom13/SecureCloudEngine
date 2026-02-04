//go:build !plan9 && !js && !race

package cachev2

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/rclone/rclone/backend/local"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/operations"
	"github.com/stretchr/testify/require"
)

func TestLegacyOptionMapping(t *testing.T) {
	m := configmap.Simple{}
	m.Set("data_chunk_size", "4M")
	m.Set("chunk_total_size", "50M")

	mapped := legacyMapper{mapper: m}
	value, ok := mapped.Get("chunk_size")
	require.True(t, ok)
	require.Equal(t, "4M", value)

	value, ok = mapped.Get("chunk_total_size")
	require.True(t, ok)
	require.Equal(t, "50M", value)

	mapped.Set("chunk_size", "6M")
	value, ok = m.Get("data_chunk_size")
	require.True(t, ok)
	require.Equal(t, "6M", value)
}

func TestCacheV2BasicOperations(t *testing.T) {
	ctx := context.Background()
	remoteDir := t.TempDir()
	cacheDir := t.TempDir()

	m := configmap.Simple{}
	m.Set("remote", remoteDir)
	m.Set("data_chunk_size", "1M")
	m.Set("data_cache_max_size", "10M")
	m.Set("data_workers", "1")
	m.Set("metadata_cache_age", "24h")
	m.Set("metadata_db_path", filepath.Join(cacheDir, "db"))
	m.Set("data_cache_path", filepath.Join(cacheDir, "chunks"))

	f, err := NewFs(ctx, "cachev2", "", m)
	require.NoError(t, err)

	_, err = operations.Rcat(ctx, f, "first.txt", io.NopCloser(bytes.NewBufferString("one")), time.Now(), nil)
	require.NoError(t, err)

	entries, err := f.List(ctx, "")
	require.NoError(t, err)
	require.Len(t, entries, 1)

	wrapper, ok := f.(fs.UnWrapper)
	require.True(t, ok)
	wrappedFs := wrapper.UnWrap()
	_, err = operations.Rcat(ctx, wrappedFs, "second.txt", io.NopCloser(bytes.NewBufferString("two")), time.Now(), nil)
	require.NoError(t, err)

	entries, err = f.List(ctx, "")
	require.NoError(t, err)
	require.Len(t, entries, 1)

	if flusher, ok := f.(fs.DirCacheFlusher); ok {
		flusher.DirCacheFlush()
	}

	entries, err = f.List(ctx, "")
	require.NoError(t, err)
	require.Len(t, entries, 2)
}
