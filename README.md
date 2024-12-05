# BoxBridge

## Overview

The `BoxBridge` library is a system that handles message processing between a Kafka broker and MongoDB using the Inbox/Outbox pattern. The system is designed to facilitate reliable message handling, where messages are first saved in MongoDB (Outbox) before being sent to Kafka, and incoming messages from Kafka are stored in MongoDB (Inbox) before processing.

The project also provides a flexible configuration mechanism for both Kafka and MongoDB, making it easy to set up and use with minimal configuration.

BoxBridge will automatically create the necessary Inbox and Outbox collections in MongoDB. These collections will serve as storage for messages, providing a mechanism for tracking and ensuring reliable message delivery. The system is designed to handle retries and failures, ensuring that messages are not lost during transmission between MongoDB and Kafka.

## Project Goals

The main goal of the BoxBridge project is to create a robust system that facilitates reliable and scalable message processing between Kafka and MongoDB. The system utilizes the Inbox/Outbox pattern, ensuring that messages are first stored in MongoDB before being sent to Kafka, and that incoming messages from Kafka are stored in MongoDB before being processed.

Key objectives of the project include:

1. **Reliable Message Handling**:
   - Ensure that all messages are reliably saved in MongoDB's Outbox before being sent to Kafka.
   - Use MongoDB's Inbox to temporarily store incoming messages from Kafka for processing.

2. **Simplified Configuration**:
   - Provide an easy-to-use configuration mechanism for both Kafka and MongoDB, so that users can quickly set up and use the system without requiring complex setups.

3. **Customizable Message Processing**:
   - Allow users to define custom logic for handling incoming messages by wrapping the Kafka consumer handler function, ensuring that each message is processed according to specific business rules.

4. **Flexible Integration**:
   - Enable the integration of Kafka and MongoDB into microservice architectures, where messages are reliably produced, consumed, and processed between services.

5. **Fault Tolerance**:
   - Ensure the system can handle message delivery failures, such as retries when saving to MongoDB or sending messages to Kafka, and provide meaningful error handling.

6. **Automatic Resource Management**:
   - Handle the initialization, operation, and shutdown of Kafka and MongoDB resources automatically, allowing for smooth and efficient resource management.

7. **Scalability**:
   - Design the system to be scalable, capable of handling increasing message volumes by efficiently managing Kafka consumers and producers.

By achieving these goals, the BoxBridge project aims to provide a flexible, easy-to-use framework for reliable message processing between Kafka and MongoDB, which can be seamlessly integrated into various service architectures.


## Limitations For Now

1. **Multi-Database Support**:
   - The current implementation supports MongoDB, but future versions will include support for additional databases (such as PostgreSQL, MySQL, or Redis). This will enable BoxBridge to integrate with a wider range of systems and provide flexibility in handling messages in different environments.

2. **Work On just one kafka cluster**:
   - Need to support multiple kafka cluster at once for producer and consumer.


## How to Use

