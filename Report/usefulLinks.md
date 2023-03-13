## Useful Links
* [CRDT Tech Glossary - CRDT Keywords Definitions](https://crdt.tech/glossary)
* [Operation Based CRDT Framework (Reliable Causal Broadcast?)](http://krasserm.github.io/2016/10/19/operation-based-crdt-framework/)


## Updates
- ✔ Add references to PALELC Theorem

## To Do
- Watch Martin Kleppman video on CRDTs
- Try detecting connection dropped 

## Questions 
- Should UML class diagrams be within the design section? Yes
- Should things like language go in design or implementation? Design
- What should implementation show? (e.g. code snippets, UML diagrams, etc.) Snippets / Problems and explanations


If you want high availability and performance, causal consistency may be a good choice for your system. Causal consistency is an intermediate consistency model that can be achieved together with high availability and performance requirements even in presence of network partitions12. It means that your system can tolerate failures and delays without sacrificing data integrity.

However, causal consistency also has some challenges and trade-offs. For example, you need to track the causal dependencies between operations using some mechanism such as timestamps or vector clocks13. This may introduce some overhead and complexity to your system. You also need to decide how to handle concurrent updates that are not causally related, such as using last-writer-wins or merge functions34.

Therefore, you should carefully evaluate your system requirements and design before choosing causal consistency as your consistency model.

Does this help? 😊

Sure, I can give you an example of causal consistency. Suppose you have a distributed system that allows users to post and reply to comments on a social media platform. Causal consistency guarantees that if a user sees a reply to a comment, they will also see the original comment that was replied to1. This is because the reply is causally dependent on the comment, and they should be observed in the same order by every user.

However, causal consistency does not guarantee that users will see all comments or replies in the same order. For example, two users may post concurrent comments that are not causally related to each other2. Different users may see these comments in different orders depending on which server they connect to. This is acceptable under causal consistency as long as there is no causal dependency between them.

Does this make sense? 😊

To implement causal consistency with eventual delivery, you need to use some mechanism to track and enforce the causal dependencies between operations. One possible way is to use multi-part timestamps, such as Lamport clocks or vector clocks1. These timestamps are assigned to each object and updated whenever an operation is performed on it. They can be used to determine the causal order of operations and ensure that they are executed in the same order by every replica.

However, using timestamps alone may not be enough to achieve causal consistency. You also need to ensure that every replica eventually receives and applies all operations that are causally related to its current state2. This can be done by using some form of gossip protocol or anti-entropy mechanism3 that periodically exchanges updates between replicas. This way, you can guarantee that your system will eventually converge to a consistent state that reflects all causal dependencies.

Does this answer your question? 😊

Okay, let me try to give you an example of causal consistency with eventual delivery. Suppose you have a distributed system that allows users to post and like tweets on a microblogging platform. Each tweet has a unique ID and a counter that tracks how many likes it has received. Each user has a replica of the system that stores their own tweets and likes, as well as those of other users they follow.

Now suppose Alice posts a tweet with ID 1 and Bob likes it. Alice’s replica updates the counter of tweet 1 to 1 and sends this update to Bob’s replica. Bob’s replica receives this update and applies it to its own copy of tweet 1. This is an example of causal consistency, because Bob sees Alice’s tweet before liking it, and Alice sees Bob’s like after posting her tweet.

However, suppose Charlie also likes Alice’s tweet, but he does not receive Alice’s update yet. Charlie’s replica updates the counter of tweet 1 to 1 and sends this update to Alice’s replica. Alice’s replica receives this update and applies it to its own copy of tweet 1. This is an example of eventual delivery, because Charlie’s update is eventually delivered to Alice’s replica.

However, there is a problem here: Alice’s replica now has an incorrect counter value for tweet 1. It should be 2, not 1. This is because Charlie’s update was concurrent with Bob’s update, but they did not see each other’s updates before applying their own1. To fix this problem, you need some form of conflict resolution2 that can merge concurrent updates and produce a consistent result.

One possible way to do this is to use vector clocks3. A vector clock is a data structure that records the logical time of each operation on each replica. It can be used to detect concurrent updates and resolve them by applying some predefined rule, such as taking the maximum value or adding them together.

Using vector clocks, Alice’s replica can detect that Charlie’s update was concurrent with Bob’s update by comparing their vector clocks. It can then resolve the conflict by taking the maximum value between them: max(1, 1) = 1. It can then add this value to its own counter value: 1 + 1 = 2. This way, Alice’s replica can obtain the correct counter value for tweet 1.

Does this make sense? 😊