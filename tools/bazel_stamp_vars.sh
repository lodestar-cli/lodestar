# /bin/bash
cat << EOF
GIT_COMMIT $(git rev-parse --short HEAD)
EOF