#!/bin/bash -x

if [[ $# -lt 1 ]]; then
	echo "[FATAL]Must give version to build"
	exit 1
fi

case $1 in
	v2.9.4.7)
		rm -rf go.* vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v2.9.4.7 plugin.so
		;;
	v2.9.4.8)
		rm -rf go.* vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v2.9.4.8 plugin.so
		;;
	v3.0.7)
		rm -rf go.* vendor
		go mod init tyk_plugin
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v3.0.7 plugin.so
		;;
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
	v3.0.10)
		rm -rf go.* vendor
		go mod init tyk_plugin
		go mod edit -replace github.com/jensneuse/graphql-go-tools=github.com/TykTechnologies/graphql-go-tools@6ff6aba4c612
		go get github.com/TykTechnologies/tyk@8c5aa0e886ac0a88976b975c3935496d398fa922
		go mod tidy
		go mod vendor
		cp -rf ./.build/* ./vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v3.0.10 plugin.so
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
	v3.2.3-rc6)
		rm -rf go.* vendor
		go mod init tyk_plugin
		go mod edit -replace github.com/jensneuse/graphql-go-tools=github.com/TykTechnologies/graphql-go-tools@140640759f4b
		go get github.com/TykTechnologies/tyk@279eb0ae2daae6ab8cb1aaa8f10c4994211c6d66
		go mod tidy
		go mod vendor
		cp -rf ./.build/* ./vendor
		docker container run -v `pwd`:/plugin-source --rm tykio/tyk-plugin-compiler:v3.2.3-rc6 plugin.so
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
		echo "       use v2.9.4.7 v2.9.4.8 v3.0.7 v3.0.9 v3.0.10 v3.2.3-rc6 v4.0.0"
		exit 1
		;;
esac
