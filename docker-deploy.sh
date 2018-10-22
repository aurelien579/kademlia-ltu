sudo docker run -d -p 4000:4000 bootstrap_node

for i in {1..5}
do
    echo $i
	sudo docker run -d kademlia_node
done
