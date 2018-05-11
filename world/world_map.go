package world
import (
	"math/rand"
	"time"
	"fmt"
) 

func (WM *WorldMap) UpdateCity(cityName string, f func(City)City) *WorldMap {
	city, b := WM.cities.Get(cityName)
	
	if b != true {
		return WM
	} 
	
	city1 := f(city)
	WM.cities.Put(cityName, city1) 
	return WM 
	
}


//Filters the list and returns the elements that meet the predicate
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


//Transforms Aliens map to a list, then returns the values that meet the Predicate 
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






func (C *City) AddConn(conn Connection) City {
	nconns := append(C.connections, conn)
	return City{C.name, nconns}	
}


func (C *City) RmConn(city_name string) City {
	fn := func(C Connection) bool {
		return C.city != city_name
	}

	
	nconns := FilterConns(C.connections, fn)
	return City{C.name, nconns}	
} 



//The basis for the random walk ish algorithm
func (C *City) RandomConn() *Connection  {

	
	rand.Seed(time.Now().Unix())
	
	if len(C.connections) != 0 { 
		n := rand.Intn(len(C.connections))
		return &C.connections[n]
	} else {return nil}

	
}




func (WM *WorldMap) RemoveAlien(alien Alien) {

	WM.aliens.Delete(alien.id)
	alien.ch <- RMSG{disconnect}
}


//Blow up the city with your over 9000 power level
func (WM *WorldMap) DestroyCity(cityName string) *WorldMap {

	
	city, e := WM.cities.Get(cityName)

	
	if e != true {
		return WM
	} 

	
	conns := city.connections
	
	fn := func(city City) City {
		return city.RmConn(cityName) 
	}

	
	WM.cities.Delete(cityName)

	
	for _, x := range conns {
		WM.UpdateCity(x.city, fn)
	}

	
	return WM 
	
}



//Forces the Invaders in the same city to brutally murder eachother
func (WM *WorldMap) RemoveRivals(city_name string) bool {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			return
		}
    }()
	
	
	pred := func (a Alien) bool {return a.location == city_name}

	
	rivals := FilterAliens(WM.aliens.Copy(), pred)

	
	
	b := len(rivals) >= 2

	if b != true {
		return false 
	} 
	
	var ids []int 


	for _, x := range rivals {
		id := x.id
		ids = append(ids, id)
	} 

	
	for _, x := range rivals {
		WM.RemoveAlien(x)
	}

	
	WM.DestroyCity(city_name)

	fmt.Printf("Aliens %v died in the battle of %s and the city was destroyed \n", ids, city_name)
	return true
}  

