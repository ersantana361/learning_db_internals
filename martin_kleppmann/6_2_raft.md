---
title: 6.2 Raft
draft: false
tags:
  - system-design
  - distributed-systems
  - martin-kleppmann
  - intermediate
  - reference
  - consensus
  - algorithms
---
### **Introduction**  
- **Title**: [Distributed Systems 6.2: Raft](https://www.youtube.com/watch?v=uXEYuDwm7e4)  
- **Overview**:  
  This video provides a detailed walkthrough of the Raft consensus algorithm, focusing on its mechanisms for achieving total order broadcast in distributed systems. The core objectives include explaining Raft’s leader election process, log replication, and commitment protocols. The structural flow progresses from node states (follower, candidate, leader) to the intricacies of term-based voting, log consistency checks, and fault tolerance. Key themes include the role of terms in ensuring leadership legitimacy, the importance of log synchronization for consistency, and the use of quorums to commit entries. The video emphasizes Raft’s design for clarity and its practical applications in real-world distributed systems.  

---

### **Chronological Analysis**  

#### **[Introduction to Raft’s State Machine and Leader Election]**  
[Timestamp: 0:01](https://youtu.be/uXEYuDwm7e4?t=1)  
> *"Every node can be in one of three states: follower, candidate, or leader... A node starts as a follower and transitions to a candidate if it detects no leader heartbeat."*  
> *"A higher term number always takes precedence... Candidates step down upon discovering a higher term."*  

**Analysis**:  
- **Technical Explanation**: Raft’s state machine ensures only one leader exists per term. Nodes begin as followers; if they time out waiting for a leader’s heartbeat, they become candidates, increment their term, and request votes. A majority (quorum) of votes is required to become a leader.  
- **Contextualization**: This segment establishes Raft’s foundational mechanics, highlighting how terms enforce leadership transitions and prevent split-brain scenarios.  
- **Significance**: The term-based system ensures liveness and safety by prioritizing newer terms, which is critical during network partitions.  
- **Real-World Implications**: Randomized election timeouts reduce contention, a design choice used in systems like Kubernetes and etcd.  

---

#### **[Voting Mechanism and Log Consistency Checks]**  
[Timestamp: 2:15](https://youtu.be/uXEYuDwm7e4?t=135)  
> *"Followers grant votes only if the candidate’s log is at least as up-to-date as theirs... Logs are compared via term numbers and lengths."*  
> *"A candidate’s last log term must be ≥ the follower’s, or their log must be longer if terms match."*  

**Analysis**:  
- **Technical Explanation**: Vote requests include the candidate’s last log term and index. Followers deny votes if their own log is more recent, preventing outdated leaders.  
- **Contextualization**: This ensures leaders possess the most complete log, preserving system consistency.  
- **Significance**: The "up-to-date" log rule prevents data loss by favoring candidates with newer or equally complete logs.  
- **Connections**: Links to leader legitimacy—only nodes with authoritative logs can propagate entries.  

---

#### **[Log Replication and Commitment Protocol]**  
[Timestamp: 7:28](https://youtu.be/uXEYuDwm7e4?t=448)  
> *"The leader appends new entries to its log, then replicates them to followers... Entries are committed once a quorum acknowledges them."*  
> *"Followers truncate conflicting log entries and adopt the leader’s suffix."*  

**Analysis**:  
- **Technical Explanation**: Leaders broadcast log entries to followers. If a follower’s log diverges, it truncates entries after the last agreed index and appends the leader’s entries.  
- **Contextualization**: This ensures all nodes converge on the same log state, even after failures.  
- **Significance**: The append-only log with term identifiers guarantees eventual consistency.  
- **Real-World Applications**: Used in databases like CockroachDB for replication.  

---

#### **[Handling Failures and Log Truncation]**  
[Timestamp: 18:00](https://youtu.be/uXEYuDwm7e4?t=1080)  
> *"If a follower’s log conflicts with the leader’s, it discards entries beyond the last agreed index... Leaders retry sending entries until logs match."*  

**Analysis**:  
- **Technical Explanation**: The `AppendEntries` RPC includes the leader’s term, previous log index/term, and new entries. Followers validate consistency before appending.  
- **Contextualization**: This "rollback and repair" mechanism handles network partitions or stale leaders.  
- **Significance**: Ensures logs are identical across nodes for committed entries, critical for linearizability.  

---

#### **[Commitment and Total Order Delivery]**  
[Timestamp: 23:33](https://youtu.be/uXEYuDwm7e4?t=1413)  
> *"A log entry is committed once a majority of nodes acknowledge it... The leader then delivers committed entries to the application."*  

**Analysis**:  
- **Technical Explanation**: The leader tracks acknowledgments via `sentLength` and `ackedLength`. Entries are committed when a quorum confirms receipt, after which they are applied to state machines.  
- **Contextualization**: Commitment finalizes the order of operations, enabling total order broadcast.  
- **Real-World Implications**: This ensures all nodes process requests in the same order, vital for systems like financial ledgers.  

---

### **Conclusion**  
The video methodically progresses from Raft’s leader election to log commitment, emphasizing its reliability and simplicity. Key milestones include term-based leadership transitions, log synchronization via append-only entries, and quorum-driven commitment. The algorithm’s practical importance lies in its fault tolerance and clarity, making it a cornerstone of modern distributed systems like distributed databases and coordination services. By ensuring logs are consistent and operations are totally ordered, Raft provides a robust framework for building systems that prioritize both availability and consistency. The learning outcome is a deep understanding of how consensus algorithms underpin reliable distributed coordination.