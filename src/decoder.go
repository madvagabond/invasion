package invasion
import (
	"github.com/BurntSushi/toml"
	"errors"
	"strings"
	"math/rand"
	"time"
)


type Server struct {
	state *WorldMap
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
func (WM *WorldMap) Validate() *WorldMap {


	check := func(L string, R string) bool {
		City := WM.cities[R]

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
			Neighbor := WM.cities[x.cityName]
		
			
			if e != nil {
				WM.cities[C.name] = C.RmConn(x.cityName)
				continue 
			}


			if b == true {continue} else {
				nconn := Connection{cityName: C.name, direction: direction}
				WM.cities[x.cityName] = Neighbor.AddConn(nconn)  
			}  

		}

		
	}  
	
	
	for _, v := range WM.cities {
		symmetry(v) 
	}

	return WM
	
}





func DecodeWorldMap(fileName string) (*WorldMap, error) {
	
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
	n := rand.Int() %  len(cities)
	return cities[n]
}  



func (WM *WorldMap) initAliens(num int) *WorldMap {
	aliens := make(map[int]Alien)
	
	for i := 1; i <= num; i++ {
		ch := make(chan RMSG)
		city := WM.RandomCity() 

		alien := Alien{id: i, moveCtr: 0, ch: ch, location: city}
		
		aliens[i] = alien
	}

	WM.aliens = aliens
	return WM

}




