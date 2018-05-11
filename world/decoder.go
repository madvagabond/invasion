package world

import (
	"errors"
	"strings"
	"math/rand"
	"io/ioutil"
	"time"

)



func Opposite(direction string) (string, error)  {
	direct := strings.ToLower(direction)
	switch direct {

	case north:
		return south, nil

	case south:
		return north, nil

	case east:
		return west, nil

	case west:
		return east, nil

	default:
		return "", errors.New("Not A Direction")
		
	}


}





func (C *City) Contains(f func(Connection) bool ) bool {
	var b bool
	b = false

	conns := C.connections
	
	for _, x := range conns {

		if f(x) == true {
			b = true 
			break
		} else {
			continue
		}  


	}

	return b 

}




	//Isn't thread safe due to only being run at time of construction
func (WM *WorldMap) Validate() *WorldMap {


	check := func(L string, R string) bool {
		City, _ := WM.cities.Get(R)

		pred := func(C Connection) bool {
			return C.city == L
		}

		return City.Contains(pred)		
	}




	
	symmetry := func(C City) {

		
		conns := C.connections

		for _, x := range conns {
			b := check(C.name, x.city)

			direction, e := Opposite(x.direction)
			Neighbor, _ := WM.cities.Get(x.city)
		
			
			if e != nil {
				WM.cities.Map[C.name] = C.RmConn(x.city)
				continue 
			}


			if b == true {continue} else {
				nconn := Connection{city: C.name, direction: direction}
				WM.cities.Put(x.city, Neighbor.AddConn(nconn))
			}  

		}

		
	}  
	
	
	for _, v := range WM.cities.Copy() {
		symmetry(v) 
	}

	return WM
	
}



func ParseCity(text string) City  {
	tokens := strings.Split(text, " ")

	parseConnection := func (s string) Connection {
		items := strings.Split(s, "=")
		return Connection{direction: items[0], city: items[1]}
	}  
	
	name := tokens[0]
	rest := tokens[1:]

	var conns []Connection
	
	for _, x := range rest {
		conn := parseConnection(x)
		conns = append(conns, conn)
	}


	return City{name: name, connections: conns}
}




func ParseCities(file string) ([]City, error) {
	dat, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	text := string(dat)
	lines := strings.Split(text, "\n")

	var cities []City

	for _, x := range lines {
		city := ParseCity(x)
		cities = append(cities, city) 
	}

	return cities, nil

}  




func DecodeMap(fileName string) (*WorldMap, error) {
	
	cities, e := ParseCities(fileName)  

	if e != nil {
		return nil, e
	}

	nodes := InitCityMap() 
	
	for _, x := range cities {
		nodes.Put(x.name, x)
	}

	aliens := InitAlienMap()
	s := &WorldMap{nodes,aliens}

	s.Validate() 

	return s, e
	
	
}


func (WM *WorldMap) RandomCity() string {
	var cities []string
	mapCopy := WM.cities.Copy()
	
	for k, _ := range mapCopy {
		cities = append(cities, k) 
	}


	rand.Seed(time.Now().Unix())
	n := rand.Intn( len(cities) )
	return cities[n]
}  






func (WM *WorldMap) InitAliens(num int) *WorldMap {
	
	for i := 1; i <= num; i++ {
		ch := make(chan RMSG, 10)
		city := WM.RandomCity() 

		alien := Alien{id: i, ctr: 1, ch: ch, location: city}
		
		WM.aliens.Put(i, alien)
	}


	cityCopy := WM.cities.Copy()

	for k, _ := range cityCopy {
		WM.RemoveRivals(k) 
	}  

	return WM

}
