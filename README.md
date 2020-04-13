# Git-trace-analysis

Leverage recent implementation of Git Trace2 Event format, I build a Unix Socket Stream listener that convert Git Trace events into relevant git performance metrics.

Metrics should be categorized by repository and operation.

## TODO

- [X] Implement Unix Socket (Stream) Listener
- [X] Deserialize events json format
- [ ] Aggregate related events into session
- [ ] Export completed sessions into metrics
- [ ] Handle orphaned sessions

## Requirements

- OS which support Unix Socket (tested with MacOS + Centos)
- Rust
- Git (tested with 2.26.0)

## How to use

```
# Run this project
git clone ... git-trace-analysis
cd git_trace-analysis

# Start the project
cargo run

# Set the following git configurations
git config --global trace2.eventTarget af_unix:stream:/tmp/git_trace.sock
git config --global trace2.eventBrief true
```

# References

- [Git Trace2 Documentation](https://git-scm.com/docs/api-trace2)
- [Git Config Documentation](https://git-scm.com/docs/git-config#Documentation/git-config.txt-trace2eventTarget)
