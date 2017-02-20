package controllers

import (
	"fmt"
	"net/http"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	api "github.com/soprasteria/docktor/model"
	"github.com/soprasteria/docktor/model/types"
	"github.com/soprasteria/docktor/server/daemons"
	"github.com/soprasteria/docktor/server/redisw"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"sync"
	"time"
)

// Daemons contains all daemons handlers
type Daemons struct {
}

type Action string

const RequestDaemonInfo Action = "REQUEST_DAEMON_INFO"
const ReceiveDaemonInfo Action = "RECEIVE_DAEMON_INFO"
const InvalidRequestDaemonInfo Action = "INVALID_REQUEST_DAEMON_INFO"

type GetInfosData struct {
	Daemon types.Daemon `json:"daemon"`
	Force  bool         `json:"force"`
}

type GetInfosDataResponse struct {
	Daemon types.Daemon        `json:"daemon"`
	Info   *daemons.DaemonInfo `json:"info"`
}

type MessageFromClient struct {
	Action Action      `json:"action"`
	Data   interface{} `json:"data"`
}

type MessageToClient struct {
	Action `json:"action"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
}

//OpenWS open websocket for daemons
func (d *Daemons) OpenWS(c echo.Context) error {
	// Upgrade the HTTP connection to Websocket protocol
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.WithError(err).Error("Can't upgrade daemons websocket")
		return c.String(http.StatusBadRequest, "Can't upgrade daemons websocket : "+err.Error())
	}
	// Channel receiveing event to write response to client
	responseChannel := make(chan MessageToClient)
	var wgGracefulClose sync.WaitGroup
	defer func() {
		err := conn.Close()
		if err != nil {
			log.WithError(err).Error("Unable to close the daemon websocket")
		}
		// Waiting for all current goroutine to end before closing channel, avoiding panic when goroutine send data on a closed channel
		wgGracefulClose.Wait()
		close(responseChannel)
		log.Debug("Channel for daemon responses is closed")
	}()

	// Gettings responses to write in channel infinite loop
	// Because writing in a websocket is not thread-safe. Only one goroutine for writer and one for reader is allowed.
	go func() {
		for {
			if response, ok := <-responseChannel; ok {
				err := conn.WriteJSON(response)
				if err != nil {
					if err != websocket.ErrCloseSent {
						log.WithError(err).WithField("message", response).Error("Unable to sent message to client")
					} else {
						log.WithField("message", response).Debug("Unable to send message to a closed websocket")
					}
				} else {
					log.WithField("message", response).Debug("Sent message to client")
				}
			} else {
				return
			}
		}
	}()

ContinuousRead:
	for {
		var message MessageFromClient
		// Read message synchronously
		err := conn.ReadJSON(&message)
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				log.Debugf("Websocket %p is deconnected.", conn)
				return err
			default:
				log.WithError(err).Warn("Message received is not a JSON message")
				continue ContinuousRead
			}
		}

		log.WithField("message", message).Debug("Received message from client")

		// But dispatch processing in goroutines to avoid other message to stack
		switch message.Action {
		case RequestDaemonInfo:
			go d.GetInfoFromMessage(c, message, conn, responseChannel, &wgGracefulClose)
			break
		default:
			log.WithField("action", message.Action).WithField("data", message.Data).Warn("Unrecognized type of message")
			break
		}

	}
}

func (d *Daemons) GetInfoFromMessage(c echo.Context, message MessageFromClient, conn *websocket.Conn, responseChannel chan MessageToClient, wgGracefulClose *sync.WaitGroup) {
	// Tell the main process, the execution of one action begins.
	wgGracefulClose.Add(1)
	// Tell the main process, the execution of one action just ended
	defer wgGracefulClose.Done()

	redisClient := redisw.GetRedis(c)
	// Parse and convert the message
	var data GetInfosData
	bytes, err := json.Marshal(message.Data)
	if err != nil {
		log.WithError(err).WithField("action", RequestDaemonInfo).Error("JSON received is syntaxically correct but not valid semantically")
		return
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.WithError(err).WithField("action", RequestDaemonInfo).Error("JSON received is syntaxically correct but not valid semantically")
		return
	}
	// Get daemon info (status, number of containers,...)
	infos, err := daemons.GetInfo(data.Daemon, redisClient, data.Force)
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	response := MessageToClient{
		Action: ReceiveDaemonInfo,
		Data:   GetInfosDataResponse{Daemon: data.Daemon, Info: infos},
	}
	if err != nil {
		response.Error = err.Error()
	}

	// Send the response to the channel
	responseChannel <- response
}

//GetAll daemons from docktor
func (d *Daemons) GetAll(c echo.Context) error {
	docktorAPI := c.Get("api").(*api.Docktor)
	daemons, err := docktorAPI.Daemons().FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all daemons")
	}
	return c.JSON(http.StatusOK, daemons)
}

//Save daemon into docktor
func (d *Daemons) Save(c echo.Context) error {
	docktorAPI := c.Get("api").(*api.Docktor)
	var daemon types.Daemon
	err := c.Bind(&daemon)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error while binding daemon: %v", err))
	}

	// If the ID is empty, it's a creation, so generate an object ID
	if daemon.ID.Hex() == "" {
		daemon.ID = bson.NewObjectId()
	}

	res, err := docktorAPI.Daemons().Save(daemon)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while saving daemon: %v", err))
	}
	return c.JSON(http.StatusOK, res)
}

//Delete daemon into docktor
func (d *Daemons) Delete(c echo.Context) error {
	docktorAPI := c.Get("api").(*api.Docktor)
	id := c.Param("daemonID")
	res, err := docktorAPI.Daemons().Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while remove daemon: %v", err))
	}
	return c.String(http.StatusOK, res.Hex())
}

//Get daemon from docktor
func (d *Daemons) Get(c echo.Context) error {
	daemon := c.Get("daemon").(types.Daemon)
	return c.JSON(http.StatusOK, daemon)
}

// GetInfo : get infos about daemon from docker
func (d *Daemons) GetInfo(c echo.Context) error {
	daemon := c.Get("daemon").(types.Daemon)
	forceParam := c.QueryParam("force")
	redisClient := redisw.GetRedis(c)

	infos, err := daemons.GetInfo(daemon, redisClient, forceParam == "true")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, infos)
}
