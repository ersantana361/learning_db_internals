---
title: 7.1 Two-phase commit
draft: false
tags:
  - system-design
  - distributed-systems
  - martin-kleppmann
  - intermediate
  - reference
  - transactions
  - commit-protocols
---
### **Introduction**  
- **Title**: [Distributed Systems 7.1: Two-phase commit](https://www.youtube.com/watch?v=-_rdWB9hN1c)  
- **Overview**:  
  This video explores the Two-Phase Commit (2PC) protocol, a foundational algorithm for ensuring atomicity in distributed transactions. The core objectives include clarifying the atomic commitment problem, explaining 2PC’s mechanics, and addressing its limitations, particularly coordinator failures. The video contrasts 2PC with consensus algorithms, introduces a fault-tolerant variant using total order broadcast, and emphasizes the interplay between atomicity and consistency in distributed systems. Key themes include the role of coordinators, the trade-offs between availability and safety, and the integration of consensus mechanisms to enhance fault tolerance.  

---

### **Chronological Analysis**  

#### **[Clarifying Consistency in Distributed Systems]**  
[Timestamp: 0:01](https://youtu.be/-_rdWB9hN1c?t=1)  
> *"Consistency means many things... In ACID, it refers to database state invariants; in replication, it means replicas agreeing on updates."*  
> *"Atomic commit ensures transactions either fully commit or abort across all nodes."*  

**Analysis**:  
- **Technical Explanation**: The video distinguishes between ACID consistency (preserving database invariants) and distributed consistency (replica synchronization). Atomic commit ensures all nodes agree on transaction outcomes.  
- **Contextualization**: Establishes why atomicity is critical—partial updates violate invariants, necessitating protocols like 2PC.  
- **Significance**: Highlights the necessity of atomicity as the foundation for higher-level consistency models.  
- **Real-World Implications**: Used in distributed databases (e.g., Apache Kafka transactions) to maintain data integrity.  
- **Connections**: Links to prior lectures on ACID properties and consensus algorithms like Raft.  

---

#### **[Atomic Commitment Problem and 2PC Basics]**  
[Timestamp: 3:44](https://youtu.be/-_rdWB9hN1c?t=224)  
> *"Atomic commit requires all nodes to agree: commit if all vote ‘yes,’ abort if any vote ‘no.’"*  
> *"2PC phases: prepare (voting) and commit/abort (decision)."*  

**Analysis**:  
- **Technical Explanation**: 2PC ensures atomicity via a coordinator managing two phases:  
  1. **Prepare**: Nodes vote on readiness to commit.  
  2. **Commit/Abort**: Coordinator finalizes based on unanimous votes.  
- **Contextualization**: Unlike consensus, 2PC requires unanimity and cannot proceed with partial node failures.  
- **Significance**: Guarantees no partial commits, critical for cross-node transactions (e.g., financial systems).  
- **Real-World Implications**: Traditional databases use 2PC but face blocking if nodes fail.  

---

#### **[Two-Phase Commit Protocol Mechanics]**  
[Timestamp: 7:28](https://youtu.be/-_rdWB9hN1c?t=448)  
> *"Coordinator writes decision to disk; replicas lock resources during prepare."*  
> *"A ‘yes’ vote binds replicas to commit if instructed."*  

**Analysis**:  
- **Technical Explanation**: During **prepare**, nodes persist changes and lock resources. A "yes" vote is irrevocable, ensuring they can commit later.  
- **Contextualization**: The coordinator’s disk write ensures recovery after crashes, preserving decision durability.  
- **Significance**: Locks prevent conflicting updates but risk blocking systems during failures.  
- **Real-World Applications**: Used in distributed SQL databases (e.g., Google Spanner) for cross-shard transactions.  

---

#### **[Coordinator Failure and Blocking Issue]**  
[Timestamp: 10:30](https://youtu.be/-_rdWB9hN1c?t=630)  
> *"If the coordinator crashes after prepare, replicas are blocked until recovery."*  
> *"Nodes cannot unilaterally abort—violates atomicity."*  

**Analysis**:  
- **Technical Explanation**: Coordinator crashes leave nodes in a "prepared" state, holding locks indefinitely. This blocks progress until recovery.  
- **Contextualization**: Highlights 2PC’s lack of fault tolerance compared to quorum-based consensus.  
- **Significance**: Demonstrates the trade-off: atomicity requires sacrificing availability during partitions.  
- **Real-World Implications**: Systems like XA transactions face downtime if coordinators fail.  

---

#### **[Fault-Tolerant 2PC with Total Order Broadcast]**  
[Timestamp: 18:00](https://youtu.be/-_rdWB9hN1c?t=1080)  
> *"Total order broadcast disseminates votes; failure detectors time out unresponsive nodes."*  
> *"First vote per replica decides outcome; replicas self-resolve via consensus."*  

**Analysis**:  
- **Technical Explanation**: Integrates total order broadcast (e.g., Raft) to propagate votes. Nodes use failure detectors to infer crashes and vote "abort" on behalf of unresponsive peers.  
- **Contextualization**: Eliminates coordinator dependency, allowing nodes to autonomously reach consistent decisions.  
- **Significance**: Enhances availability by leveraging consensus for fault tolerance.  
- **Connections**: Builds on earlier lessons about total order broadcast and quorum systems.  
- **Real-World Implications**: Modern systems (e.g., CockroachDB) use hybrid models for distributed transactions.  

---

### **Conclusion**  
The video progresses from defining consistency challenges to solving atomic commitment via 2PC and its fault-tolerant variant. Key milestones include:  
1. **Atomic Commitment Problem**: Necessitates unanimous agreement for cross-node transactions.  
2. **2PC Mechanics**: Coordinator-led prepare/commit phases ensure atomicity but risk blocking.  
3. **Fault-Tolerant 2PC**: Integrates consensus (total order broadcast) to resolve coordinator bottlenecks.  

**Practical Importance**: 2PC remains foundational for distributed transactions but requires enhancements (e.g., consensus) for real-world resilience. **Theoretical Insight**: Balancing atomicity and availability demands trade-offs, addressed through hybrid models. **Learning Outcome**: Understanding 2PC’s role in distributed systems clarifies how modern databases achieve ACID guarantees while mitigating coordinator dependence.