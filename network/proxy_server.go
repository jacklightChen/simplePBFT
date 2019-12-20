package network

import (
	"fmt"
	"github.com/simplePBFT/common"
	. "github.com/simplePBFT/config"
	. "github.com/simplePBFT/protocol"
	pb "github.com/simplePBFT/protos"
	"net"
)

type Server struct {
	listener net.Listener
	node *Node
	committees map[string]*common.Committee // committee is a sharding
	localCmt *common.Committee              // node belongs to localCmt
	config *Config
	internalMsgQueue chan pb.Message

	listener1 net.Listener
}

func NewServer(config *Config) *Server {

	server := &Server{}
	server.config = config
	server.CreateCommittees(config)
	node := NewNode(server)
	server.node = node
	server.internalMsgQueue = make(chan  pb.Message)

	return server
}

func (server *Server) Start() {
	fmt.Printf("Server will be started at %s...\n", server.config.LocalURL)
	if ln, err := net.Listen("tcp", string(server.config.LocalURL)); err == nil {
		server.listener = ln
		go server.Accept()
		go server.Receive()
	}
	for _,cmt := range server.committees {
		go cmt.Dails(server.config.LocalNodeID)
	}
}


func (server *Server) CreateCommittees(config *Config) {
	server.committees = make(map[string]*common.Committee)
	for _, cmtConfig := range config.CommitteeConfigs {
		cmt := server.CreateCommittee(cmtConfig)
		server.committees[cmtConfig.ID] = cmt
		if cmtConfig.ID == config.LocalCommittee {
			server.localCmt = cmt
		}
	}
}

func (server *Server) CreateCommittee(cmtConfig CommitteeConfig) *common.Committee {
	cmt := &common.Committee{
		CommitteeID: cmtConfig.ID,
		NodeInfos: make([]common.NodeInfo, len(cmtConfig.Nodes)),
	}
	for idx, node := range cmtConfig.Nodes {
		cmt.NodeInfos[idx].NodeID = node.NodeID
		cmt.NodeInfos[idx].NodeURL = node.URL
	}
	return cmt
}


func (server *Server) Accept() {
	for {
		fmt.Println("in accept")
		conn, err := server.listener.Accept()
		if err != nil {
			panic(err)
		}
		go server.handleConnect(conn)
	}
}

func (server *Server) handleConnect(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go server.reader(conn,readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			//	Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

func (server *Server) reader(conn net.Conn,readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			//fmt.Println(string(data))
			if msg, err := Unserial(data); err == nil {
				//fmt.Println("insert")
				server.internalMsgQueue <- msg
			}
		}
	}
}

func (server *Server) Receive() {
	for {
		select {
		case msg := <- server.internalMsgQueue:
			fmt.Println("receive msg1")
			if msg.Payload != nil {
				server.node.MsgEntrance <- msg
			}
		}
	}
}

func (server *Server) send(cmt string, id string, msg []byte) bool {
	if cmt == "" {
		return server.localCmt.Send(id, msg)
	}
	if cmt, ok := server.committees[cmt]; ok {
		return cmt.Send(id, msg)
	}
	return false
}