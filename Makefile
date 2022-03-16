VERSION=0.0.3

docker-image: bin/bugsim
	make -C deploy/docker/bugsim VERSION=$(VERSION)

bin/bugsim: cmd/bugsim/main.go
	CGO_ENABLED=0 go build -o bin/bugsim -ldflags="-extldflags=-static" ./cmd/bugsim

helm-lint:
	docker run -it --rm -v $(PWD):/data quay.io/helmpack/chart-testing:v3.5.0 \
	  ct lint --charts /data/deploy/helm/cloud-native-demo \
	    --chart-repos grafana=https://grafana.github.io/helm-charts \
		--chart-repos linkerd=https://helm.linkerd.io/stable \
		--debug --validate-maintainers=false

kind-load: docker-image
	docker tag github.com/cypherfox/cloud-native-demo/bugsim:$(VERSION) localhost:5001/bugsim:$(VERSION)
	kind load docker-image localhost:5001/bugsim:$(VERSION)

.PHONEY: kind-load docker-image helm-lint