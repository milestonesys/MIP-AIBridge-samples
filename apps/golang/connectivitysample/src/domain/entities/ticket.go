package entities

// Holds data for a websocket request
type Ticket struct {
	token    string
	innerURL string
}

func NewTicket(token, innerURL string) *Ticket {
	return &Ticket{
		token:    token,
		innerURL: innerURL,
	}
}

func (t *Ticket) Token() string {
	return t.token
}

func (t *Ticket) InnerURL() string {
	return t.innerURL
}
