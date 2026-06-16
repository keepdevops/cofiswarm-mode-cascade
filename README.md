# cofiswarm-mode-cascade

Cofiswarm component: `mode-cascade`.

- Layout: [REPO-STANDARD-LAYOUT](https://github.com/keepdevops/cofiswarmdev/blob/main/docs/REPO-STANDARD-LAYOUT.md)
- Migration: [MIGRATION-SPRINTS](https://github.com/keepdevops/cofiswarmdev/blob/main/docs/MIGRATION-SPRINTS.md)

## FHS paths

| Path | Purpose |
|------|---------|
| `/etc/cofiswarm/mode-cascade/` | config |
| `/var/lib/cofiswarm/mode-cascade/` | state |
| `/var/log/cofiswarm/mode-cascade/` | logs |

## Test

```bash
./test/scripts/assert-layout.sh mode-cascade
```
