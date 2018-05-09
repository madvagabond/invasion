//maintains the internal state of the game
//note that with update functions instead of returning an error value when type does not exist it just returns the map unchanged, this might be slightly problematic. 
package state

import (
	"sync/atomic"
	"math/rand"
	"time"
	"sync"
)



const (
	north = "north"
	south = "south"
	east = "east"
	west = "west"
)





type Connection struct {
	cityName string
	direction string
}




type Alien struct {
	id int
	moveCtr int32
}




//Simple Node with an adjacency list, and the Direction as the label, and what invaders are there
type City struct {
	name string
	connections []Connection
	invaders []Alien
}





type CityMap struct {
	nodes map[string]City
	mu sync.RWMutex 
}



//adds edge to city
func (C *City) addConn(conn Connection) City {
	nconns := append(C.connections, conn)
	return City{C.name, nconns, C.invaders}	
}



//Only retains elements of list, where F(X) -> true 
func filterConns(conns []Connection, F func (Connection) bool ) []Connection {
	var retVal []Connection
	
	for _, x := range conns {

		if F(x) {
		retVal = append(retVal, x)

		} else {
			continue
		}
		
	}

	return retVal
}




func filterAliens(aliens []Alien, F func(Alien) bool) []Alien {
	var retVal []Alien
	
	for _, x := range aliens {
		if F(x) {
			retVal = append(retVal, x)

		} else {
			continue
		}  

	}

	return retVal
}  



func (C *City) rmConn(city_name string) City {
	fn := func(C Connection) bool {
		return C.cityName != city_name
	}

	
	nconns := filterConns(C.connections, fn)
	return City{C.name, nconns, C.invaders}	
} 





func (C *City) rmAlien(alien Alien) City {

	f := func (a Alien) bool {
		return alien.id != a.id 
	}

	ninvaders := filterAliens(C.invaders, f)
	return City{C.name, C.connections, ninvaders} 
	
}




func (C *City) addAlien(alien Alien) City {
	ninvaders := append(C.invaders, alien)
	return City{C.name, C.connections, ninvaders}	
} 



func (C *City) randomConn() Connection  {
	rand.Seed(time.Now().Unix())
	n := rand.Int() %  len(C.connections)
	return C.connections[n]
	
}




func (CM *CityMap) updateCity(cityName string, f func(City)City) *CityMap {

	CM.mu.RLock() 
	city := CM.nodes[cityName]
	CM.mu.RUnlock()
	

	city1 := f(city)

	CM.mu.Lock()
	CM.nodes[cityName] = city1
	CM.mu.Unlock()

	
	return CM 
	
}






func (CM *CityMap) destroyCity(cityName string) *CityMap {
	city := CM.nodes[cityName]

	conns := city.connections

	fn := func(city City) City {
		return city.rmConn(cityName) 
	}
	
	for _, x := range conns {
		CM.updateCity(x.cityName, fn)
	}

	delete(CM.nodes, cityName)
	return CM 
	
}



func (CM *CityMap) move(alien Alien, currentCity City) *CityMap {

	rmfn := func(city City) City {
		return city.rmAlien(alien)
	}

	CM.updateCity(currentCity.name, rmfn)

	atomic.AddInt32(&alien.moveCtr, 1)
	conn := currentCity.randomConn()

	city := CM.nodes[conn.cityName]

	fn := func(city City) City {
		return city.addAlien(alien)
	}

	CM.updateCity(city.name, fn)

	return CM
	
} 
