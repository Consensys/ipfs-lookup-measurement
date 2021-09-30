build:
	go build cmd/controller/*

nodes-list.out :
	printf "localhost:5001 \n localhost:5001 \n localhost:5001 \n localhost:5001 \n " > nodes-list.out

itest: nodes-list.out
	./controller simple -l nodes-list.out
