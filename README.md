# Kafkaesque
Attempt at "Build my own Kafka"

So learnings:
An important reason for this project is to learn Golang and in golang one of the 
things that make it good and unique is, Goroutines
So ill try to "get" them
so here we use Mutex, mutual exclusion like from OS, lock lets only one goroutine access

update: and whoa go is so cool like, we can use the defer keyword to schedule smth like a mutex.Unlock whenever the function exits, such clean, such safe

here the use is

Multiple producers may be publishing at the same time.

Multiple consumers may be subscribing or receiving at the same time.

Using sync.RWMutex ensures that:

Readers don’t block each other (fast broadcasts).

Writers block readers only when necessary (safe subscription changes).


So now a basic broker is working, consumers getting messages from topics

New example showcasing actual use case of broker:

```
[Producer-2] published: payment-1
[Producer-1] published: order-1
[Consumer-B] received from topic=orders : order-1
[Consumer-C] received from topic=payments : payment-1
[Consumer-A] received from topic=orders : order-1
[Producer-1] published: order-2
[Consumer-B] received from topic=orders : order-2
[Consumer-A] received from topic=orders : order-2
[Producer-2] published: payment-2
[Consumer-C] received from topic=payments : payment-2
[Producer-1] published: order-3
[Consumer-B] received from topic=orders : order-3
[Consumer-A] received from topic=orders : order-3
```


update 2:

what the helly golang??
` In Go, only identifiers starting with an uppercase letter are exported and accessible from other packages` ??? why
lolz