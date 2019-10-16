# Sensor Demo
This is a small application that is meant to demo use of serial devices with an "on-prem" Kubernetes cluster. It is not meant to be anything more than a proof of concept.

## Usage
There are Kubernetes manifests located in the `/manifests` directory located in the root of this repository. It includes a `kustomization.yml` that can be referenced as a base with git and then kustomized further to your liking.

## Docker Images
* [web](https://cloud.docker.com/u/phillebaba/repository/docker/phillebaba/sensor-demo-web)
* [client/server](https://cloud.docker.com/u/phillebaba/repository/docker/phillebaba/sensor-demo)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
