load('ext://ko', 'ko_build')

ko_build('ephr-image',
    './cmd/ephr',
    deps=['./cmd/ephr', './internal'])

postgres_deploy = '''
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
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
        env:
        - name: POSTGRES_PASSWORD
          value: test
        - name: POSTGRES_DB
          value: ephr
        - name: PGUSER
          value: postgres
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
    resource_deps=['postgres'])

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