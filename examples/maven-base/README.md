Angenommen:
- die Dockerfiles befinden sich in einem git repository
- wir haben das repository gerade mit 1.3.4 getaged

Wir erhalten folgende images:
- based on java-7:
  - private.docker.registry:1234/docker-common/maven-base:1.3.4-java-7
  - private.docker.registry:1234/docker-common/maven-base:1.3-java-7
  - private.docker.registry:1234/docker-common/maven-base:1-java-7
  - private.docker.registry:1234/docker-common/maven-base:java-7
- bases on java-8:
  - private.docker.registry:1234/docker-common/maven-base:1.3.4-java-8
  - private.docker.registry:1234/docker-common/maven-base:1.3-java-8
  - private.docker.registry:1234/docker-common/maven-base:1-java-8
  - private.docker.registry:1234/docker-common/maven-base:java-8
  - private.docker.registry:1234/docker-common/maven-base:latest
- bases on java-9:
  - private.docker.registry:1234/docker-common/maven-base:1.3.4-java-9
  - private.docker.registry:1234/docker-common/maven-base:1.3-java-9
  - private.docker.registry:1234/docker-common/maven-base:1-java-9
  - private.docker.registry:1234/docker-common/maven-base:java-9