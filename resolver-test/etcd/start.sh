docker run --name etcd-test \
	-p 2379:2379 \
	-p 2380:2380 \
	-v $(pwd)/etcd-data:/home/etcd/etcd-data \
	-v $(pwd)/etcd.conf.yaml:/home/etcd/etcd.conf.yaml \
	-e ALLOW_NONE_AUTHENTICATION=yes \
	bitnami/etcd etcd --config-file /home/etcd/etcd.conf.yaml