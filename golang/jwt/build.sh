#!/bin/bash -x

if [[ $# -lt 1 ]]; then
	echo "[FATAL]Must give version to build"
	exit 1
fi

case $1 in
	v3.0.9)
		rm -rf go.* vendor
		go mod init tyk_plugin
		go mod edit -replace github.com/jensneuse/graphql-go-tools=github.com/TykTechnologies/graphql-go-tools@6ff6aba4c612
		go get github.com/TykTechnologies/tyk@92e1486ab27653d059b23e02feee71683125ab68
		go mod tidy
		go mod vendor
		cp -rf ./.build/* ./vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v3.0.9 plugin.so
		;;
	v3.0.10-rc5)
		rm -rf go.* vendor
		go mod init tyk_plugin
		go mod edit -replace github.com/jensneuse/graphql-go-tools=github.com/TykTechnologies/graphql-go-tools@6ff6aba4c612
		go get github.com/TykTechnologies/tyk@b75124550124197d8ca34aac0868d793a67387d1
		go mod tidy
		go mod vendor
		cp -rf ./.build/* ./vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v3.0.10-rc5 plugin.so
		;;
	v4.0.0)
		rm -rf go.* vendor
		go mod init tyk_plugin
		go mod edit -replace github.com/jensneuse/graphql-go-tools=github.com/TykTechnologies/graphql-go-tools@6ff6aba4c612
		go get github.com/TykTechnologies/tyk@272b8ba0fcb3148ddcff8ac996ffd402d6e53933
		go mod tidy
		go mod vendor
		cp -rf ./.build/* ./vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v4.0.0 plugin.so
		;;
	*)
		echo "[FATAL]Unknown version"
		echo "       use v3.0.9 v3.0.10-rc5 v4.0.0"
		exit 1
		;;
esac
