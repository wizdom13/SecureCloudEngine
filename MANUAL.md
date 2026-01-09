# rclone serve ftp
Serve remote:path over FTP.
Run a basic FTP server to serve a remote over FTP protocol.
This can be viewed with a FTP client or you can make a remote of
type FTP to read and write it.
## Server options
Use --addr to specify which IP address and port the server should
listen on, e.g. --addr 1.2.3.4:8000 or --addr :8080 to listen to all
IPs.  By default it only listens on localhost.  You can use port
:0 to let the OS choose an available port.
If you set --addr to listen on a public or LAN accessible IP address
then using Authentication is advised - see the next section for info.
### Authentication
By default this will serve files without needing a login.
You can set a single username and password with the --user and --pass flags.
## Auth Proxy

If you supply the parameter `--auth-proxy /path/to/program` then
rclone will use that program to generate backends on the fly which
then are used to authenticate incoming requests.  This uses a simple
JSON based protocol with input on STDIN and output on STDOUT.

**PLEASE NOTE:** `--auth-proxy` and `--authorized-keys` cannot be used
together, if `--auth-proxy` is set the authorized keys option will be
ignored.

There is an example program
[bin/test_proxy.py](https://github.com/rclone/rclone/blob/master/bin/test_proxy.py)
in the rclone source code.

The program's job is to take a `user` and `pass` on the input and turn
those into the config for a backend on STDOUT in JSON format.  This
config will have any default parameters for the backend added, but it
won't use configuration from environment variables or command line
options - it is the job of the proxy program to make a complete
config.

This config generated must have this extra parameter

- `_root` - root to use for the backend

And it may have this parameter

- `_obscure` - comma separated strings for parameters to obscure

If password authentication was used by the client, input to the proxy
process (on STDIN) would look similar to this:

```json
{
  "user": "me",
  "pass": "mypassword"
}
```

If public-key authentication was used by the client, input to the
proxy process (on STDIN) would look similar to this:

```json
{
  "user": "me",
  "public_key": "AAAAB3NzaC1yc2EAAAADAQABAAABAQDuwESFdAe14hVS6omeyX7edc...JQdf"
}
```

And as an example return this on STDOUT

```json
{
  "type": "sftp",
  "_root": "",
  "_obscure": "pass",
  "user": "me",
  "pass": "mypassword",
  "host": "sftp.example.com"
}
```

This would mean that an SFTP backend would be created on the fly for
the `user` and `pass`/`public_key` returned in the output to the host given.  Note
that since `_obscure` is set to `pass`, rclone will obscure the `pass`
parameter before creating the backend (which is required for sftp
backends).

The program can manipulate the supplied `user` in any way, for example
to make proxy to many different sftp backends, you could make the
`user` be `user@example.com` and then set the `host` to `example.com`
in the output and the user to `user`. For security you'd probably want
to restrict the `host` to a limited list.

Note that an internal cache is keyed on `user` so only use that for
configuration, don't use `pass` or `public_key`.  This also means that if a user's
password or public-key is changed the cache will need to expire (which takes 5 mins)
before it takes effect.

This can be used to build general purpose proxies to any kind of
backend that rclone supports.

rclone serve ftp remote:path [flags]
      --addr string                            IPaddress:Port or :Port to bind server to (default "localhost:2121")
      --auth-proxy string                      A program to use to create the backend from the auth
      --cert string                            TLS PEM key (concatenation of certificate and CA certificate)
  -h, --help                                   help for ftp
      --key string                             TLS PEM Private key
      --pass string                            Password for authentication (empty value allow every password)
      --passive-port string                    Passive port range to use (default "30000-32000")
      --public-ip string                       Public IP address to advertise for passive connections
      --user string                            User name for authentication (default "anonymous")
# rclone serve http
Serve the remote over HTTP.
Run a basic web server to serve a remote over HTTP.
This can be viewed in a web browser or you can make a remote of type
http read from it.

You can use the filter flags (e.g. `--include`, `--exclude`) to control what
is served.

The server will log errors.  Use `-v` to see access logs.

`--bwlimit` will be respected for file transfers.  Use `--stats` to
control the stats printing.
Use `--addr` to specify which IP address and port the server should
listen on, eg `--addr 1.2.3.4:8000` or `--addr :8080` to listen to all
If you set `--addr` to listen on a public or LAN accessible IP address
You can use a unix socket by setting the url to `unix:///path/to/socket`
or just by using an absolute path name.

`--addr` may be repeated to listen on multiple IPs/ports/sockets.
Socket activation, described further below, can also be used to accomplish the same.

`--server-read-timeout` and `--server-write-timeout` can be used to
control the timeouts on the server.  Note that this is the total time
for a transfer.

`--max-header-bytes` controls the maximum number of bytes the server will
accept in the HTTP header.

`--baseurl` controls the URL prefix that rclone serves from.  By default
rclone will serve from the root.  If you used `--baseurl "/rclone"` then
rclone would serve from a URL starting with "/rclone/".  This is
useful if you wish to proxy rclone serve.  Rclone automatically
inserts leading and trailing "/" on `--baseurl`, so `--baseurl "rclone"`,
`--baseurl "/rclone"` and `--baseurl "/rclone/"` are all treated
identically.

`--disable-zip` may be set to disable the zipping download option.

### TLS (SSL)

By default this will serve over http.  If you want you can serve over
https.  You will need to supply the `--cert` and `--key` flags.
If you wish to do client side certificate validation then you will need to
supply `--client-ca` also.

`--cert` must be set to the path of a file containing
either a PEM encoded certificate, or a concatenation of that with the CA
certificate. `--key` must be set to the path of a file
with the PEM encoded private key. If setting `--client-ca`,
it should be set to the path of a file with PEM encoded client certificate
authority certificates.

`--min-tls-version` is minimum TLS version that is acceptable. Valid
values are "tls1.0", "tls1.1", "tls1.2" and "tls1.3" (default "tls1.0").

## Socket activation

Instead of the listening addresses specified above, rclone will listen to all
FDs passed by the service manager, if any (and ignore any arguments passed
by `--addr`).

This allows rclone to be a socket-activated service.
It can be configured with .socket and .service unit files as described in
<https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html>.

Socket activation can be tested ad-hoc with the `systemd-socket-activate`command

```console
systemd-socket-activate -l 8000 -- rclone serve
```

This will socket-activate rclone on the first connection to port 8000 over TCP.

### Template

`--template` allows a user to specify a custom markup template for HTTP
and WebDAV serve functions.  The server exports the following markup
to be used within the template to server pages:

| Parameter   | Subparameter | Description |
| :---------- | :----------- | :---------- |
| .Name       |              | The full path of a file/directory. |
| .Title      |              | Directory listing of '.Name'. |
| .Sort       |              | The current sort used. This is changeable via '?sort=' parameter. Possible values: namedirfirst, name, size, time (default namedirfirst). |
| .Order      |              | The current ordering used. This is changeable via '?order=' parameter. Possible values: asc, desc (default asc). |
| .Query      |              | Currently unused. |
| .Breadcrumb |              | Allows for creating a relative navigation. |
|             | .Link        | The link of the Text relative to the root. |
|             | .Text        | The Name of the directory. |
| .Entries    |              | Information about a specific file/directory. |
|             | .URL         | The url of an entry. |
|             | .Leaf        | Currently same as '.URL' but intended to be just the name. |
|             | .IsDir       | Boolean for if an entry is a directory or not. |
|             | .Size        | Size in bytes of the entry. |
|             | .ModTime     | The UTC timestamp of an entry. |

The server also makes the following functions available so that they can be used
within the template. These functions help extend the options for dynamic
rendering of HTML. They can be used to render HTML based on specific conditions.

| Function   | Description |
| :---------- | :---------- |
| afterEpoch  | Returns the time since the epoch for the given time. |
| contains    | Checks whether a given substring is present or not in a given string. |
| hasPrefix   | Checks whether the given string begins with the specified prefix. |
| hasSuffix   | Checks whether the given string end with the specified suffix. |

You can either use an htpasswd file which can take lots of users, or
set a single username and password with the `--user` and `--pass` flags.

Alternatively, you can have the reverse proxy manage authentication and use the
username provided in the configured header with `--user-from-header`  (e.g., `--user-from-header=x-remote-user`).
Ensure the proxy is trusted and headers cannot be spoofed, as misconfiguration
may lead to unauthorized access.

If either of the above authentication methods is not configured and client
certificates are required by the `--client-ca` flag passed to the server, the
client certificate common name will be considered as the username.

Use `--htpasswd /path/to/htpasswd` to provide an htpasswd file.  This is
in standard apache format and supports MD5, SHA1 and BCrypt for basic
authentication.  Bcrypt is recommended.

To create an htpasswd file:

```console
touch htpasswd
htpasswd -B htpasswd user
htpasswd -B htpasswd anotherUser
```

The password file can be updated while rclone is running.

Use `--realm` to set the authentication realm.

Use `--salt` to change the password hashing salt from the default.
rclone serve http remote:path [flags]
      --addr stringArray                       IPaddress:Port or :Port to bind server to (default 127.0.0.1:8080)
      --allow-origin string                    Origin which cross-domain request (CORS) can be executed from
      --baseurl string                         Prefix for URLs - leave blank for root
      --client-ca string                       Client certificate authority to verify clients with
      --disable-zip                            Disable zip download of directories
  -h, --help                                   help for http
      --htpasswd string                        A htpasswd file - if not provided no authentication is done
      --max-header-bytes int                   Maximum size of request header (default 4096)
      --min-tls-version string                 Minimum TLS version that is acceptable (default "tls1.0")
      --pass string                            Password for authentication
      --realm string                           Realm for authentication
      --salt string                            Password hashing salt (default "dlPL2MqE")
      --server-read-timeout Duration           Timeout for server reading data (default 1h0m0s)
      --server-write-timeout Duration          Timeout for server writing data (default 1h0m0s)
      --template string                        User-specified template
      --user string                            User name for authentication
      --user-from-header string                User name from a defined HTTP header
# rclone serve nfs
Serve the remote as an NFS mount
Create an NFS server that serves the given remote over the network.
This implements an NFSv3 server to serve any rclone remote via NFS.
The primary purpose for this command is to enable the [mount
command](https://rclone.org/commands/rclone_mount/) on recent macOS versions where
installing FUSE is very cumbersome.
This server does not implement any authentication so any client will be
able to access the data. To limit access, you can use `serve nfs` on
the loopback address or rely on secure tunnels (such as SSH) or use
firewalling.
For this reason, by default, a random TCP port is chosen and the
loopback interface is used for the listening address by default;
meaning that it is only available to the local machine. If you want
other machines to access the NFS mount over local network, you need to
specify the listening address and port using the `--addr` flag.
Modifying files through the NFS protocol requires VFS caching. Usually
you will need to specify `--vfs-cache-mode` in order to be able to
write to the mountpoint (`full` is recommended). If you don't specify
VFS cache mode, the mount will be read-only.
`--nfs-cache-type` controls the type of the NFS handle cache. By
default this is `memory` where new handles will be randomly allocated
when needed. These are stored in memory. If the server is restarted
the handle cache will be lost and connected NFS clients will get stale
handle errors.
`--nfs-cache-type disk` uses an on disk NFS handle cache. Rclone
hashes the path of the object and stores it in a file named after the
hash. These hashes are stored on disk the directory controlled by
`--cache-dir` or the exact directory may be specified with
`--nfs-cache-dir`. Using this means that the NFS server can be
restarted at will without affecting the connected clients.
`--nfs-cache-type symlink` is similar to `--nfs-cache-type disk` in
that it uses an on disk cache, but the cache entries are held as
symlinks. Rclone will use the handle of the underlying file as the NFS
handle which improves performance. This sort of cache can't be backed
up and restored as the underlying handles will change. This is Linux
only. It requires running rclone as root or with `CAP_DAC_READ_SEARCH`.
You can run rclone with this extra permission by doing this to the
rclone binary `sudo setcap cap_dac_read_search+ep /path/to/rclone`.
`--nfs-cache-handle-limit` controls the maximum number of cached NFS
handles stored by the caching handler. This should not be set too low
or you may experience errors when trying to access files. The default
is `1000000`, but consider lowering this limit if the server's system
resource usage causes problems. This is only used by the `memory` type
cache.
To serve NFS over the network use following command:
```sh
rclone serve nfs remote: --addr 0.0.0.0:$PORT --vfs-cache-mode=full
This specifies a port that can be used in the mount command. To mount
the server under Linux/macOS, use the following command:
```sh
mount -t nfs -o port=$PORT,mountport=$PORT,tcp $HOSTNAME:/ path/to/mountpoint
Where `$PORT` is the same port number used in the `serve nfs` command
and `$HOSTNAME` is the network address of the machine that `serve nfs`
was run on.
If `--vfs-metadata-extension` is in use then for the `--nfs-cache-type disk`
and `--nfs-cache-type cache` the metadata files will have the file
handle of their parent file suffixed with `0x00, 0x00, 0x00, 0x01`.
This means they can be looked up directly from the parent file handle
is desired.
This command is only available on Unix platforms.
rclone serve nfs remote:path [flags]
      --addr string                            IPaddress:Port or :Port to bind server to
  -h, --help                                   help for nfs
      --nfs-cache-dir string                   The directory the NFS handle cache will use if set
      --nfs-cache-handle-limit int             max file handles cached simultaneously (min 5) (default 1000000)
      --nfs-cache-type memory|disk|symlink     Type of NFS handle cache to use (default memory)
