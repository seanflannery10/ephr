# -*- mode: Python -*-

load('ext://ko', 'ko_build')

ko_build('ephr-image',
         './cmd/ephr',
         deps=['./cmd/ephr', './internal'])

k8s_yaml('deployments/kubernetes.yaml')
k8s_resource('ephr', port_forwards=4000)