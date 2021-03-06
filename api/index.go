
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"strconv"
	"log"
	"time"
)

// Respuesta del bot
type Response struct {
	Msg string `json:"text"`
	ChatID int64 `json:"chat_id"`
	Method string `json:"method"`
}

// Representa un conjunto de objetos Inmueble
type Inmuebles struct {
	Inmuebles []Inmueble `json:"inmuebles"`
}

// Representa un objeto Inmueble individual
type Inmueble struct {
	Superficie   float32  `json:"superficie"`
	Habitaciones int      `json:"habitaciones"`
	Precio		 float32  `json:"precio"`
	Calle  		 string   `json:"calle"`
	Portal 		 int      `json:"portal"`
	Piso   		 int      `json:"piso"`
	Letra  		 string   `json:"letra"`
	Propietario  string   `json:"propietario"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson (url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Funcion serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	
	//Obtiene el payload
	defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var update tgbotapi.Update
    if err := json.Unmarshal(body,&update); err != nil {
        log.Fatal("Error en el update →", err)
	}
	
		// Si el mensaje es un comando
		if update.Message.IsCommand() {
			// Mensaje respuesta
			text := ""
			switch update.Message.Command() {
				// Establece el error
				case "autor":
					text = "Raúl Del Pozo Moreno, estudiante de la Universidad de Granada. Bot desarrollado para la asignatura IV."
				// Indica la cantidad de inmuebles disponibles
				case "cantidad":
					// Obtiene el fichero de datos
					var inmuebles Inmuebles
					getJson("https://inmobiliv.herokuapp.com/inmuebles", &inmuebles)

					// Responde segun cantidad
					if len(inmuebles.Inmuebles) == 0 {
						text = "No hay inmuebles registrados"
					} else {
						text = "La cantidad de inmuebles disponibles es de " + strconv.Itoa(len(inmuebles.Inmuebles))
					}

				// Indica todos los inmuebles disponibles
				case "todo":
					// Obtiene el fichero de datos
					var inmuebles Inmuebles
					getJson("https://inmobiliv.herokuapp.com/inmuebles", &inmuebles)

					if len(inmuebles.Inmuebles) == 0 {
						text = "No hay inmuebles registrados"
					} else {
						// Rellena la respuesta con los datos leidos
						text = "---\n"
						for i := 0 ; i < len(inmuebles.Inmuebles) ; i++ {
							var sup string = fmt.Sprintf("%.2f", inmuebles.Inmuebles[i].Superficie) + "m²"
							var hab string = strconv.Itoa(inmuebles.Inmuebles[i].Habitaciones)
							var pre string = fmt.Sprintf("%.2f", inmuebles.Inmuebles[i].Precio) + "€"
							var cal string = inmuebles.Inmuebles[i].Calle
							var por string = strconv.Itoa(inmuebles.Inmuebles[i].Portal)
							var pis string = strconv.Itoa(inmuebles.Inmuebles[i].Piso)
							var let string = inmuebles.Inmuebles[i].Letra
							var pro string = inmuebles.Inmuebles[i].Propietario
							text += "Direccion: " + cal + " Nº" + por + " " + pis + "º" + let
							text += "\nSuperficie: " + sup + "\nHabitaciones: " + hab + "\nPrecio: " + pre + "\nPropietario: " + pro
							text += "\n---\n"
						}
					}

				// Indica que no es un comando valido si no esta indicado arriba
				default:
					text = "No es valido"

			}

			// Genera la respuesta
			data := Response{ Msg: text, Method: "sendMessage", ChatID: update.Message.Chat.ID }

			// Genera el mensaje
			msg, _ := json.Marshal( data )
			log.Printf("Response %s", string(msg))
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w,string(msg))
		}

}

/*func main() {

	http.ListenAndServe(":3000", http.HandlerFunc(Handler))

}*/

