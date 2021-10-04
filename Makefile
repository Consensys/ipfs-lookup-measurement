docker:
	cd monitor; docker build -t ipfs-monitor .; cd ..
	cd node; docker build -t ipfs-node .; cd ..