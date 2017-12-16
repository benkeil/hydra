Angenommen:
- die Dockerfiles befinden sich in einem git repository
- wir haben das repository gerade mit 1.3.4 getaged

Wir erhalten folgende images:
- based on root:
  - private.docker.registry:1234/docker-common/nginx-base:1.3.4
  - private.docker.registry:1234/docker-common/nginx-base:1.3
  - private.docker.registry:1234/docker-common/nginx-base:1
  - private.docker.registry:1234/docker-common/nginx-base:latest