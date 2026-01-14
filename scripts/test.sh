#!/bin/sh
set -eu

REPO_ROOT="$(git rev-parse --show-toplevel)"
REPORTS_DIR="${REPO_ROOT}/reports"

# Defaults
mod=""
reports=false

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Run tests on Go modules.

Options:
  -m, --mod MOD    Target module (default: all modules)
  --reports        Generate JUnit and Cobertura reports
  -h, --help       Show this help

Examples:
  $(basename "$0")                    # Test all modules
  $(basename "$0") -m ./exporter      # Test specific module
  $(basename "$0") --reports          # Test all with reports
  $(basename "$0") --reports -m .     # Test root module with reports
EOF
}

# Find all modules in the repository
find_modules() {
    find . -type f -name go.mod -exec dirname {} \;
}

# Run tests for a single module
test_module() {
    _mod="$1"
    echo "Testing ${_mod}/..."
    go test -C "$_mod" ./...
}

# Run tests with reports for a single module
test_module_reports() {
    _mod="$1"
    _outdir="$_mod"

    echo "Testing ${_mod}/..."
    mkdir -p "${REPORTS_DIR}/${_outdir}"
    cd "$_mod"
    go tool -modfile="${REPO_ROOT}/go.tool.mod" gotestsum \
        --junitfile="${REPORTS_DIR}/${_outdir}/junit.xml" \
        --format=testname \
        -- \
        -coverprofile="${REPORTS_DIR}/${_outdir}/cover.out" -covermode=atomic ./...
    go tool -modfile="${REPO_ROOT}/go.tool.mod" gocover-cobertura \
        < "${REPORTS_DIR}/${_outdir}/cover.out" \
        > "${REPORTS_DIR}/${_outdir}/cobertura.xml"
    cd - > /dev/null
}

# Parse arguments
while [ $# -gt 0 ]; do
    case "$1" in
        -m|--mod)
            mod="$2"
            shift 2
            ;;
        --reports)
            reports=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

# Run tests
if [ -n "$mod" ]; then
    # Single module
    if [ "$reports" = true ]; then
        test_module_reports "$mod"
    else
        test_module "$mod"
    fi
else
    # All modules
    find_modules | while read -r m; do
        if [ "$reports" = true ]; then
            test_module_reports "$m"
        else
            test_module "$m"
        fi
    done
fi
