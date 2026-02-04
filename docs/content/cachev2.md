---
title: "Cache v2"
description: "Rclone docs for cachev2 remote"
versionIntroduced: "v1.70"
status: Experimental
---

# {{< icon "fa fa-archive" >}} Cache v2

The `cachev2` backend is the successor to the deprecated `cache` backend. It
wraps another remote and provides:

- Metadata caching (directory listings, object info).
- Chunked read caching for partial file access.
- Optional write-through caching (`data_writes`).
- Optional temporary upload staging (`temp_upload_path`).
- Optional Plex integration (same `plex_*` options).

## Legacy cache backend inventory

The legacy cache backend registers itself as the `cache` provider in
`backend/cache/cache.go` with the following core capabilities and defaults:

- Metadata cache age (`info_age`): **6h**.
- Chunked read cache size (`chunk_size`): **5M**.
- Chunk cache max size (`chunk_total_size`): **10G**.
- Chunk cleanup interval (`chunk_clean_interval`): **1m**.
- Read retries (`read_retries`): **10**.
- Chunk workers (`workers`): **4**.
- In-memory chunks (`chunk_no_memory`): **false** (memory enabled).
- Source rate limit (`rps`): **-1** (disabled).
- Write-through caching (`writes`): **false**.
- Temp upload wait (`tmp_wait_time`): **15s**.
- DB wait time (`db_wait_time`): **1s**.
- Default cache paths: `--cache-dir` + `/cache-backend` for DB and chunks.
- Runtime RC commands: `cache/expire`, `cache/stats`, `cache/fetch`.

## Configuration

Create a cachev2 remote with `rclone config` or config file entries:

```text
[my-cachev2]
type = cachev2
remote = myremote:path
metadata_cache_age = 6h
data_chunk_size = 5M
data_cache_max_size = 10G
```

### Core options

- `remote`: wrapped remote to cache.
- `metadata_cache_age`: how long metadata stays valid.
- `data_chunk_size`: chunk size for cached reads.
- `data_cache_max_size`: upper limit for cached chunk data.
- `data_cache_path` / `metadata_db_path`: where cache metadata/chunks live.

### Upload staging

To buffer uploads locally and write back later:

```text
temp_upload_path = /path/to/tmp
temp_wait_time = 15s
```

### Rate limiting and retries

- `source_rps`: rate limit requests to the wrapped remote.
- `data_read_retries`: read retry count when cache data is missing.

## Migration from `cache`

Rename legacy keys to the cachev2 names:

| Cache (legacy) | Cachev2 |
| --- | --- |
| `chunk_size` | `data_chunk_size` |
| `chunk_total_size` | `data_cache_max_size` |
| `chunk_clean_interval` | `data_cache_clean_interval` |
| `chunk_path` | `data_cache_path` |
| `chunk_no_memory` | `data_chunk_no_memory` |
| `db_path` | `metadata_db_path` |
| `db_purge` | `metadata_db_purge` |
| `db_wait_time` | `metadata_db_wait_time` |
| `info_age` | `metadata_cache_age` |
| `read_retries` | `data_read_retries` |
| `rps` | `source_rps` |
| `tmp_upload_path` | `temp_upload_path` |
| `tmp_wait_time` | `temp_wait_time` |
| `workers` | `data_workers` |
| `writes` | `data_writes` |

`remote` and `plex_*` options remain unchanged.

## Commands

Cachev2 supports the same `rc` commands as the legacy backend:

- `cache/expire`: expire a cached file or directory.
- `cache/stats`: view cache statistics.
- `cache/fetch`: prefetch cached chunks.

## VFS-first workflows

For mount/serve workflows, prefer VFS caching instead of a cache backend. See
[VFS documentation](/vfs/) for details.
