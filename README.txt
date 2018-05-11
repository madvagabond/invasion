
Code Organization 

The package invasion is the binary 
The package world is the library 

Most of the state change operations go on in world_map.go 
The data structures used are in structs.go
The goroutines and the actual simulation are in server.go 

The functions used for instantiation / parsing the maps 
Are in decoder.go 



My Critiques 
If you look through the commits you will probably notice the following things. 


It took me a while to shift from how I usually code, to something more idiomatic in Go
I tried to set it up in a goroutine per request manner, which indicates a tendency to PMO, even when it's not necessary. 
My preference for immutability and well typed code,  makes my go code, awkward, and clunky in some places, and leads to a lot of random variables, and I violated Dry for the Filter functions. 


Also it's not as efficient as I would have liked it to be. 

