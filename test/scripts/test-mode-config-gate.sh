#!/usr/bin/env bash
set -euo pipefail
ROLE="${1:?usage: test-mode-config-gate.sh <role>}"
YAML="$(cd "$(dirname "$0")/../standalone/etc/cofiswarm/${ROLE}" && pwd)/${ROLE}.yaml"
[[ -f "$YAML" ]] || { echo "missing $YAML"; exit 1; }
grep -qE 'dispatch_url:' "$YAML" || { echo "missing dispatch_url in $YAML"; exit 1; }
grep -qE 'slot_manager_url:' "$YAML" || { echo "missing slot_manager_url in $YAML"; exit 1; }
echo "ok: ${ROLE} config points at dispatch + slot-manager"
