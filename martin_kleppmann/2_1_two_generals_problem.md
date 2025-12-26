---
title: 2.1 The Two generals problem
draft: false
tags:
  - system-design
  - distributed-systems
  - martin-kleppmann
  - intermediate
  - reference
  - fault-tolerance
  - reliability
---
### **Introduction**  
- **Title**: [Distributed Systems 2.1: The two generals problem](https://www.youtube.com/watch?v=MDuWnzVnfpI)  
- **Overview**:  
  The video explores foundational challenges in distributed systems through the lens of the *two generals problem*, a classic thought experiment. It begins by defining system models and their role in formalizing assumptions about failures (e.g., message loss, node crashes). Using the two generals analogy, the presenter illustrates the impossibility of achieving guaranteed consensus in unreliable networks. The analysis transitions to real-world applications, such as e-commerce payment systems, to demonstrate how theoretical limitations are mitigated in practice. Key themes include uncertainty in communication, the concept of *common knowledge*, and pragmatic solutions like revocable actions.

---

### **Chronological Analysis**  

#### **[Introduction to System Models and the Two Generals Problem]**  
[Timestamp: 0:00–1:59](https://youtu.be/MDuWnzVnfpI?t=0)  
> *"A system model is very important because... we have to be precise about what failures we are assuming are possible."*  
> *"The two generals can't just talk to each other... they can only communicate via messengers [that] might get captured."*  

**Analysis**:  
The video establishes the necessity of system models to formalize assumptions in distributed systems. The two generals problem is introduced as a metaphor for nodes (generals) needing to coordinate an action (attack) over an unreliable channel (messengers). The critical constraint is that both must act simultaneously to succeed, but messages may be lost, creating uncertainty. This segment contextualizes distributed systems as inherently probabilistic, where nodes lack shared state certainty. The problem’s significance lies in its illustration of *consensus impossibility* under unreliable communication—a cornerstone of distributed systems theory.  

---

#### **[Analyzing Communication Strategies and Infinite Regress]**  
[Timestamp: 2:00–7:42](https://youtu.be/MDuWnzVnfpI?t=120)  
> *"General one does not know whether... the initial message didn’t get through or the response was lost."*  
> *"You end up with infinite chains of... messages before there’s any certainty."*  

**Analysis**:  
The presenter dissects potential solutions, such as retrying messages or requiring acknowledgments, but demonstrates their futility. If General 1 attacks unconditionally, they risk acting alone. If they wait for confirmation, General 2 faces the same dilemma, leading to infinite recursive dependencies. This reflects the concept of *common knowledge* in distributed systems: even with multiple acknowledgments, nodes cannot achieve mutual certainty of intent. The segment underscores the theoretical impossibility of perfect consensus in asynchronous networks, a result formalized by the FLP impossibility theorem. The infinite regress problem highlights the trade-off between safety (avoiding incorrect actions) and liveness (progressing despite failures).  

---

#### **[Application to Online Shop and Payment Service]**  
[Timestamp: 7:43–10:31](https://youtu.be/MDuWnzVnfpI?t=463)  
> *"The online shop dispatches goods if and only if the payment service charges the card... analogous to the two generals problem."*  
> *"Messages might get lost... it’s impossible to achieve certainty."*  

**Analysis**:  
The video bridges theory to practice using an e-commerce example. Here, the shop (General 1) and payment service (General 2) must ensure atomicity: goods ship iff payment succeeds. However, network failures mirror messenger loss, risking inconsistencies (e.g., charging without shipping). The presenter emphasizes that real-world systems avoid deadlock by relaxing strict atomicity. For instance, payments are made *revocable* (via refunds), allowing post-facto corrections. This mirrors *compensating transactions* in distributed databases. The example shows how theoretical impossibility is circumvented through pragmatic, fault-tolerant design.  

---

#### **[Real-World Mitigations and Idempotent Operations]**  
[Timestamp: 10:32–11:31](https://youtu.be/MDuWnzVnfpI?t=632)  
> *"The payment service will... charge the card because charges can be refunded."*  
> *"Check if the payment succeeded... using idempotent requests."*  

**Analysis**:  
The video concludes with practical strategies to handle uncertainty. Revocable actions (e.g., refunds) and idempotent operations (retrying safely) allow systems to progress despite message loss. By designing operations to be reversible or repeatable, real-world applications sidestep the need for absolute consensus. For example, idempotent payment requests ensure duplicate messages don’t overcharge users. This segment ties back to system models: by assuming eventual consistency and building recovery mechanisms, engineers mitigate the two generals’ theoretical constraints. The emphasis shifts from impossibility to probabilistic reliability, reflecting real-world distributed systems like payment gateways and cloud services.  

---

### **Conclusion**  
The video progresses from abstract theory (two generals problem) to concrete applications (payment systems), emphasizing the tension between ideal guarantees and practical feasibility. Key milestones include:  
1. **System Models**: Framing assumptions about failures is critical for algorithm design.  
2. **Consensus Impossibility**: Unreliable communication inherently limits certainty.  
3. **Pragmatic Adaptations**: Real-world systems use revocable actions and idempotency to achieve *eventual consistency*.  

The theoretical exploration underscores the importance of understanding distributed systems’ limitations, while the practical examples demonstrate how engineers innovate within these bounds. The learning outcome is a nuanced appreciation of trade-offs: while perfect consensus is unattainable, probabilistic reliability and fault tolerance enable functional, large-scale systems.