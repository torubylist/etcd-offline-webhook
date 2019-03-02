#etcd-offline-webhook

##Dependencies
- go 1.10.4
- kubernetes 1.9.6

## Test
- go test ./...

## Run
./gen-certs.sh
./ca-bundle.sh
kubectl apply -f manifest.yaml