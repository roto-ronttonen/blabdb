# A pretty bad document database for go projects

Uses badgerdb, has very limited query capabilities and bad performance

I wouldn't use for anything important

But it's nice to have when doing quick prototypes

Requires gcc. So for windows you have to install https://jmeubank.github.io/tdm-gcc/ and use MinGW (installed from tdm-gcc)

## How to use

Uses json marshalling so when defining models add `json` struct tags example `type Person struct { Name string `json:"name"` }`.

Open database with func Open(), use initialized database to select a collection with Collection method. Use Collection methods to query and mutate database. List all collections with db.ListCollections(). List all document keys with collection.GetAllKeys()

Key is embedded to returned documents as `key`. If you want your struct to have access it use `Key string `json:"key"``. This also means that if you are storing a value name key it will be overwritten

Find uses writeTo object for having a struct reference for unmarshalling. writeTo will be equal to last value found by query. Result is array of query results

### Query

Inside Where is and. Eeach Where block is an or

Limitations:

- Only == operator supported
- Cant traverse nested structures

If no limit is provided its set at 10. Limit is capped at 100 because if someone forgets to add a limit to their query their computer might explode or something

### Benchmarks

I haven't done any benchmarks but im sure they suck
