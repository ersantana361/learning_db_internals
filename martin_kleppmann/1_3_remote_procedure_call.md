---
title: 1.3 Remote Procedure Call
draft: false
tags:
  - system-design
  - distributed-systems
  - martin-kleppmann
  - intermediate
  - reference
  - networking
  - protocols
---
### **Introduction**  
- **Title**: [Distributed Systems 1.3: RPC (Remote Procedure Call)](https://www.youtube.com/watch?v=S2osKiqQG9s)  
- **Overview**:  
  This lecture demystifies Remote Procedure Call (RPC) as a foundational abstraction for distributed systems, using real-world examples like online payment processing. Martin Kleppmann explains how RPC frameworks simulate local function calls across networks, addressing challenges like marshaling, fault tolerance, and interoperability. Key themes include the mechanics of stubs, historical RPC implementations (CORBA, gRPC), and modern practices (REST, microservices). The lecture progresses from practical use cases to technical intricacies, emphasizing how RPC enables scalable, cross-language communication in distributed architectures.

---

### **Chronological Analysis**  

#### **[RPC in Practice: Online Payment Example]**  
[Timestamp: 0:00–4:07](https://youtu.be/S2osKiqQG9s?t=0)  
> *"The payment service... is usually provided by a different company... RPC allows the online shop to send a message to this service."*  
> *"What looks like a local function call is translated into network communication via stubs."*  

**Analysis**:  
Kleppmann introduces RPC through an **online payment scenario**, where a shop (client) communicates with a third-party payment service (server). The client calls a `processPayment` function, which is actually a **stub** that marshals data (e.g., card details) into a network message. This segment contextualizes RPC as a mechanism for cross-organizational coordination, abstracting low-level networking into familiar programming constructs. The example underscores RPC’s role in modular systems, enabling services to specialize (e.g., payment processing) while maintaining interoperability. Challenges like **idempotency** (avoiding duplicate charges) are hinted at, linking to later discussions on fault tolerance.  

---

#### **[RPC Mechanics: Stubs and Marshaling]**  
[Timestamp: 4:07–9:55](https://youtu.be/S2osKiqQG9s?t=247)  
> *"The stub function... sends a message to the service. Marshaling translates programming language arguments into network messages."*  
> *"RPC frameworks handle encoding/decoding, allowing functions to return values across nodes."*  

**Analysis**:  
Here, **stubs** and **marshaling** (serialization) are explained as core RPC components. The client-side stub serializes arguments (e.g., JSON/Protocol Buffers), while the server-side stub deserializes them into native types. Kleppmann highlights **location transparency**—the illusion that remote functions behave like local ones. However, he contrasts this with reality: network calls introduce latency, partial failures, and non-determinism. This segment bridges theory (abstraction) and practice (network realities), setting the stage for fault-tolerance discussions.  

---

#### **[Challenges: Network Failures and Fault Tolerance]**  
[Timestamp: 9:55–16:28](https://youtu.be/S2osKiqQG9s?t=595)  
> *"Messages might be lost... How do we handle retries without charging a card twice?"*  
> *"RPC must address timeouts, crashes, and indeterminate states."*  

**Analysis**:  
Kleppmann dissects RPC’s **failure modes**: lost messages, delays, and server crashes. Unlike local calls, RPC requires handling **timeouts** and idempotent retries to prevent duplicate transactions (e.g., charging a card multiple times). This segment emphasizes the need for **idempotency tokens** or transactional safeguards in payment systems. The discussion ties back to earlier lectures on network unreliability, stressing that RPC’s simplicity masks underlying complexity, necessitating robust error handling in distributed systems.  

---

#### **[Historical and Modern RPC Implementations]**  
[Timestamp: 16:28–21:10](https://youtu.be/S2osKiqQG9s?t=988)  
> *"CORBA, Java RMI, gRPC... REST is often used but differs philosophically."*  
> *"REST leverages HTTP for interoperability but shares RPC’s core goal: invoking remote logic."*  

**Analysis**:  
The lecture reviews RPC’s evolution, from 1980s **CORBA** to modern **gRPC** and **REST**. While REST (using HTTP methods like POST/GET) is contrasted with traditional RPC, Kleppmann notes both enable remote execution. **gRPC** is highlighted for its efficiency (binary Protocol Buffers) and cross-language support. This segment contextualizes REST as a subset of RPC principles, optimized for web compatibility. The shift to **microservices** is noted, where RPC facilitates decoupled, polyglot services in enterprises.  

---

#### **[Interoperability and IDLs in Microservices]**  
[Timestamp: 21:10–27:35](https://youtu.be/S2osKiqQG9s?t=1270)  
> *"Interface Definition Languages (IDLs)... enable cross-language communication in microservices."*  
> *"Protocol Buffers define message types and services, generating stubs in multiple languages."*  

**Analysis**:  
Kleppmann introduces **Interface Definition Languages (IDLs)** like gRPC’s Protocol Buffers to solve cross-language compatibility. IDLs specify data types and service interfaces abstractly, allowing code generation for Java, Python, etc. This enables **microservices** in heterogeneous environments (e.g., legacy COBOL systems interacting with modern Go services). The example IDL for a payment service demonstrates structured type definitions, ensuring consistent marshaling/unmarshaling. This segment underscores RPC’s role in scalable, maintainable distributed systems.  

---

### **Conclusion**  
The lecture progresses from practical RPC use cases to technical depth, emphasizing:  
1. **Abstraction vs. Reality**: RPC simplifies distributed communication but must handle network failures.  
2. **Evolution**: From CORBA to gRPC, RPC adapts to technological shifts while maintaining core principles.  
3. **Interoperability**: IDLs and REST enable scalable, polyglot microservices architectures.  

**Practical Implications**: RPC is vital for modern systems (e.g., payment gateways, cloud services), requiring careful error handling. **Theoretical Importance**: The dichotomy between local and remote calls exposes fundamental distributed systems challenges (e.g., partial failure). By mastering RPC mechanics and tools like Protocol Buffers, developers can build resilient, cross-platform systems, balancing simplicity with the complexities of networked environments.