```bash
  boxBridge := pkg.NewBoxBridge(pkg.NewConfigBuilder().
		WithMongoDBURL("mongodb://localhost:27017").
		WithKafkaURL("localhost:9092").
		WithOutboxCollection("outbox"). //default as outbox if u need use this and change it
		WithInboxCollection("inbox"). //default as inbox if u need use this and change it
		WithRetryAttempts(3).
		Build())

	producerConfig := model.ProducerConfig{
		TopicName: "my-topic",
		ClientID:  "my-producer",
	}

	boxBridge.AddProducer(producerConfig)

	boxBridge.Produce(producerConfig, "key1", Foo{
		Name: "Erdem Köşk",
		Note: "Hey , box-bridge is amazing mate!",
	})

	// Create handler function for consumer each or one
	handlerFunc := func(msg *model.KafkaMessage) error {
		log.Printf("Received message: %s", string(msg.Value))
		return nil
	}

	consumerConfig := model.ConsumerConfig{
		TopicName:   "my-topic",
		GroupID:     "my-consumer-group",
		HandlerFunc: handlerFunc,
	}

	boxBridge.AddConsumer(consumerConfig)
```
 **Console For Handle Consumer With Success:**
 ```bash
  2024/12/01 22:21:47 BOX-BRIDGE: MongoDB connection successfully established.
%4|1733080907.611|CONFWARN|rdkafka#producer-1| [thrd:app]: Configuration property group.id is a consumer property and will be ignored by this producer instance
%4|1733080907.611|CONFWARN|rdkafka#producer-1| [thrd:app]: Configuration property auto.offset.reset is a consumer property and will be ignored by this producer instance
2024/12/01 22:21:47 BOX-BRIDGE: Producer for topic my-topic initialized
2024/12/01 22:21:47 BOX-BRIDGE: Message delivered to my-topic[0]@11
2024/12/01 22:21:47 Message successfully sent to Kafka!
2024/12/01 22:21:47 BOX-BRIDGE: Consumer for topic my-topic started
2024/12/01 22:21:50 Received message: {"name":"Erdem Köşk","note":"Hey , box-bridge is amazing mate!"}
2024/12/01 22:21:57 Shutting down Kafka manager
2024/12/01 22:21:57 BOX-BRIDGE: Closing consumer for topic my-topic
2024/12/01 22:21:57 BOX-BRIDGE: Closing producer for topic my-topic
2024/12/01 22:21:57 BOX-BRIDGE: Kafka Manager shut down gracefully
2024/12/01 22:21:58 BOX-BRIDGE: MongoDB connection successfully closed.
```   

 **Console For Cannot Handle Consumer Message So It Wont Commited:**
 ```bash
2024/12/05 23:56:16 BOX-BRIDGE: Collection 'inbox' correlationId index created successfully
2024/12/05 23:56:16 BOX-BRIDGE: Collection 'outbox' correlationId index created successfully
2024/12/05 23:56:16 BOX-BRIDGE: MongoDB connection successfully established.
%4|1733432176.936|CONFWARN|rdkafka#producer-1| [thrd:app]: Configuration property group.id is a consumer property and will be ignored by this producer instance
%4|1733432176.936|CONFWARN|rdkafka#producer-1| [thrd:app]: Configuration property enable.auto.commit is a consumer property and will be ignored by this producer instance
%4|1733432176.936|CONFWARN|rdkafka#producer-1| [thrd:app]: Configuration property auto.offset.reset is a consumer property and will be ignored by this producer instance
2024/12/05 23:56:16 BOX-BRIDGE: Producer for topic my-topic initialized
2024/12/05 23:56:17 BOX-BRIDGE: Message delivered to my-topic[0]@5
2024/12/05 23:56:17 Message successfully sent to Kafka!
2024/12/05 23:56:17 BOX-BRIDGE: Consumer for topic my-topic started
2024/12/05 23:56:20 Received message: {"name":"Erdem Köşk","note":"Hey , box-bridge is amazing mate!"}
2024/12/05 23:56:20 Error in handler function for message: something went wrong , it should not commit any offset!
2024/12/05 23:56:27 Shutting down Kafka manager
2024/12/05 23:56:27 BOX-BRIDGE: Closing consumer for topic my-topic
2024/12/05 23:56:27 BOX-BRIDGE: Closing producer for topic my-topic
2024/12/05 23:56:27 BOX-BRIDGE: Kafka Manager shut down gracefully
2024/12/05 23:56:27 BOX-BRIDGE: MongoDB connection successfully closed.
```  


**Auto Created Collections**
![Auto Created Collections](https://i.imgur.com/8W5J0ek.png)
**Outbox Collection**
![Outbox](https://i.imgur.com/nYn2CK5.png)
**Inbox Collection**
![Inbox](https://i.imgur.com/gXoH5R5.png)
**Receive message but if u cannot processed it wont commited on inbox**
![Inbox](https://i.imgur.com/yOpyMOB.png)
**Receive message and processed it commited and mark as Processed**
![Inbox](https://i.imgur.com/CAMyRnJ.png)

