#!/bin/sh
V=$(git rev-list --count master)
V=$(expr ${V} + 1)
LINE=$(grep 'const Version int64' bas/core_init.go)
sed -i '' -e "s|${LINE}|const Version int64 = ${V}|" bas/core_init.go
