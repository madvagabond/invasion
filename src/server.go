package state
import (
	"github.com/BurntSushi/toml"
	"errors"
	"strings"
)


type Server struct {
	state *CityMap
} 


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
func (CM *CityMap) Validate() *CityMap {


	check := func(L string, R string) bool {
		City := CM.nodes[R]

		pred := func(C Connection) bool {
			return C.cityName == L
		}

		return City.Contains(pred)		
	}




	
	symmetry := func(C City) {
		conns := C.connections

		for _, x := range conns {
			b := check(C.name, x.cityName)

			direction, e := Opposite(x.direction)
			Neighbor := CM.nodes[x.cityName]
		
			
			if e != nil {
				CM.nodes[C.name] = C.rmConn(x.cityName)
				continue 
			}


			if b == true {continue} else {
				nconn := Connection{cityName: C.name, direction: direction}
				CM.nodes[x.cityName] = Neighbor.addConn(nconn)  
			}  

		}

		
	}  
	
	
	for _, v := range CM.nodes {
		symmetry(v) 
	}

	return CM
	
}





func DecodeCityMap(fileName string) (*CityMap, error) {
	
	var cities []City
	_, e := toml.DecodeFile(fileName, &cities) 

	if e != nil {
		return nil, e
	}

	nodes := make(map[string]City)
	
	for _, x := range cities {
		nodes[x.name] = x
	}

	aliens := make(map[int]Alien)

	
	s := &CityMap{nodes: nodes, aliens: aliens}

	s.Validate() 

	return s, e
	
	
}  


func (CM *CityMap) initAliens(num int) *CityMap {
	aliens := make(map[int]Alien)
	
	for i := 1; i <= num; i++ {
		alien := Alien{id: i, moveCtr: 0} 
		aliens[i] = alien
	}

	CM.aliens = aliens
	return CM

}




