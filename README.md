# cofiswarm-mode-cascade

Mode plugin service (cascade) — HTTP execute endpoint for `cofiswarm-dispatch`.

- SDK: [cofiswarm-mode-sdk](../cofiswarm-mode-sdk)
- Legacy C++: `legacy/cpp/`
- Config gate: `dispatch_url` + `slot_manager_url` in `test/standalone/etc/cofiswarm/mode-cascade/mode-cascade.yaml`

## Run

```bash
make build
./bin/cofiswarm-mode-cascade -config test/standalone/etc/cofiswarm/mode-cascade/mode-cascade.yaml
```

Default listen: `:8023`
