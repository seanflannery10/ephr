# Vars
POSTGRES_USER = 'postgres'
POSTGRES_PASSWORD = 'test'
POSTGRES_DB = 'ephr'

# Extensions
load('ext://tests/golang', 'test_go')
load('ext://restart_process', 'docker_build_with_restart')

# Tests
test_go('test-ephr-cmd', './cmd/...', './cmd')
test_go('test-ephr-internal', './internal/...', './internal')

# Apply Migrations
local_resource(
  'migrations',
  cmd='dbmate --url postgres://{USER}:{PASS}@localhost:5432/{DB}?sslmode=disable up'.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB),
  resource_deps=['postgres']
)

# Build App
local_resource(
  'ephr-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o ./bin/ephr ./cmd/ephr',
   deps=['./cmd/ephr/', './internal/'],
)

# Run App
dockerfile='''
FROM golang:1.19.3-alpine
COPY /bin/ephr /
'''

docker_build_with_restart(
  'ephr-image',
  '.',
  entrypoint='/ephr',
  dockerfile_contents=dockerfile,
  only=[
    './bin/'
  ],
  live_update=[
    sync('./bin/', '/'),
  ],
)

ephr = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: ephr
  labels:
    run: ephr
data:
  DB_URL: 'postgres://{USER}:{PASS}@postgres:5432/{DB}?sslmode=disable'
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ephr
  labels:
    app: ephr
spec:
  selector:
    matchLabels:
      app: ephr
  template:
    metadata:
      labels:
        app: ephr
    spec:
      containers:
        - name: ephr
          image: ephr-image
          envFrom:
            - configMapRef:
                name: ephr
          ports:
            - containerPort: 4000
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(ephr))
k8s_resource('ephr', port_forwards=['4000'], resource_deps=['postgres', 'ephr-compile'])

# Run Debug App
dockerfile='''
FROM golang:1.19.3-alpine
ENV PORT="4001"
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY /bin/ephr /
'''

docker_build_with_restart(
  'ephr-debug-image',
  '.',
  entrypoint='dlv --listen=:4009 --headless=true --api-version=2 --accept-multiclient exec /ephr',
  dockerfile_contents=dockerfile,
  only=[
    './bin/'
  ],
  live_update=[
    sync('./bin/', '/'),
  ],
)

ephr_debug = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: ephr-debug
  labels:
    run: ephr-debug
data:
  DB_URL: 'postgres://{USER}:{PASS}@postgres:5432/{DB}?sslmode=disable'
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ephr-debug
  labels:
    app: ephr-debug
spec:
  selector:
    matchLabels:
      app: ephr-debug
  template:
    metadata:
      labels:
        app: ephr-debug
    spec:
      containers:
        - name: ephr-debug
          image: ephr-debug-image
          envFrom:
            - configMapRef:
                name: ephr-debug
          ports:
            - containerPort: 4001
            - containerPort: 4009
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(ephr_debug))
k8s_resource('ephr-debug', port_forwards=['4001', '4009'], resource_deps=['postgres', 'ephr-compile'])

# Run Logto
logto = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: logto
  labels:
    run: logto
data:
  TRUST_PROXY_HEADER: 'true'
  DB_URL: 'postgres://{USER}:{PASS}@postgres:5432/logto?sslmode=disable'
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logto
  labels:
    run: logto
spec:
  selector:
    matchLabels:
      run: logto
  template:
    metadata:
      labels:
        run: logto
    spec:
      containers:
        - name: logto
          image: ghcr.io/logto-io/logto:1.0.0-beta.14
          command: [ 'sh', '-c', 'sleep 3 && npm run cli db seed -- --swe && npm start' ]
          envFrom:
            - configMapRef:
                name: logto
          ports:
            - containerPort: 3001
---
apiVersion: v1
kind: Service
metadata:
  name: logto
  labels:
    run: logto
spec:
  ports:
  - port: 3001
    protocol: TCP
  selector:
    run: logto
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD)

k8s_yaml(blob(logto))
k8s_resource('logto', port_forwards=3001, resource_deps=['postgres'])

# Run Postgres
postgres = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres
  labels:
    run: postgres
data:
  POSTGRES_USER: {USER}
  POSTGRES_PASSWORD: {PASS}
  POSTGRES_DB: {DB}
  PGUSER: {USER}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  labels:
    run: postgres
spec:
  selector:
    matchLabels:
      run: postgres
  template:
    metadata:
      labels:
        run: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:15
          args:
            - postgres
            - -c
            - log_statement=all
          envFrom:
            - configMapRef:
                name: postgres
          ports:
            - containerPort: 5432
          startupProbe:
            exec:
              command:
                - /bin/sh
                - -c
                - exec pg_isready -h localhost
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  labels:
    run: postgres
spec:
  ports:
    - port: 5432
      protocol: TCP
  selector:
    run: postgres
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(postgres))
k8s_resource('postgres', port_forwards=5432)