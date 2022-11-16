# -*- mode: Python -*-

load('ext://ko', 'ko_build')

docker_compose('docker-compose.yml')

#local_resource('dbmate up', cmd='make migrations')

ko_build('ephr-image',
         './cmd/ephr',
         deps=['./cmd/ephr', './internal'])