# dna-test

## Store
### Overview
Usually, for entities like `Account` I use relational databases 
like Postgres because those models perfectly fit in table 
structure and with the increase of entities, which may be 
connected with each other, foreign key constraints will be 
an ideal way to bind them together. However, as you 
recommended using in-memory storage in the task, I have 
tried to implement it to make it seem somehow reasonable. 
So, here in-memory storage is used as a data buffer to 
increase performance. Retrieved from the database data will 
be stored in cache, so the next request for the same data 
will be taken from the in-memory storage.

### Interfaces
Here I implemented `DbClient` and `CacheClient` interfaces 
with their corresponding functions, so they can be easily 
replaced with different realizations. In this project I 
used `GORM` and `Redis`, however, there is a small example 
of `DbClient` implemented with the help of a pure Postgres 
driver.

## Models
Here field `type` of the model `Account` is just a string. 
It is enough for creating, retrieving and making some 
operations depending on the value of that field. However, 
if `type` covers bigger functionality, it could be treated 
as a separate entity with its own model and table. In that 
case, it would be a foreign key for the `Accounts` table.

## Structure
Usually, in my work projects, I do not use `DbClient` 
interface with general sql commands. Even using `GORM`, 
we create a repo layer and writes all the sql queries for 
each unique database request and strictly parse data into 
the structure. For my individual projects, I always try to 
use something new, like store clients or error handler 
middleware here. By the way, I am aware that some practices 
that I might use could be not the best approach to solve 
certain problems, but I am always experimenting. However, 
I can say that I am very flexible in my code writing 
behaviour. I can easily adapt to the corporate coding 
style and keep all the conventions that were established 
at my work.