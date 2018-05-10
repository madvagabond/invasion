package invasion

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
		City := WM.cities[R]

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
			Neighbor := WM.cities[x.city]
		
			
			if e != nil {
				WM.cities[C.name] = C.RmConn(x.city)
				continue 
			}


			if b == true {continue} else {
				nconn := Connection{city: C.name, direction: direction}
				WM.cities[x.city] = Neighbor.AddConn(nconn)  
			}  

		}

		
	}  
	
	
	for _, v := range WM.cities {
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




func DecodeWorldMap(fileName string) (*WorldMap, error) {
	
	cities, e := ParseCities(fileName)  

	if e != nil {
		return nil, e
	}

	nodes := make(map[string]City)
	
	for _, x := range cities {
		nodes[x.name] = x
	}

	aliens := make(map[int]Alien)

	
	s := &WorldMap{cities: nodes, aliens: aliens}

	s.Validate() 

	return s, e
	
	
}


func (WM *WorldMap) RandomCity() string {
	var cities []string
	
	for k, _ := range WM.cities {
		cities = append(cities, k) 
	}


	rand.Seed(time.Now().Unix())
	n := rand.Intn( len(cities) )
	return cities[n]
}  






func (WM *WorldMap) InitAliens(num int) *WorldMap {
	
	aliens := make(map[int]Alien)
	
	for i := 1; i <= num; i++ {
		ch := make(chan RMSG)
		city := WM.RandomCity() 

		alien := Alien{id: i, moveCtr: 1, ch: ch, location: city}
		
		aliens[i] = alien
	}


	for k, _ := range WM.cities {
		p := func(a Alien) bool {return a.location == k}
		rivals := FilterAliens(WM.aliens, p) 

		if len(rivals) > 1 {
			WM.RemoveRivals(rivals, k)
		}
		
	}  

	
	WM.aliens = aliens
	return WM

}




