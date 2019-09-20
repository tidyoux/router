MK := ./make/make.sh

all: server agent

server: web
	$(MK) server routersvr

web:
	gopherjs build -m github.com/tidyoux/router/server/web -o ./bin/static/main.min.js

agent:
	$(MK) agent routeragt

.PHONY: all server web agent
