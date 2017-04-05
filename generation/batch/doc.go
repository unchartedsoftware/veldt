/*
Package batch provides a set of tiles that batches tiling requests together, feeding
them to another tile generation package.

That other package must provide TileFactoryCtors instead of TileCtors, but the 
interface is nearly identical, and both interfaces can be provided by the same class.
The easiest way to do this is:

  * Have Parse(...) (from the Tile interface) merely store parameters for later -
    actual parameter inspection should happen during the CreateTiles call
  * Have Create(...) (from the Tile interface) call CreateTiles(...) (from the 
    TileFactory interface)

The former is a fairly serious limitation, and if done carelessly, may 
obfuscate some errors.  But carefully done, these two restrictions should cause
little problem.

Tiles are returned using one-use buffered channels.  Any user _must_ return 
something to all channels in a batched tile request.  Forgotten channels may 
use up the applications request buffer, limiting and eventually eliminating 
the ability of the application to request tiles at all, let alone in batches.

For an example of a generation package that can be used both for batch and single
tiles, see the salt generation package.
*/
package batch
