---
title: 7.3 Eventual consistency
draft: false
tags:
  - system-design
  - distributed-systems
  - martin-kleppmann
  - intermediate
  - reference
  - consistency
---
### **Introduction**  
- **Title**: [Distributed Systems 7.3: Eventual Consistency](https://www.youtube.com/watch?v=9uCP3qHNbWw)  
- **Overview**:  
  This video explores the trade-offs between consistency models in distributed systems, focusing on **eventual consistency** as an alternative to linearizability. It begins by critiquing linearizability’s limitations—high latency, scalability bottlenecks, and availability issues during network partitions—and introduces eventual consistency through practical examples and theoretical frameworks like the CAP theorem. The presentation transitions from strong consistency models to weaker guarantees, emphasizing availability, conflict resolution, and convergence properties. Key themes include the interplay between system assumptions (e.g., quorum communication), real-world applications (e.g., calendar synchronization), and the role of commutative operations in achieving consistency.

---

### **Chronological Analysis**  

#### **[Limitations of Linearizability]**  
[Timestamp: 0:00](https://youtu.be/9uCP3qHNbWw?t=0)  
> *"Linearizability is a very strong consistency model... [but] it’s actually quite expensive to implement."*  
> *"If you can't contact a quorum of nodes... neither reading nor writing is possible."*  

**Analysis**:  
- **Technical Explanation**: Linearizability ensures all operations appear atomic and ordered, masking replication. However, it requires protocols like Raft, which involve leader-centric sequencing and quorum communication (2 round trips per operation). This creates bottlenecks and limits scalability.  
- **Context**: The video positions linearizability as impractical for systems requiring high availability during network partitions. The reliance on quorums means operations halt if nodes are unreachable.  
- **Significance**: This segment establishes why distributed systems often need weaker models. Linearizability’s trade-offs—consistency vs. availability—motivate the exploration of eventual consistency.  
- **Real-World Implications**: Financial systems or real-time databases might prioritize linearizability, but applications like collaborative tools or offline-first apps cannot tolerate its downtime.

These consistency challenges become particularly complex in partitioned systems where cross-shard coordination is required. For practical implementation of eventual consistency patterns in partitioned architectures, see [[System Design/data_partitioning_and_sharding]].  

---

#### **[Calendar Example and CAP Theorem]**  
[Timestamp: 2:12](https://youtu.be/9uCP3qHNbWw?t=132)  
> *"In the presence of a network partition... you have to choose between linearizability and availability."*  
> *"The CAP theorem... illustrates a fundamental choice during partitions."*  

**Analysis**:  
- **Technical Explanation**: The calendar app example demonstrates concurrent updates during a partition. Conflicts arise (e.g., title vs. time changes), resolved via "last writer wins." The CAP theorem formalizes this trade-off: during partitions, systems must choose consistency (C) or availability (A).  
- **Context**: The example humanizes eventual consistency, showing how disconnected replicas can operate independently. The CAP theorem underpins this by explaining why strict consistency is unattainable during partitions.  
- **Significance**: Highlights the inevitability of conflicts in distributed systems and the need for resolution strategies.  
- **Real-World Applications**: Calendar/synchronization tools (Google Calendar, Dropbox) use similar conflict resolution, prioritizing user experience over strict consistency.  

---

#### **[Defining Eventual and Strong Eventual Consistency]**  
[Timestamp: 6:35](https://youtu.be/9uCP3qHNbWw?t=395)  
> *"Strong eventual consistency guarantees convergence: replicas with the same updates reach the same state."*  
> *"Commutative operations allow replicas to apply updates in any order."*  

**Analysis**:  
- **Technical Explanation**: Eventual consistency (EC) ensures replicas eventually converge if updates stop. Strong eventual consistency (SEC) adds guarantees that replicas with identical update sets are consistent, even with reordered operations. This relies on commutative operations (e.g., incrementing counters) and conflict-free replicated data types (CRDTs).  
- **Context**: SEC addresses EC’s ambiguity by formalizing convergence, enabling practical implementations like collaborative editing or IoT device synchronization.  
- **Significance**: SEC’s mathematical rigor makes it viable for critical systems needing predictable outcomes.  
- **Connections**: Links to causal broadcast protocols, where message ordering respects causality but not total order.  

---

#### **[Conflict Resolution and Model Comparisons]**  
[Timestamp: 9:13](https://youtu.be/9uCP3qHNbWw?t=553)  
> *"Merge algorithms resolve conflicts by combining concurrent updates... into a clean final state."*  
> *"Eventual consistency operates under the weakest system assumptions—no quorums or synchrony."*  

**Analysis**:  
- **Technical Explanation**: Conflict resolution strategies range from simple (last-writer-wins) to complex (CRDT merges). The video compares consistency models by their assumptions: atomic commit (strongest, requires all nodes) vs. EC (weakest, no quorums).  
- **Context**: Systems like DynamoDB or Riak use vector clocks for conflict detection, while CRDTs enable automatic merging.  
- **Significance**: Emphasizes the scalability and fault tolerance of EC, contrasting it with linearizability’s constraints.  
- **Real-World Implications**: EC suits global-scale systems (social media feeds, DNS) where latency and partition tolerance are critical.  

---

### **Conclusion**  
The video progresses from linearizability’s limitations to eventual consistency’s pragmatic trade-offs. Key milestones include:  
1. **Trade-off Awareness**: Linearizability’s unavailability during partitions vs. EC’s "always-on" capability.  
2. **Conflict Management**: Real-world examples (calendar app) and theoretical frameworks (CAP theorem) illustrate the necessity of conflict resolution.  
3. **Model Hierarchy**: Atomic commit > consensus > linearizable get/set > EC, ordered by weakening assumptions.  

**Practical Importance**: EC enables scalable, fault-tolerant systems but requires careful handling of conflicts. Strong eventual consistency adds rigor, making it suitable for applications needing predictable convergence.  

**Learning Outcomes**: Viewers gain a nuanced understanding of consistency models, their trade-offs, and the role of system design in balancing availability, latency, and correctness. The video underscores that no "one-size-fits-all" model exists—choice depends on use-case requirements.