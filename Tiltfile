# -*- mode: Python -*-

load('ext://ko', 'ko_build')

docker_compose('docker-compose.yml')

ko_build('ephr-image',
         './cmd/ephr',
         deps=['./cmd/ephr', './internal'])