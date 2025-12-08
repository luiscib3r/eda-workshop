load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://helm_remote', 'helm_remote')

#=======================================================================
# Open Telemetry Collector
#=======================================================================
k8s_yaml('k8s.local/otel/config.yaml')
k8s_yaml('k8s.local/otel/jaeger.yaml')
k8s_yaml('k8s.local/otel/prometheus.yaml')
k8s_resource('jaeger', labels='otel', port_forwards=[16686])
k8s_resource('prometheus', labels='otel')

#===========================================================
# NATS
#===========================================================
helm_remote(
    'nats',
    repo_name='nats',
    repo_url='https://nats-io.github.io/k8s/helm/charts',
    values=['k8s.local/nats/values.yaml'],
    namespace='eda-workshop',
    set=['upgrade.install=true']
)

k8s_yaml('k8s.local/nats/config.yaml')
k8s_resource('nats', labels='nats', port_forwards=[4222])

#=======================================================================
# Storage
#=======================================================================
k8s_yaml('k8s.local/seaweedfs/config.yaml')
k8s_yaml('k8s.local/seaweedfs/master.yaml')
k8s_yaml('k8s.local/seaweedfs/volume.yaml')
k8s_yaml('k8s.local/seaweedfs/filer.yaml')
k8s_yaml('k8s.local/seaweedfs/s3.yaml')
k8s_resource('seaweedfs-master', labels='storage')
k8s_resource('seaweedfs-volume', labels='storage', resource_deps=['seaweedfs-master'])
k8s_resource('seaweedfs-filer', labels='storage', resource_deps=['seaweedfs-master', 'seaweedfs-volume'], port_forwards=[8888])
k8s_resource('seaweedfs-s3', labels='storage', resource_deps=['seaweedfs-filer'])

#=======================================================================
# Nginx for some CORS handling
#=======================================================================
k8s_yaml('k8s.local/nginx/config.yaml')
k8s_yaml('k8s.local/nginx/deployment.yaml')
k8s_resource('nginx', labels='nginx', port_forwards=[8333], resource_deps=['seaweedfs-s3'])

#=======================================================================
# Postgres
#=======================================================================
k8s_yaml('k8s.local/postgres/scripts.yaml')
k8s_yaml('k8s.local/postgres/secret.yaml')
k8s_yaml('k8s.local/postgres/deployment.yaml')
k8s_yaml('k8s.local/postgres/service.yaml')
k8s_resource('postgres', port_forwards=5432, labels='postgres')

#=======================================================================
# Frontend
#=======================================================================
docker_build(
    'frontend',
    context='./frontend',
    dockerfile='k8s.local/Dockerfile.web',
    ignore=['./frontend/dist/', './frontend/node_modules/'],
    live_update=[
        fall_back_on('./frontend/vite.config.ts'),
        sync('./frontend/', '/app/'),
        run('bun install --frozen-lockfile', trigger=[
            './frontend/package.json', './frontend/bun.lock'
        ]),
    ]
)

k8s_yaml('k8s.local/frontend/deployment.yaml')
k8s_resource('frontend', port_forwards=['5173:5173'], labels='frontend')

#=======================================================================
# Go Service Deployment Function
#=======================================================================
def deploy_service(service_name, main_path, port_forwards, resource_deps=[], labels=[]):
    build_name = '{}-build'.format(service_name)
    build_cmd = 'CGO_ENABLED=0 GOOS=linux go build -o ./build/{} -gcflags "-N -l" {}'.format(
        service_name, 
        main_path
    )

    if len(labels) == 0:
        labels = [service_name]

    local_resource(
        build_name, 
        build_cmd, 
        labels=labels,
        deps=['./backend'], 
    )

    docker_build_with_restart(
        service_name,
        context='./build',
        entrypoint=[
            '/dlv', 
            '--listen=:40000', 
            '--api-version=2', 
            '--headless=true', 
            '--accept-multiclient', 
            'exec', 
            '--continue', 
            '/app/{}'.format(service_name)
        ],
        dockerfile='k8s.local/Dockerfile.svc',
        build_args={
            'BINARY': service_name
        },
        only=service_name,
        live_update=[
            sync('./build/{}'.format(service_name), '/app/{}'.format(service_name))
        ],
    )

    
    k8s_yaml('k8s.local/{}/deployment.yaml'.format(service_name))
    
    k8s_resource(
        service_name, 
        labels=labels,
        resource_deps=[build_name] + resource_deps,
        port_forwards=port_forwards,
    )

# ===========================================================
# Backend Service
# ===========================================================
k8s_yaml('k8s.local/backend/config.yaml')
k8s_yaml('k8s.local/backend/service.yaml')
deploy_service(
    service_name='backend',
    main_path='./backend/cmd/api',
    port_forwards=['40000:40000', '8080:8080'],
    resource_deps=['nats', 'postgres', 'nginx']
)
# ===========================================================
# Telegram Notifier Service
# ===========================================================
deploy_service(
    service_name='telegram',
    main_path='./backend/cmd/telegram',
    port_forwards=['40001:40000'],
    resource_deps=['backend', 'ocr', 'ocr-image'],
    labels=['backend']
)

# ===========================================================
# OCR Service
# ===========================================================
deploy_service(
    service_name='ocr',
    main_path='./backend/cmd/ocr',
    port_forwards=['40002:40000'],
    resource_deps=['backend'],
    labels=['backend']
)

# ===========================================================
# OCR Image Service
# ===========================================================
docker_build(
    'ocr-image',
    context='./backend',
    dockerfile='k8s.local/Dockerfile.ocr',
)

k8s_yaml('k8s.local/ocr-image/deployment.yaml')
k8s_resource('ocr-image', labels='backend', resource_deps=['ocr'])