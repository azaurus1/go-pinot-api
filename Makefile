


.PHONY: kafka
kafka:
	docker compose -f docker-compose.yaml up -d kafka schema-registry


.PHONY: pinot
pinot:
	docker compose -f docker-compose.yaml up -d pinot-controller pinot-broker pinot-server pinot-zookeeper


.PHONY: generate-data
generate-data:
	cd ./example/data-gen && \
	go run generate_data.go && \
	cd ../..


.PHONY: consume-data
consume-data:
	kcat -b localhost:29092 \
		-t ethereum.mainnet.blocks -C \
        -s value=avro -r http://localhost:8081 \
        -f '\nKey (%K bytes): %k\t\n%s\nTimestamp: %T\tPartition: %p\tOffset: %o\n--\n' \
        -o beginning

