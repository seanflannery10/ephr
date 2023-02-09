# Load ENV Settings
load('ext://dotenv', 'dotenv')
dotenv()

POSTGRES_USER = os.getenv('POSTGRES_USER')
POSTGRES_PASSWORD = os.getenv('POSTGRES_PASSWORD')
POSTGRES_DB = os.getenv('POSTGRES_DB')

# Tests
load('ext://tests/golang', 'test_go')
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
COPY /bin/ /
'''

load('ext://restart_process', 'docker_build_with_restart')
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
k8s_resource('ephr', port_forwards='4000', resource_deps=['postgres', 'ephr-compile'])

# Run App Migrations
migrations_dockerfile='''
FROM amacneil/dbmate
COPY /db/migrations/ /db/migrations/
'''

docker_build(
  'ephr-migrations-image',
  '.',
  dockerfile_contents=migrations_dockerfile,
  only=['./db/migrations/'],
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
        command: ["/bin/sh", "-c", 'dbmate down; dbmate up']
        envFrom:
          - configMapRef:
              name: ephr-migrations
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(ephr_migrations))
k8s_resource('ephr-migrations', resource_deps=['postgres'])

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
k8s_resource('postgres', port_forwards='5432')
