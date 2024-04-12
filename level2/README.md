# Level 2

## Goal
The goal of this task is to implement a concurrent queue system with multiple readers and writers, showcasing the usage of CSP, channels, and Go's concurrency model.

## Requirements
1. Implement a queue with a fixed capacity of `x` elements.
2. Multiple writer goroutines will generate and enqueue elements into the queue.
   - Each element should be a struct containing `{id: <writer-id>, val: <data>}`.
   - Writers should block if the queue is full and wait until space becomes available.
3. Multiple reader goroutines will dequeue elements from the queue and process them.
   - Readers should block if the queue is empty and wait until elements become available.
   - Readers should store the processed elements in a shared map `{id: <writer-id>, val: <data>}`.
4. Implement a batch writing mechanism for the shared map.
   - The batch writer goroutine should periodically write the contents of the shared map to a persistent storage (e.g., file or database).
   - The batch size and write interval should be configurable.
