package invasion

type Server struct {
	World *WorldMap
	ch chan TMSG
	sig chan bool 
} 



func (Srv *Server) Handler() {
	
	for len(Srv.World.aliens) > 0 {
		
		msg := <- Srv.ch
		alien := Srv.World.aliens[msg.id]
		go Srv.World.Move(alien) 

	}

	
}




func (Alien *Alien) Worker(conn chan TMSG) {

	
	for { 
		req := TMSG{id: Alien.id}
		conn <- req

		rep := <- Alien.ch
		switch rep.status {

		case disconnect:
			close(Alien.ch)
			break


		case map_change:
			req1 := TMSG{id: Alien.id}
			conn <- req1
		}		
	}

}  


