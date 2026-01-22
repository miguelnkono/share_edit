package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"share_edit/models"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Sender", "Viewer"}

type server struct {
	User models.User

	startTime    int
	server_entry net.Listener
	cursor       int
	choice       string
	showMenu     bool
}

func (s server) Init() tea.Cmd {
	if s.startTime > 0 && !s.showMenu{
		return tick
	}

	return nil
}

func (s server) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl + q":
			return s, tea.Quit

		// cursor position
		case "enter":
			if s.showMenu {
				s.choice = choices[s.cursor]
				return s, tea.Quit
			}

		case "down", "j":
			if s.showMenu {
				s.cursor++
				if s.cursor >= len(choices) {
					s.cursor = 0
				}
				return s, nil
			}

		case "up", "k":
			if s.showMenu {
				s.cursor--
				if s.cursor < 0 {
					s.cursor = len(choices) - 1
				}
				return s, nil
			}
		}

		// the tick message
	case tickMsg:
		s.startTime--
		if s.startTime <= 0 {
			// ask the user to choice what his role in the session of communication is
			s.showMenu = true
			return s, tea.Quit
		}
		return s, tick
	}

	return s, nil
}

func (s server) View() string {
	if !s.showMenu {
		return fmt.Sprintf("Your share_edit program will start in %d ...\n\nTo quick now press ctrl + q\n", s.startTime)
	}

	str := strings.Builder{}
	str.WriteString("Who are you?\n\n")
	for i := range choices {
		if s.cursor == i {
			str.WriteString("(•) ")
		} else {
			str.WriteString("( ) ")
		}
		str.WriteString(choices[i])
		str.WriteString("\n")
	}
	str.WriteString("\n(Press ctrl + q to quit)\n")

	return str.String()
}

func main() {
	logfilePath := os.Getenv("BUBBLETEA_LOG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "logs"); err != nil {
			log.Fatal(err)
		}
	}

	program := tea.NewProgram(initModel())

	// Run returns the model as a tea.Model.
	m, err := program.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(server); ok && m.choice != "" {
		fmt.Printf("\n---\nYou chose %s!\n", m.choice)
	}

}

// initialize the server model
func initModel() server {

	return server{startTime: 5, cursor: 0, showMenu: false}
}

// the tick msg
type tickMsg struct{}

func tick() tea.Msg {
	time.Sleep(1 * time.Second)
	return tickMsg{}
}

// the error message
type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

// func init_server() tea.Msg {

// 	// listener for new connections
// 	listener, err := net.Listen("tcp", "localhost:8080")
// 	if err != nil {
// 		return errMsg{err}
// 	}
// 	defer listener.Close()

// 	// accept new connections
// 	for {
// 		connection, err := listener.Accept()
// 		if err != nil {
// 			log.Fatal("Error while waiting for new connections: ", err)
// 			continue
// 		}

// 		// handle each connection in a new goroutine
// 		fmt.Printf("%s connected\n", connection.LocalAddr().String())
// 		go handle_connection(connection)
// 	}

// 	return server{server_entry: listener}
// }

// func main() {
// }

// func handle_connection(connection net.Conn) {
// 	defer connection.Close()

// 	reader := bufio.NewReader(connection)
// 	message, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal("Error while reading data: %v", err)
// 		return
// 	}

// 	message_to_send := strings.ToUpper(strings.TrimSpace(message))
// 	response := fmt.Sprintf("ACK: %s\n", message_to_send)
// 	_, err = connection.Write([]byte(response))
// 	if err != nil {
// 		log.Fatal("Error while sending the response to the client: %v", err)
// 	}
// }
