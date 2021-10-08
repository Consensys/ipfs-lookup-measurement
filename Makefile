docker:
	cd controller; make agent && cp -vp agent ../node
	cd monitor; docker build -t ipfs-monitor .; cd ..
	cd node; docker build --no-cache -t ipfs-node .; cd ..
