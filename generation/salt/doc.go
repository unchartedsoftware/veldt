/*
Package salt manages provision of tiles to Veldt using Salt/Spark as a data source.

Because all requests have to be transmitted to the salt/spark data server, where they 
will have to be parsed independenly anyway, there is little use for parsing and 
validation on the GO side for salt-backed tile querries.

The go side of the salt tile system therefore just forwards messages over to salt as 
is as much as possible.

Salt works best with early notice as to the datasets it is using, so it can cache 
them as much as possible.  The salt tile server therefore separates out datasets 
and queries.  Queries note the dataset against which they are querying by ID, without
fine details.  They require that the dataset in question should first be registered.

Dataset registration must happen during a call to NewSaltTile.  Any number of datasets
can be configured in a single call.  Since this data is simply passed to the salt 
server, one could pass the information to set up dataset A in one call to NewSaltTile, 
and the information to dataset B in another, and use the tile constructors returned
indiscriminately, but obviously this would be somewhat confusing to the developer; one
is encouraged ot use datasets only in the constructor returned by the call to 
NewSaltTile in which they were initiated.

There are 3 points of communications between the GO and Scala side of the Salt tile 
system:

	salt.meta.Create (default_meta.go)
		Connects to Salt to request the metadata for a dataset

	NewSaltTile (salt_tile.go)
		Creates a constructor for salt tiles, and connects to Salt to register
		datasets [TODO: Perhaps separate out registration portion?]

	salt.Tile.Create (salt_tile.go)
		Connects to Salt to request a given tile or tiles

All use NewConnection in salt.go to make the actual connection.  NewConnection caches
its connection for reuse, so in typical use, should return the same connection every
time.

*/
package salt
