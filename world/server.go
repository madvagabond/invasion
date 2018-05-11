package world
import "fmt"


func (Srv *Server) Handler() {
	

	for {

		if Srv.World.aliens.Len() == 0 {
			fmt.Println("Server exiting")
			Srv.Sig <- true 
			return
		}
		
	
		msg := <-Srv.ch
		Srv.Move(msg.id)
			
	}
		//numAliens :=
	
	
}

	

	



func MakeServer(worldMap *WorldMap) *Server {
	return &Server{worldMap, make(chan TMSG, 30000), make(chan bool)}
}   




func (Alien *Alien) Worker(conn chan TMSG) {

	
	for { 
		req := TMSG{id: Alien.id}
		conn <- req

		rep := <- Alien.ch

		
		switch rep.status {

		case map_change:
			req1 := TMSG{id: Alien.id}
			conn <- req1

		case disconnect:
			return

		default:
			return

		}

		
	}

}  

func (Srv *Server) SpawnWorkers() {
	m1 := Srv.World.aliens.Copy()
	
	for _, a := range m1 {
		a1 := a
		go a1.Worker(Srv.ch) 
	}

}  






//The main simulation function
func (Srv *Server) Move(alienID int) {

	WM := Srv.World


	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			return
		}
	}()

	alien, e := WM.aliens.Get(alienID)
	if e != true {return} 


	
	city_name := alien.location


	current_city, e1 := WM.cities.Get(city_name)

	if e1 != true {
	//	fmt.Printf("Alien %d died in the destruction of %s\n", alien.id, city_name)
		WM.RemoveAlien(alien)
		return
	}
	


	if alien.ctr >= 10000 {
		WM.RemoveAlien(alien)
		fmt.Printf("Alien %d was removed due to making 10000 moves \n", alien.id) 
		return 
	}

	new_conn := current_city.RandomConn()
	


	if new_conn == nil {
		WM.RemoveAlien(alien)
		fmt.Printf("Alien %d, was trapped in %s and starved to death \n", alien.id, city_name) 
		return 
	}

	new_cname := new_conn.city
	alien.ctr++
	a1 := Alien{alien.id, alien.ctr, new_cname, alien.ch}

	
	WM.aliens.Put(a1.id, a1)

	b := WM.RemoveRivals(new_cname)
	
	if b != true {
		a1.ch <- RMSG{map_change}
	}

	
}  
