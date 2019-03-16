package socks5

import (
	"log"
	"net"

	"github.com/txthinking/x"
)

// Connect remote conn which u want to connect with your dialer
// Error or OK both replied.
func (r *Request) Connect(c *net.TCPConn, d x.Dialer) (*net.TCPConn, error) {
	if Debug {
		log.Println("Call:", r.Address())
	}
	tmp, err := d.Dial("tcp", r.Address())
	if err != nil {
		var p *Reply
		if r.Atyp == ATYPIPv4 || r.Atyp == ATYPDomain {
			p = NewReply(RepHostUnreachable, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = NewReply(RepHostUnreachable, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, err
	}
	rc := tmp.(*net.TCPConn)

	a, addr, port, err := ParseAddress(rc.LocalAddr().String())
	if err != nil {
		var p *Reply
		if r.Atyp == ATYPIPv4 || r.Atyp == ATYPDomain {
			p = NewReply(RepHostUnreachable, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = NewReply(RepHostUnreachable, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, err
	}
	p := NewReply(RepSuccess, a, addr, port)
	if err := p.WriteTo(c); err != nil {
		return nil, err
	}

	return rc, nil
}
