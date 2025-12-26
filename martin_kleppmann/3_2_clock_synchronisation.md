---
title: 3.2 Clock Synchronisation
draft: false
tags:
  - system-design
  - distributed-systems
  - martin-kleppmann
  - intermediate
  - reference
  - time
  - synchronization
---
### **Introduction**  
- **Title**: [Distributed Systems 3.2: Clock Synchronisation](https://www.youtube.com/watch?v=mAyW-4LeXZo)  
- **Overview**:  
  This lecture addresses the challenges of synchronizing clocks in distributed systems, focusing on the Network Time Protocol (NTP). It begins by explaining clock skew due to quartz clock drift, introduces NTP’s hierarchical structure (strata), and details its synchronization mechanisms. The video transitions into practical implications, including slewing/stepping clocks and the critical distinction between monotonic and time-of-day clocks for software reliability. Key themes include balancing accuracy with practical constraints and mitigating synchronization pitfalls in real-world systems.

---

### **Chronological Analysis**  

#### **[Introduction to Clock Skew and NTP]**  
[Timestamp: 0:02](https://youtu.be/mAyW-4LeXZo?t=2)  
> *"every computer pretty much contains a quartz clock... we have to somehow try and make the clocks reasonably accurate."*  
> *"clock skew is... the difference between those two [clocks] at the same instant in time."*  

- **Technical Explanation**: Quartz clocks drift due to temperature sensitivity and manufacturing variances, causing **clock skew**—the difference between two clocks at the same moment. NTP mitigates this by synchronizing with authoritative time sources (atomic clocks/GPS).  
- **Context**: Establishes the problem of unreliable physical clocks in distributed systems and introduces NTP as a solution.  
- **Significance**: Highlights the foundational challenge of time coordination in asynchronous networks, where perfect synchronization is impossible.  
- **Real-World Application**: NTP’s ubiquity in operating systems (e.g., macOS settings) ensures global timekeeping consistency for logs, security certificates, and distributed transactions.  
- **Connection**: Sets the stage for discussing NTP’s hierarchical design and synchronization techniques.  

---

#### **[NTP Stratum Hierarchy and Error Mitigation]**  
[Timestamp: 1:04](https://youtu.be/mAyW-4LeXZo?t=64)  
> *"servers are arranged into... strata... stratum 1 is connected directly to a stratum 0 time source."*  
> *"statistical techniques... query multiple servers... exclude outliers."*  

- **Technical Explanation**: NTP uses a **stratum hierarchy** (0-15) to propagate time. Stratum 0 devices (atomic clocks/GPS) feed stratum 1 servers, which cascade to lower strata. Redundancy and outlier rejection improve reliability.  
- **Context**: Explains how NTP balances accuracy and scalability while minimizing single points of failure.  
- **Significance**: Stratum levels trade precision for accessibility, enabling global synchronization without requiring every device to connect directly to atomic clocks.  
- **Real-World Application**: Enterprises use stratum 2/3 servers for internal synchronization, reducing load on primary time sources.  
- **Connection**: Introduces statistical methods critical for NTP’s skew estimation, covered next.  

---

#### **[NTP’s Clock Skew Estimation Mechanism]**  
[Timestamp: 4:00](https://youtu.be/mAyW-4LeXZo?t=240)  
> *"estimate the clock skew... using timestamps t1-t4... assume network latency is symmetric."*  

- **Technical Explanation**: NTP calculates skew using four timestamps: client send (t1), server receive (t2), server reply (t3), and client receive (t4). Total delay (Δ = (t4 - t1) - (t3 - t2)) assumes symmetric network latency (Δ/2 for each direction). Skew θ = [(t2 - t1) + (t3 - t4)] / 2.  
- **Context**: Demonstrates how NTP compensates for variable network delays without synchronized clocks.  
- **Significance**: The symmetric latency assumption simplifies calculations but introduces error in asymmetric networks (e.g., congested links).  
- **Real-World Application**: NTP clients use iterative sampling to refine estimates, achieving sub-10ms accuracy in stable networks.  
- **Connection**: Informs subsequent discussion on slewing vs. stepping to correct skew.  

---

#### **[Slewing vs. Stepping Clocks]**  
[Timestamp: 7:58](https://youtu.be/mAyW-4LeXZo?t=478)  
> *"slewing the clock... slightly speed up or slow down... stepping forcibly adjusts the clock."*  

- **Technical Explanation**:  
  - **Slewing**: Gradually adjusts clock rate (≤500 ppm) for small skews (<125 ms), avoiding abrupt jumps.  
  - **Stepping**: Instantly corrects large skews (>125 ms) but risks software errors (e.g., negative durations).  
- **Context**: Balances precision with stability—slewing maintains continuity, while stepping handles severe drift.  
- **Significance**: Slewing preserves application integrity (e.g., timers), whereas stepping is a last resort.  
- **Real-World Application**: Financial systems prioritize slewing to avoid timestamp anomalies in high-frequency trades.  
- **Connection**: Leads into software design considerations for handling clock adjustments.  

---

#### **[Software Implications: Monotonic vs. Time-of-Day Clocks]**  
[Timestamp: 9:32](https://youtu.be/mAyW-4LeXZo?t=572)  
> *"monotonic clock... moves forward at a near constant rate... unaffected by NTP stepping."*  

- **Technical Explanation**:  
  - **Monotonic clocks** (e.g., Java’s `nanoTime()`) measure intervals unaffected by NTP adjustments.  
  - **Time-of-day clocks** (e.g., `currentTimeMillis()`) provide wall-clock time but risk jumps during stepping.  
- **Context**: Guides developers to use monotonic clocks for elapsed time measurements and time-of-day for cross-system coordination.  
- **Significance**: Prevents bugs like negative durations or performance metric inaccuracies.  
- **Real-World Application**: Databases use monotonic clocks for query timeouts; TLS certificates rely on synchronized time-of-day clocks.  
- **Connection**: Reinforces the lecture’s theme: clock synchronization is both a network and software challenge.  

---

### **Conclusion**  
The video progresses from the problem of quartz clock drift to NTP’s hierarchical and statistical solutions, culminating in software best practices:  
1. **Key Milestones**:  
   - Clock skew as an inherent challenge in distributed systems.  
   - NTP’s stratum hierarchy and timestamp-based skew estimation.  
   - Slewing/stepping trade-offs and the critical role of monotonic clocks.  
2. **Practical Importance**:  
   - NTP enables global time coordination but requires careful handling of network asymmetry and software clocks.  
   - Monotonic clocks ensure reliable performance metrics, while time-of-day clocks enable cross-system event ordering.  
3. **Learning Outcomes**:  
   - Understand NTP’s mechanisms and limitations.  
   - Apply monotonic clocks for interval measurements to avoid NTP-induced errors.  
   - Design systems tolerant of clock skew and synchronization discontinuities.  

This lecture underscores that while perfect synchronization is unattainable, strategic protocols and coding practices mitigate risks, ensuring robust distributed systems.