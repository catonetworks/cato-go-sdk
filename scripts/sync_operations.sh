#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
sdk_root="$(cd -- "${script_dir}/.." && pwd)"
cli_root="${1:-${CLI_ROOT:-"${sdk_root}/../cato-cli"}}"
cli_root="$(cd -- "${cli_root}" && pwd)"

if [[ ! -d "${cli_root}/queryPayloads" ]]; then
	echo "CLI checkout has no queryPayloads directory: ${cli_root}" >&2
	exit 1
fi

checked_out_cli_sha="$(git -C "${cli_root}" rev-parse --verify HEAD^{commit})"
cli_commit_sha="${CLI_COMMIT_SHA:-${checked_out_cli_sha}}"
if [[ ! "${cli_commit_sha}" =~ ^[0-9a-f]{40}$ ]]; then
	echo "CLI_COMMIT_SHA must be a full 40-character commit SHA." >&2
	exit 1
fi
if [[ "${checked_out_cli_sha}" != "${cli_commit_sha}" ]]; then
	echo "CLI checkout ${checked_out_cli_sha} does not match CLI_COMMIT_SHA=${cli_commit_sha}." >&2
	exit 1
fi

dry_run_output="$(
	cd -- "${sdk_root}"
	go run ./cmd/gqlops import \
		--cli-root "${cli_root}" \
		--cli-commit-sha "${cli_commit_sha}" \
		--expected 0 \
		--dry-run
)"
printf '%s\n' "${dry_run_output}"

canonical_count="$(awk -F '[= ]' '/^canonical=/ { print $2; exit }' <<<"${dry_run_output}")"
if [[ ! "${canonical_count}" =~ ^[1-9][0-9]*$ ]]; then
	echo "Unable to determine canonical operation count from importer output." >&2
	exit 1
fi

expected_count="${EXPECTED_OPERATIONS:-${canonical_count}}"
if [[ ! "${expected_count}" =~ ^[1-9][0-9]*$ ]]; then
	echo "EXPECTED_OPERATIONS must be a positive integer." >&2
	exit 1
fi
if [[ "${expected_count}" != "${canonical_count}" ]]; then
	echo "EXPECTED_OPERATIONS=${expected_count} does not match CLI operation count ${canonical_count}." >&2
	exit 1
fi

make -C "${sdk_root}" operations-import \
	CLI_ROOT="${cli_root}" \
	CLI_COMMIT_SHA="${cli_commit_sha}" \
	EXPECTED_OPERATIONS="${expected_count}"
make -C "${sdk_root}" generate EXPECTED_OPERATIONS="${expected_count}"
make -C "${sdk_root}" generate-check EXPECTED_OPERATIONS="${expected_count}"
(
	cd -- "${sdk_root}"
	go test ./...
	go build ./...
)
