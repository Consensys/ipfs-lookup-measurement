docker:
	printf "8twhpZqypAdqrbhD8feb mCl6L9aX5UaokvwxCLcM\n" > node/.key
	cp -p node/.key controller/.key
	cd controller; make agent && cp -vp agent ../node
	cd monitor; docker build -t ipfs-monitor .; cd ..
	cd node; docker build --no-cache -t ipfs-node .; cd ..
