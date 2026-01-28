# payments-core
payments-core is the authoritative core ledger service responsible for maintaining account balances and executing atomic transfers.


# Api Flows 
## Transfers API - POST {url}/transactions 
### Functional Requirements
* Idempotency Guarantee 
* Atomicity Guarantee 

### Non-Functional Requirements 
I want to ensure my system has CA (of CAP Theorem)
* highly available
* highly consistent . We dont want the system to have inconsistencies in the data we show to users . **Eventual consistency** is out of question. 
> NOTE : Since this is a sample application I did not emphasize on the cutting edge design decisions , For example I'm taking row lock. this is not an ideal solution in a distributed system . I'd prefer 'redis' for ensuring we guarantee atomicity. 

![transfers api plant UML](design/transfers_api.png)


## Transfers API - POST {url}/accounts 

![transfers api plant UML](design/create_account.png)
