# A pretty bad document database for go projects

Uses badgerdb, has very limited query capabilities and bad performance

I wouldn't use for anything super important

But it's nice to have when doing quick prototypes

Requires gcc. So for windows you have to install https://jmeubank.github.io/tdm-gcc/ and use MinGW (installed from tdm-gcc)

## How to use

Uses json marshalling so when defining models add `json` struct tags example `type Person struct { Name string `json:"name"` }`.

Open database with func Open(), use initialized database to select a collection with Collection method. Use Collection methods to query and mutate database. List all collections with db.ListCollections(). List all document keys with collection.GetAllKeys()

Structs returned from database don't have the key embedded. Sorry
#   b l a b d b 
 
 
