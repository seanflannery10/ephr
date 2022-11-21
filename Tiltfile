load('ext://ko', 'ko_build')
load('ext://tests/golang', 'test_go')

ko_build('ephr-image',
    './cmd/ephr',
    deps=['./cmd/ephr', './internal'])

test_go('test-ephr-cmd', './cmd/...', './cmd')
test_go('test-ephr-internal', './internal/...', './internal')

ephr_blob = '''
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
          env:
            - name: DB_URL
              value: 'postgres://postgres:test@postgres:5432/ephr?sslmode=disable'
          ports:
            - containerPort: 4000'''

k8s_yaml(blob(ephr_blob))
k8s_resource('ephr', port_forwards=4000, resource_deps=['postgres'])

logto_deploy = '''
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
          command: [ "sh", "-c", "sleep 3 && npm run cli db seed -- --swe && npm start" ]
          env:
            - name: TRUST_PROXY_HEADER
              value: "true"
            - name: DB_URL
              value: postgres://postgres:test@postgres:5432/logto?sslmode=disable
          ports:
            - containerPort: 3001'''

logto_service = '''
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
    run: logto'''

k8s_yaml(blob(logto_deploy))
k8s_yaml(blob(logto_service))
k8s_resource('logto', port_forwards=3001, resource_deps=['postgres'])

postgres_deploy = '''
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
          env:
            - name: PGUSER
              value: postgres
            - name: POSTGRES_DB
              value: ephr
            - name: POSTGRES_PASSWORD
              value: test
          ports:
            - containerPort: 5432
          startupProbe:
            exec:
              command:
                - /bin/sh
                - -c
                - exec pg_isready -h localhost
            periodSeconds: 5'''

postgres_service = '''
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
    run: postgres'''

k8s_yaml(blob(postgres_deploy))
k8s_yaml(blob(postgres_service))
k8s_resource('postgres', port_forwards=5432)

local_resource('migrations',
    cmd='dbmate --url postgres://postgres:test@localhost:5432/ephr?sslmode=disable up',
    resource_deps=['postgres']
)