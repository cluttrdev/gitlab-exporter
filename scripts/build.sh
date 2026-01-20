#!/bin/sh
set -eu

REPO_ROOT="$(git rev-parse --show-toplevel)"
BIN_DIR="${REPO_ROOT}/bin"
DIST_DIR="${REPO_ROOT}/dist"
APPS=$(ls "${REPO_ROOT}/cmd/")
DEFAULT_BINARY_PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"
DEFAULT_IMAGE_PLATFORMS="linux/amd64 linux/arm64"

# Defaults
app=""
platform=""
tag=""
multiplatform=false
all=false
dist=false

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS] COMMAND

Commands:
  binary    Build application binary
  image     Build Docker image

Options:
  -a, --app APP         Target application
  -p, --platform PLAT   Target platform(s) - comma-separated os/arch pairs
  -t, --tag TAG         Image tag (default: version from git)
  --all                 Build for all apps and platforms
  --dist                Create distribution archives (binary only)
  --multiplatform       Use docker buildx for multi-arch images
  -h, --help            Show this help

Examples:
  $(basename "$0") binary -a gitlab-exporter
  $(basename "$0") binary --all --dist
  $(basename "$0") binary --all -p linux/amd64,linux/arm64
  $(basename "$0") image -a gitlab-exporter
  $(basename "$0") image -a gitlab-exporter -p linux/amd64,linux/arm64 --multiplatform
  $(basename "$0") image --all -p linux/amd64,darwin/arm64
  $(basename "$0") image --all --multiplatform -p linux/amd64,linux/arm64
EOF
}

# Get version from git tags
version() {
    git describe --tags --exact-match 2>/dev/null || \
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

# Create distribution archive for specified app and platform.
# Args: app, platform (os/arch)
create_dist() {
    _app="$1"
    _platform="$2"

    _os_arch=$(echo "${_platform}" | tr '/' '_')
    _ver="$(version)"

    mkdir -p "${DIST_DIR}"
    _archive="${_app}_${_ver}_${_os_arch}.tar.gz"

    echo "Creating ${_archive}..."
    tar -czf "${DIST_DIR}/${_archive}" -C "${BIN_DIR}/${_os_arch}" "${_app}"
    (cd "${DIST_DIR}" && sha256sum "${_archive}" > "${_archive}.sha256")
}

# Command is required
command="$1"
shift
if [ -z "${command}" ]; then
    echo "Error: command required" >&2
    usage >&2
    exit 1
fi

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
        --dist)
            dist=true
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

case "${command}" in
    binary)
        if [ "${all}" = true ]; then
            if [ -n "${platform}" ]; then
                _platforms=$(echo "${platform}" | tr ',' ' ')
            else
                _platforms="${DEFAULT_BINARY_PLATFORMS}"
            fi
            for a in ${APPS}; do
                for p in ${_platforms}; do
                    build_binary "${a}" "${p}"
                    if [ "${dist}" = true ]; then
                        create_dist "${a}" "${p}"
                    fi
                done
            done
        else
            if [ -z "${app}" ]; then
                echo "Error: --app required for single binary build" >&2
                exit 1
            fi
            if [ -z "${platform}" ]; then
                platform="$(go env GOOS)/$(go env GOARCH)"
            fi
            build_binary "${app}" "${platform}"
            if [ "${dist}" = true ]; then
                create_dist "${app}" "${platform}"
            fi
        fi
        ;;
    image)
        if [ "${all}" = true ]; then
            if [ "${multiplatform}" = true ]; then
                if [ -n "${platform}" ]; then
                    _platform="${platform}"
                else
                    _platform=$(echo "${DEFAULT_IMAGE_PLATFORMS}" | tr ' ' ',')
                fi
                for a in ${APPS}; do
                    build_image_multiplatform "${a}" "${_platform}" "${tag}"
                done
            else
                if [ -n "${platform}" ]; then
                    _platforms=$(echo "${platform}" | tr ',' ' ')
                else
                    _platforms="${DEFAULT_IMAGE_PLATFORMS}"
                fi
                for a in ${APPS}; do
                    for p in ${_platforms}; do
                        build_image "${a}" "${p}" "${tag}"
                    done
                done
            fi
        else
            if [ -z "${app}" ]; then
                echo "Error: --app required for single image build" >&2
                exit 1
            fi
            if [ -z "${platform}" ]; then
                platform="$(go env GOOS)/$(go env GOARCH)"
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
