# fsnotifyr

fsnotifyr is an efficient recursive file watcher that accepts [double-star globs](https://pkg.go.dev/github.com/bmatcuk/doublestar).

In builds on [fsnotify](https://pkg.go.dev/github.com/fsnotify/fsnotify) providing sorely needed recursive functionality while retaining it's excellent cross-platform support and stability.

## Value Proposition

There does not appear to be a good alternative for this. The goal is:

- efficient and flexible globbing
- good cross-platform support
- simple, intuitive API
- usable as a package, or a binary

<img src="jimenju.png" alt="Watch Tree" title="Watch Tree" />

> [!NOTE]
> In active development. Not suitable for production. YMMV
