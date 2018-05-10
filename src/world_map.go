//maintains the internal state of the game
//note that with update functions instead of returning an error value when type does not exist it just returns the map unchanged, this might be slightly problematic. 
package invasion

import (
	"sync/atomic"
	"math/rand"
	"time"
	"fmt"
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
	state_change *Alien
} 



type Connection struct {
	cityName string
	direction string
}




type Alien struct {
	id int
	moveCtr int32
	
	location string
	ch chan RMSG
}







//Simple Node with an adjacency list, and the Direction as the label, and what invaders are there
type City struct {
	name string
	connections []Connection
}





type WorldMap struct {

	cities map[string]City
	aliens map[int]Alien
	
	mu sync.RWMutex 
}



//adds edge to city
func (C *City) AddConn(conn Connection) City {
	nconns := append(C.connections, conn)
	return City{C.name, nconns}	
}



//Only retains elements of list, where F(X) -> true 
func FilterConns(conns []Connection, F func (Connection) bool ) []Connection {
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





func (C *City) RmConn(city_name string) City {
	fn := func(C Connection) bool {
		return C.cityName != city_name
	}

	
	nconns := FilterConns(C.connections, fn)
	return City{C.name, nconns}	
} 






func (C *City) RandomConn() Connection  {
	rand.Seed(time.Now().Unix())
	n := rand.Int() %  len(C.connections)
	return C.connections[n]
	
}




func (WM *WorldMap) UpdateCity(cityName string, f func(City)City) *WorldMap {

	WM.mu.RLock() 
	city := WM.cities[cityName]
	WM.mu.RUnlock()
	

	city1 := f(city)

	WM.mu.Lock()
	WM.cities[cityName] = city1
	WM.mu.Unlock()

	
	return WM 
	
}






func (WM *WorldMap) DestroyCity(cityName string) *WorldMap {

	WM.mu.RLock() 
	city := WM.cities[cityName]
	WM.mu.RUnlock()

	
	conns := city.connections
	
	fn := func(city City) City {
		return city.RmConn(cityName) 
	}
	
	for _, x := range conns {
		WM.UpdateCity(x.cityName, fn)
	}

	WM.mu.Lock()
	delete(WM.cities, cityName)
	WM.mu.Unlock()
	
	return WM 
	
}



func FilterAliens(aliens map[int]Alien, F func(Alien) bool ) []Alien {

	var retVal []Alien
	
	for _, v := range aliens {

		if F(v) {
		retVal = append(retVal, v)

		} else {
			continue
		}
		
	}

	return retVal
}   



func (WM *WorldMap) RemoveAlien(alien Alien) {

	WM.mu.Lock()
	delete(WM.aliens, alien.id)
	WM.mu.Unlock()

	alien.ch <- RMSG{disconnect, nil} 
}




func (WM *WorldMap) Move (alien Alien) {
	
	ctr := alien.moveCtr
	atomic.AddInt32(&ctr, 1)

	if alien.moveCtr >= 1000 {

		WM.RemoveAlien(alien)
		fmt.Printf("Alien %d was removed due to making 1000 moves", alien.id) 

		return 
	}


	city_name := alien.location

	WM.mu.RLock() 
	current_city := WM.cities[city_name]
	WM.mu.RUnlock() 

	new_cname := current_city.RandomConn().cityName
	pred := func (a Alien) bool {return a.location == new_cname}

	rivals := FilterAliens(WM.aliens, pred)

	if len(rivals) > 1 {

		var ids []int 


		for _, x := range rivals {
			id := x.id
			ids = append(ids, id)
		} 

		
		for _, x := range rivals {
			WM.RemoveAlien(x)
		} 


		
		WM.DestroyCity(city_name)

		fmt.Printf("Aliens %v died in battle, and city %s was destroy", ids, city_name)

		return 
	}

	
}  
