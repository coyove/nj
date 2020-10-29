#!/bin/sh
V=$(git rev-list --count master)
V=$(expr ${V} + 1)
LINE=$(grep 'const Version int64' lib.go)
sed -i '' -e "s|${LINE}|const Version int64 = ${V}|" lib.go
