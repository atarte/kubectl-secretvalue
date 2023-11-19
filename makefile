BINARY=kubectl-secretvalue
build:
	go build -o ${BINARY}

clean:
	rm ${BINARY}

deploy:
	VAR="default-value"
	kubectl create ns ns-test
	kubectl create secret generic my-secret -n ns-test --from-literal=my-key=$(VAR)

undeploy:
	kubectl delete secret my-secret -n ns-test
	kubectl delete ns ns-test

build_script:
	cp ./scripts/kubectl-secretvalue.sh ./${BINARY}
	chmod +x $(BINARY)