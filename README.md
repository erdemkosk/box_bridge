# BoxBridge

## Overview

The `BoxBridge` library is a system that handles message processing between a Kafka broker and MongoDB using the Inbox/Outbox pattern. The system is designed to facilitate reliable message handling, where messages are first saved in MongoDB (Outbox) before being sent to Kafka, and incoming messages from Kafka are stored in MongoDB (Inbox) before processing.

The project also provides a flexible configuration mechanism for both Kafka and MongoDB, making it easy to set up and use with minimal configuration.

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


## Future Plans

While the BoxBridge project is in active development and functional in its current state, there are several areas that are planned for improvement to make it more suitable for production environments. These include:

1. **Multi-Database Support**:
   - The current implementation supports MongoDB, but future versions will include support for additional databases (such as PostgreSQL, MySQL, or Redis). This will enable BoxBridge to integrate with a wider range of systems and provide flexibility in handling messages in different environments.

2. **Retry Mechanism Enhancements**:
   - The current retry logic is basic and will be improved to include more advanced configurations, such as exponential backoff, error categorization, and customizable retry policies. This will enhance the system’s ability to recover from temporary failures without causing delays or data loss.

3. **Advanced Logging**:
   - Future versions will include more detailed and configurable logging, allowing users to easily monitor message processing, track errors, and generate useful diagnostic information for troubleshooting. Logs will include message IDs, timestamps, retries, and success/failure statuses.

4. **Message Prioritization**:
   - There are plans to implement message prioritization, where certain messages are handled with higher priority than others. This will be especially useful for time-sensitive or critical operations, ensuring that they are processed first.

5. **Transaction Management**:
   - To improve reliability, transaction management will be introduced to ensure that messages are processed atomically. If a message fails at any point (e.g., in the Outbox or Kafka), the system will ensure that no partial transactions are committed, preventing data corruption or inconsistency.

6. **Monitoring and Metrics**:
   - The integration of monitoring tools and metrics collection will be added to provide real-time insights into system performance. This will include message throughput, error rates, and resource utilization, enabling better monitoring and scalability management.

7. **Security Enhancements**:
   - BoxBridge will be enhanced to include features such as encryption of messages in transit and at rest, authentication, and authorization mechanisms to ensure that only authorized services can interact with the system.

8. **Better Fault Tolerance**:
   - Improvements to fault tolerance will be made to ensure that the system can handle a wide variety of failures, such as network issues, database unavailability, or Kafka failures, without causing service downtime or message loss.


The project is actively being developed to meet these goals, and these enhancements will be implemented to ensure that BoxBridge can operate seamlessly in a production environment, with increased scalability, reliability, and user configurability.


## Dependencies

- Kafka: For message production and consumption.
- MongoDB: For storing messages in the Outbox and Inbox collections.
- Go: The programming language used for the implementation.

---



## Project Goals

- **Understand Design Patterns**: Deep dive into commonly used design patterns and their applications in real-world scenarios.
- **Solve Design Pattern Challenges**: Address and solve practical problems using appropriate design patterns.
- **Demonstrate Proficiency**: Show a thorough grasp of design patterns by providing clear solutions and explanations.

## Design Patterns Covered

This project explores a range of design patterns, including but not limited to:

- **Creational Patterns**: Singleton, Factory Method, Abstract Factory, Builder, Prototype
- **Structural Patterns**: Adapter, Bridge, Composite, Decorator, Facade, Flyweight, Proxy
- **Behavioral Patterns**: Chain of Responsibility, Command, Interpreter, Iterator, Mediator, Memento, Observer, State, Strategy, Template Method, Visitor
- 
## Example Problems
1) [Printer Problem](https://github.com/erdemkosk/i-know-design-paterns/blob/master/src/problems/printer/README.md)

## How to Use

1.
    ```bash
   boxBridge := pkg.NewBoxBridge(pkg.NewConfig("localhost:27017", "", "", "", 0))

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

