#!/usr/bin/env bash

set -euo pipefail

# 1. Enforce that DATABASE_URL is supplied
if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "Error: DATABASE_URL environment variable is not set." >&2
  echo "Usage: DATABASE_URL=\"postgres://...\" ./scripts/migrate.sh [up|down|status]" >&2
  exit 1
fi

COMMAND="${1:-}"

if [[ "$COMMAND" != "up" && "$COMMAND" != "down" && "$COMMAND" != "status" ]]; then
  echo "Error: Invalid or missing command." >&2
  echo "Usage: DATABASE_URL=\"...\" $0 [up|down|status]" >&2
  exit 1
fi

# Define modules in explicit dependency order:
# 'up' runs: auth -> profile -> content
# 'down' runs in reverse: content -> profile -> auth
MODULES=("auth" "profile" "content")

run_migrate() {
  local module="$1"
  local action="$2"
  local step="${3:-}"
  local migration_path="migrations/${module}/"
  local tracking_table="schema_migrations_${module}"

  if [[ ! -d "$migration_path" ]]; then
    echo "Directory $migration_path does not exist, skipping..."
    return
  fi

  # Append or merge x-migrations-table to the connection string
  local db_url
  if [[ "$DATABASE_URL" == *"?"* ]]; then
    db_url="${DATABASE_URL}&x-migrations-table=${tracking_table}"
  else
    db_url="${DATABASE_URL}?x-migrations-table=${tracking_table}"
  fi

  echo "==> Running '$action' for module: [${module}]"

  if [[ "$action" == "status" ]]; then
    migrate -path "$migration_path" -database "$db_url" version
  else
    migrate -path "$migration_path" -database "$db_url" "$action" $step
  fi
}

case "$COMMAND" in
  up)
    for mod in "${MODULES[@]}"; do
      run_migrate "$mod" "up"
    done
    ;;

  down)
    # Roll back in exact reverse dependency order
    for (( idx=${#MODULES[@]}-1; idx>=0; idx-- )); do
      run_migrate "${MODULES[idx]}" "down" "1"
    done
    ;;

  status)
    for mod in "${MODULES[@]}"; do
      run_migrate "$mod" "status"
    done
    ;;
esac