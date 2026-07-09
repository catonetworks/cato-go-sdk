#!/usr/bin/env bash
set -euo pipefail

schema_file="${SCHEMA_FILE:-cato_api.graphqls}"
patch_dir="${PATCH_DIR:-schema-patches}"
schema_url="${SCHEMA_CURL_URL:-https://system.cc.catonetworks.com/api/schema?with_undocumented=true}"
allow_custom_schema_url="${ALLOW_CUSTOM_SCHEMA_URL:-false}"
manual_work_file="${SCHEMA_PATCH_MANUAL_WORK_FILE:-schema-patch-manual-work-needed.txt}"

validate_schema_url() {
	case "${schema_url}" in
		https://system.cc.catonetworks.com/api/schema*)
			return 0
			;;
		https://*)
			if [ "${allow_custom_schema_url}" = "true" ]; then
				return 0
			fi
			;;
	esac

	echo "Refusing schema URL outside the approved Cato schema endpoint." >&2
	echo "Set ALLOW_CUSTOM_SCHEMA_URL=true only for a trusted manual Jenkins run." >&2
	exit 1
}

triage_schema_patches() {
	if [ ! -d "${patch_dir}" ]; then
		echo "No patch directory found: ${patch_dir}"
		return 0
	fi

	shopt -s nullglob
	local patches=("${patch_dir}"/*.patch)
	shopt -u nullglob

	if [ "${#patches[@]}" -eq 0 ]; then
		echo "No schema patch files found in ${patch_dir}"
		return 0
	fi

	for patch in "${patches[@]}"; do
		echo "Checking schema patch: ${patch}"

		if git apply --reverse --check "${patch}" >/dev/null 2>&1; then
			echo "Patch ${patch} is now included in the upstream schema; deleting it."
			rm -f "${patch}"
			continue
		fi

		if git apply --check "${patch}" >/dev/null 2>&1; then
			echo "Patch applies cleanly: ${patch}"
			continue
		fi

		echo "Patch no longer applies cleanly: ${patch}" >&2
		echo "Regenerate this patch against the freshly fetched ${schema_file}, then rerun." >&2
		{
			echo "Schema patch requires manual work: ${patch}"
			echo
			echo "The latest schema was fetched and normalized, but this patch no longer applies cleanly."
			echo "Regenerate the patch context against ${schema_file}, verify it with:"
			echo
			echo "  git apply --check ${patch}"
			echo
			echo "Then rerun the schema update job."
		} > "${manual_work_file}"
		git apply --check "${patch}" >&2
		exit 1
	done
}

validate_schema_url
rm -f "${manual_work_file}"

echo "Fetching and normalizing schema from ${schema_url}"
SCHEMA_CURL_URL="${schema_url}" make schema-update

triage_schema_patches

echo "Applying schema patches"
make apply-patches

echo "Generating SDK client and models"
if ! make generate; then
	echo "Code generation failed." >&2
	echo "If this is an unmapped custom scalar, add a scalar implementation and .gqlgenc.yml mapping." >&2
	exit 1
fi

echo "Verifying Go build"
go build ./...
