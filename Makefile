MODE ?= development

all: frontend run

setup:
	cd frontend && npm install

.PHONY: frontend
frontend:
	cd $@ && npx webpack --mode $(MODE)

run:
	go1.16rc1 run main.go

build:
	go1.16rc1 build .

clean:
	rm -rf frontend/dist frontend/node_modules
