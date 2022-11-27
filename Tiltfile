load('ext://dotenv', 'dotenv')
load('ext://tests/golang', 'test_go')
load('ext://restart_process', 'docker_build_with_restart')

# Load ENV Settings
dotenv()

POSTGRES_USER = os.getenv('POSTGRES_USER')
POSTGRES_PASSWORD = os.getenv('POSTGRES_PASSWORD')
POSTGRES_DB = os.getenv('POSTGRES_DB')

# Tests
test_go('test-ephr-cmd', './cmd/...', './cmd')
test_go('test-ephr-internal', './internal/...', './internal')

# Build App
local_resource(
  'ephr-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o ./bin/ephr ./cmd/ephr',
   deps=['./cmd/ephr/', './internal/'],
)

# Run App
dockerfile='''
FROM alpine
COPY /bin/ephr /
'''

docker_build_with_restart(
  'ephr-image',
  '.',
  entrypoint='/ephr',
  dockerfile_contents=dockerfile,
  only=['./bin/'],
  live_update=[sync('./bin/', '/')],
)

ephr = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: ephr
  labels:
    app: ephr
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

# Run App Migrations
migrations_dockerfile='''
FROM golang:1.19.3-alpine AS builder
RUN CGO_ENABLED=0 go install github.com/amacneil/dbmate@latest

FROM alpine
COPY --from=builder /go/bin/dbmate /
COPY /db/migrations/ /db/migrations/
'''

docker_build(
  'ephr-migrations-image',
  '.',
  dockerfile_contents=migrations_dockerfile,
  only=['./db/migrations/'],
  live_update=[sync('./db/migrations/', '/db/migrations/')],
)

ephr_migrations = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: ephr-migrations
  labels:
    app: ephr-migrations
data:
  DATABASE_URL: 'postgres://{USER}:{PASS}@postgres:5432/{DB}?sslmode=disable'
---
apiVersion: batch/v1
kind: Job
metadata:
  name: ephr-migrations
  labels:
    app: ephr-migrations
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: ephr-migrations
        image: ephr-migrations-image
        command: ["/bin/sh", "-c", '/dbmate down; /dbmate up']
        envFrom:
          - configMapRef:
              name: ephr-migrations
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(ephr_migrations))
k8s_resource('ephr-migrations', resource_deps=['postgres'])

# Run Logto
logto = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: logto
  labels:
    app: logto
data:
  TRUST_PROXY_HEADER: 'true'
  DB_URL: 'postgres://{USER}:{PASS}@postgres:5432/logto?sslmode=disable'
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logto
  labels:
    app: logto
spec:
  selector:
    matchLabels:
      app: logto
  template:
    metadata:
      labels:
        app: logto
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
    app: logto
spec:
  ports:
  - port: 3001
    protocol: TCP
  selector:
    app: logto
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
    app: postgres
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
    app: postgres
spec:
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
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
    app: postgres
spec:
  ports:
    - port: 5432
      protocol: TCP
  selector:
    app: postgres
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(postgres))
k8s_resource('postgres', port_forwards=5432)