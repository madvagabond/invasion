package world
import (
	"sync"
)


const (
	north = "north"
	south = "south"
	east = "east"
	west = "west"
)


const (
	map_change = "map_change"
	disconnect = "disconnect"
)


type TMSG struct {
	id int
}


type RMSG struct {
	status string
} 



type Connection struct {
	city string
	direction string
}




type Alien struct {
	id int
	ctr int
	
	location string
	ch chan RMSG
}




//Cities and Roads are implemented as a Vertex with an adjacency list
type City struct {
	name string
	connections []Connection
}




//The global state for the game
type WorldMap struct {
	cities *CityMap
	aliens *AlienMap
}


type Server struct {
	World *WorldMap
	ch chan TMSG
	Sig chan bool
} 


type CityMap struct {
	Map map[string]City
	Mu sync.RWMutex
}

type AlienMap struct {
	Map map[int]Alien
	Mu sync.RWMutex
}




func InitCityMap() *CityMap {
	return &CityMap{Map: make(map[string]City) }
}


func InitAlienMap() *AlienMap {
	return &AlienMap{Map: make(map[int]Alien) }
}

func (CM *CityMap) Get(key string) (City, bool) {
	CM.Mu.RLock()
	defer CM.Mu.RUnlock()
	
	v, e := CM.Map[key]
	return v, e
}

func (AM *AlienMap) Len() int {
	AM.Mu.RLock()
	defer AM.Mu.RUnlock()
	return len(AM.Map)
} 


func (CM *CityMap) Delete(key string) {
	CM.Mu.Lock()
	defer CM.Mu.Unlock()
	delete(CM.Map, key)

}

func (CM *CityMap) Put(key string, value City) {
	CM.Mu.Lock()
	defer CM.Mu.Unlock()
	CM.Map[key] = value
}





func (AM *AlienMap) Get(key int) (Alien, bool) {
	AM.Mu.RLock()
	defer AM.Mu.RUnlock()
	
	v, e := AM.Map[key]
	return v, e
}


func (AM *CityMap) Copy() map[string]City {
	AM.Mu.RLock() 
	defer AM.Mu.RUnlock()
	
	retVal := make(map[string]City) 
	
	for k, v := range AM.Map {
		retVal[k] = v
	}

	return retVal
}  



func (AM *AlienMap) Delete(key int) {
	AM.Mu.Lock()
	defer AM.Mu.Unlock()
	delete(AM.Map, key)

}



func (AM *AlienMap) Put(key int, value Alien) {
	AM.Mu.Lock()
	defer AM.Mu.Unlock()
	AM.Map[key] = value
}


func (AM *AlienMap) Copy() map[int]Alien {
	AM.Mu.RLock() 
	defer AM.Mu.RUnlock()
	
	retVal := make(map[int]Alien) 
	
	for k, v := range AM.Map {
		retVal[k] = v
	}

	return retVal
}  
