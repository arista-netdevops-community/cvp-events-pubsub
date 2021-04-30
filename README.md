# Streaming CVP Events into a pub sub broker.
![Alt text](media/overall.jpg?raw=true "overall")

### What are CVP Events 

Events are streaming from Arista EOS switches, virtual appliances or containers.  Events can be anything from high cpu utilization to BGP neighbors down.  Events will stream into CVP with the Terminattr streaming agent which using gRPC to send data to CVP.  CVP as a state device will then store the data and reach out to alerting systems once a threashold has been reached. 

### Screenshot of an event.

An Example of a CVP event is high interface utlization which is shown below. 

![Alt text](media/events.jpg?raw=true "events")

Interface Exceeded Outbound Utilization Threshold        Interface outOctets bandwidth utilization (98.80548%) exceeded threshold of 98%
kafka writer: writing 1 messages to test (partition: 0)
Interface Exceeded Outbound Utilization Threshold        Interface outOctets bandwidth utilization (1300.97%) exceeded threshold of 95%
kafka writer: writing 1 messages to test (partition: 0)

### Streaming events to your favorite pub/sub system.

In this binary today Kafka is supported.  So any new events with the UPDATED message once this binary is ran will then send to a kafka topic within the config/data.yaml file.  

```
Interface Exceeded Outbound Utilization Threshold        Interface outOctets bandwidth utilization (1300.97%) exceeded threshold of 95%
kafka writer: writing 1 messages to test (partition: 0)
```

This allows messages/events to be pushed to a kafka bus and picked up by applications like ELK stack for further alerting.

### CVP resource API's

[The CVP resource API's](https://aristanetworks.github.io/cloudvision-apis/modeling/) are API's which expose CVP data to third party programs or applications.  Currently, [Events, inventory and tags](https://aristanetworks.github.io/cloudvision-apis/models/) can either be streamed or received to an end user/third party application.  These events API's are exposed as multiple gRPC services within CVP.  The protobufs are compiled and available through [cloudvision-go](https://github.com/aristanetworks/cloudvision-go) or [cloudvision-python.](https://github.com/aristanetworks/cloudvision-python)

### Compile / Run / Change YAML file. 

Build the binary if necessary 
```
cd cmd /
go build -o ../bin/cmd main 
```

Edit the config/data.yaml file
```
kafka_broker: "127.0.0.1"
kafka_topic: "test"
```

### Local install building with go 1.16 && Kafkacat.

For local testing run zookeeper and kafka 
```
docker create --name=zookeeper -it --privileged --net=host -e ZOOKEEPER_CLIENT_PORT=2181 confluentinc/cp-zookeeper

docker create --name=kafka -it --privileged --net=host -e KAFKA_ZOOKEEPER_CONNECT=127.0.0.1:2181 -e /KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:9092 -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 confluentinc

docker start zookeeper kafka 
```

Keeping the data.conf yaml file the same pointing to local host 

```
./main --server cvpaddress:8443 --yamlfile ../config/data.yaml  -username ansible -password ansible
```

Optionally install kafkacat to watch locally kafka messages. 
```
kafkacat -b 127.0.0.1:9092 -t test

Interface outOctets bandwidth utilization (98.73116%) exceeded threshold of 98%
Interface outOctets bandwidth utilization (1300.969%) exceeded threshold of 95%
Interface inOctets bandwidth utilization (98.7302%) exceeded threshold of 98%
Interface inOctets bandwidth utilization (1300.917%) exceeded threshold of 95%
```