package main
import (
"fmt"
"net"
"os"
"strconv"
"time"
"bufio"
)
//Variáveis globais interessantes para o processo
var err string
var myPort string //porta do meu servidor
var nServers int //qtde de outros processo
var CliConn []*net.UDPConn //vetor com conexões para os servidores
 //dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo
 //mensagens dos outros processos)
var id int //numero identificador do processo
var logicalClock int

func CheckError(err1 error){
	if err1 != nil {
		fmt.Println("Erro: ", err1)
		os.Exit(0)
	}
}

func PrintError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
	}
}

func MaxInt(x1, x2 int) int {


	if x1 > x2 {
		return x1
	}
	return x2
}

func doServerJob() {
//Ler (uma vez somente) da conexão UDP a mensagem
//Escreve na tela a msg recebida


	 buf := make([]byte, 1024)
	 var aux int
	 for {

		 n, addr, err := ServConn.ReadFromUDP(buf)
		 aux ,_ = strconv.Atoi(string (buf[0:n]))
   	 	 logicalClock = MaxInt(logicalClock ,aux ) + 1
		 fmt.Println("Received ",logicalClock, " from ",addr)

		 if err != nil {
		 	fmt.Println("Error: ",err)
		 } 
 	}
}
func doClientJob(otherProcess int, i int) {
//Envia uma mensagem (com valor i) para o servidor do processo
//otherServer
    //defer Conn.Close()
    msg := strconv.Itoa(i)
     buf := []byte(msg)
     _,err := CliConn[otherProcess].Write(buf)
     if err != nil {
        fmt.Println(msg, err)
    }

}



func initConnections() {
	id, _ = strconv.Atoi(os.Args[1])
	myPort = os.Args[ id + 1]
	nServers = len(os.Args) - 2
	/*Esse 2 tira o nome (no caso Process) e tira a primeira porta
	(que é a minha). As demais portas são dos outros processos*/
	//Outros códigos para deixar ok a conexão do meu servidor
	//Outros códigos para deixar ok as conexões com os servidores
	//dos outros processos
	connections := make([]*net.UDPConn, nServers, nServers)
	
	for i:=0; i<nServers; i++ {

		port := os.Args[i+2]
		
			ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1" + string (port) )
			PrintError(err)
 
    		LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
 		   	PrintError(err)
 
    		connections[i], err = net.DialUDP("udp",LocalAddr, ServerAddr)
			PrintError(err)
		
	}
	CliConn = connections

	 /* Lets prepare a address at any address at port 10001*/   
	 ServerAddr,err := net.ResolveUDPAddr("udp", myPort)
	 CheckError(err)
	
	 /* Now listen at selected port */
	 ServConn, err = net.ListenUDP("udp", ServerAddr)
	 CheckError(err)
}

func readInput(ch chan string) {
	// Non-blocking async routine to listen for terminal input
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
	}
}
	


func main(){
	initConnections()
	//O fechamento de conexões devem ficar aqui, assim só fecha
	//conexão quando a main morrer
	defer ServConn.Close()
	for i := 0; i < nServers; i++ {
		defer CliConn[i].Close()
	}
	//Todo Process fará a mesma coisa: ouvir msg e mandar infinitos
	//i’s para os outros processos
	ch := make(chan string)
	logicalClock = 1
	go readInput(ch)
	for {
			
		//Server
		go doServerJob()
		// When there is a request (from stdin). Do it!
		select {
			case x, valid := <- ch :
				if valid {
                        compare,_ := strconv.Atoi(x)
						if ( compare != id){

							fmt.Printf("Enviado %d para %s  \n",logicalClock, x)
							doClientJob(compare-1 , logicalClock)
						} else{

							logicalClock++
							fmt.Printf("Atualizado logicalClock para %d \n",logicalClock)
						}
				} else {
						 
					fmt.Println("Channel closed!")
						
				}
				
			default:
			
				// Do nothing in the non-blocking approach.
			
				time.Sleep(time.Second * 1)
		}
			
		// Wait a while
		time.Sleep(time.Second * 1)
	}
}
		