#!/bin/sh
set -eu

REPO_ROOT="$(git rev-parse --show-toplevel)"
BIN_DIR="${REPO_ROOT}/bin"
APPS=$(ls "${REPO_ROOT}/cmd/")
DEFAULT_PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

# Defaults
app=""
platform=""
tag=""
multiplatform=false
all=false

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS] COMMAND

Commands:
  binary    Build application binary
  image     Build Docker image

Options:
  -a, --app APP         Target application
  -p, --platform PLAT   Target platform (default: current os/arch)
  -t, --tag TAG         Image tag (default: version from git)
  --multiplatform       Use docker buildx for multi-arch images
  --all                 Build for all apps and platforms
  -h, --help            Show this help

Examples:
  $(basename "$0") binary -a gitlab-exporter
  $(basename "$0") binary --all
  $(basename "$0") image -a gitlab-exporter
  $(basename "$0") image -a gitlab-exporter -p linux/amd64,linux/arm64 --multiplatform
EOF
}

# Get version from git tags
version() {
    git describe --exact-match 2>/dev/null || \
        echo "$(git describe --tags --abbrev=0)-dev.$(git rev-list --count "$(git describe --tags --abbrev=0)"..HEAD)+$(git rev-parse --short=8 HEAD)"
}

# Build binary for specified app and platform.
# Args: app, platform (os/arch)
build_binary() {
    _app="$1"
    _platform="$2"

    _goos="${_platform%/*}"
    _goarch="${_platform#*/}"
    _ver="$(version)"

    echo "Building ${_app} for ${_platform}..."
    CGO_ENABLED=0 GOOS="${_goos}" GOARCH="${_goarch}" \
    go build \
        -C "${REPO_ROOT}/cmd/${_app}" \
        -ldflags "-s -w -X 'main.version=${_ver}'" \
        -o "${BIN_DIR}/${_goos}_${_goarch}/" \
        .
}

# Build Docker image for specified app and platform.
# Args: app, platform (os/arch), tag
build_image() {
    _app="$1"
    _platform="$2"
    _tag="$3"

    _os="${_platform%/*}"
    _arch="${_platform#*/}"

    if [ -z "${_tag}" ]; then
        _tag="$(version | tr '+' '-')"
    fi

    if [ ! -f "${BIN_DIR}/${_os}_${_arch}/${_app}" ]; then
        echo "Binary ${_os}_${_arch}/${_app} not found! Run 'build.sh binary -a ${_app} -p ${_platform}' first." >&2
        exit 1
    fi

    echo "Building image ${_app}:${_tag} for ${_platform}..."
    docker build \
        --file "${REPO_ROOT}/Dockerfile" \
        --platform "${_platform}" \
        --build-arg APP="${_app}" \
        --tag "${_app}:${_tag}" \
        "${BIN_DIR}"
}

# Build multi-platform Docker image for specified app.
# Args: app, platform (comma-separated os/arch), tag
build_image_multiplatform() {
    _app="$1"
    _platform="$2"
    _tag="$3"

    if [ -z "${_tag}" ]; then
        _tag="$(version | tr '+' '-')"
    fi

    # Verify all binaries exist
    for _plat in $(echo "${_platform}" | tr ',' ' '); do
        _os="${_plat%/*}"
        _arch="${_plat#*/}"
        if [ ! -f "${BIN_DIR}/${_os}_${_arch}/${_app}" ]; then
            echo "Binary ${_os}_${_arch}/${_app} not found! Run 'build.sh binary -a ${_app} -p ${_plat}' first." >&2
            exit 1
        fi
    done

    _image="${_app}:${_tag}"
    echo "Building multiplatform image ${_image} for ${_platform}..."
    docker buildx build \
        --file "${REPO_ROOT}/Dockerfile" \
        --platform "${_platform}" \
        --build-arg APP="${_app}" \
        --output "type=image,name=${_image},push=false" \
        "${BIN_DIR}"
}

# Parse arguments
while [ $# -gt 0 ]; do
    case "$1" in
        -a|--app)
            app="$2"
            shift 2
            ;;
        -p|--platform)
            platform="$2"
            shift 2
            ;;
        -t|--tag)
            tag="$2"
            shift 2
            ;;
        --multiplatform)
            multiplatform=true
            shift
            ;;
        --all)
            all=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        -*)
            echo "Unknown option: $1" >&2
            usage >&2
            exit 1
            ;;
        *)
            break
            ;;
    esac
done

# Command is required
if [ $# -eq 0 ]; then
    echo "Error: command required" >&2
    usage >&2
    exit 1
fi

command="$1"
shift

# Set default platform if not specified
if [ -z "${platform}" ]; then
    platform="$(go env GOOS)/$(go env GOARCH)"
fi

case "${command}" in
    binary)
        if [ "${all}" = true ]; then
            for a in ${APPS}; do
                for p in ${DEFAULT_PLATFORMS}; do
                    build_binary "${a}" "${p}"
                done
            done
        else
            if [ -z "${app}" ]; then
                echo "Error: --app required for single binary build" >&2
                exit 1
            fi
            build_binary "${app}" "${platform}"
        fi
        ;;
    image)
        if [ "${all}" = true ]; then
            if [ "${multiplatform}" = true ]; then
                _platform=$(echo "${platform}" | tr ' ' ',')
                # If single platform specified, use default multi-platform set
                case "${_platform}" in
                    *,*) ;;  # already has comma, keep it
                    *)   _platform="linux/amd64,linux/arm64" ;;
                esac
                for a in ${APPS}; do
                    build_image_multiplatform "${a}" "${_platform}" "${tag}"
                done
            else
                for a in ${APPS}; do
                    build_image "${a}" "${platform}" "${tag}"
                done
            fi
        else
            if [ -z "${app}" ]; then
                echo "Error: --app required for single image build" >&2
                exit 1
            fi
            if [ "${multiplatform}" = true ]; then
                build_image_multiplatform "${app}" "${platform}" "${tag}"
            else
                build_image "${app}" "${platform}" "${tag}"
            fi
        fi
        ;;
    *)
        echo "Unknown command: ${command}" >&2
        usage >&2
        exit 1
        ;;
esac
