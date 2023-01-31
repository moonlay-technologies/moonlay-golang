package rabbitmq

import (
	"sync"
	"time"
)

type ChannelPool interface {
	SetChannelIdleTimeout(timeout time.Duration)
	SetMinIdleChannel(minIdleChannel int)
	SetMaxPool(maxPool int)
	GetChannel() (*Channel, error)
	Done(target *Channel)
	CloseAll()
}

type channelPool struct {
	connection   *Connection
	serviceCode  string
	channelIndex int
	channelList  []*Channel
	config       *Config
	capacity     int
	muLock       sync.Mutex
}

type Config struct {
	MaxPool        int
	MinIdleChannel int
	IdleTimeout    time.Duration
}

func (c *Connection) NewChannelPool(serviceCode string) (ChannelPool, error) {
	channelList := make([]*Channel, 0)

	nChannel, err := c.Channel(serviceCode)
	if err != nil {
		return nil, err
	}
	channelList = append(channelList, nChannel)

	cp := &channelPool{
		serviceCode:  serviceCode,
		connection:   c,
		config:       &Config{MaxPool: 10, IdleTimeout: 5 * time.Second, MinIdleChannel: 1},
		channelList:  channelList,
		channelIndex: 0,
	}
	go cp.validateChannel()
	return cp, nil
}

func (cp *channelPool) SetChannelIdleTimeout(timeout time.Duration) {
	cp.muLock.Lock()
	cp.config.IdleTimeout = timeout
	cp.muLock.Unlock()
}

func (cp *channelPool) SetMinIdleChannel(minIdleChannel int) {
	cp.muLock.Lock()
	cp.config.MinIdleChannel = minIdleChannel
	cp.muLock.Unlock()
}

func (cp *channelPool) SetMaxPool(maxPool int) {
	cp.muLock.Lock()
	cp.config.MaxPool = maxPool
	cp.muLock.Unlock()
}

func (cp *channelPool) addNewChannel() error {
	nChannel, err := cp.connection.Channel(cp.serviceCode)
	if err != nil {
		return err
	}
	cp.channelList = append(cp.channelList, nChannel)
	return nil
}

func (cp *channelPool) countChannel() int {
	cp.muLock.Lock()
	defer cp.muLock.Unlock()
	return len(cp.channelList)
}

func (cp *channelPool) canAddNewChannel() bool {
	cp.muLock.Lock()
	defer cp.muLock.Unlock()
	return len(cp.channelList) < cp.config.MaxPool
}

func (cp *channelPool) addIndex() {
	if cp.channelIndex < cp.config.MaxPool-1 {
		cp.channelIndex++
	} else {
		cp.channelIndex = 0
	}
}

func (cp *channelPool) GetChannel() (*Channel, error) {
	cp.muLock.Lock()
	defer cp.muLock.Unlock()
	defer cp.addIndex()
	if len(cp.channelList) < cp.config.MaxPool {
		if err := cp.addNewChannel(); err != nil {
			return nil, err
		}
	}
	ch := cp.channelList[cp.channelIndex]
	//fmt.Printf("Get Pool Object with ID: %s\n", ch.ID)
	return ch, nil
}

func (cp *channelPool) validateChannel() {
	for {
		skipUnlock := false
		cp.muLock.Lock()
		for i, obj := range cp.channelList {
			currentActiveLength := len(cp.channelList)
			//fmt.Printf("i :%d \t length : %d \t duration : %s \n", i, currentActiveLength, time.Since(obj.lastUsed))
			if i >= currentActiveLength || currentActiveLength == cp.config.MinIdleChannel {
				cp.muLock.Unlock()
				skipUnlock = true
				break
			}
			if time.Since(obj.lastUsed) >= cp.config.IdleTimeout {
				if err := cp.channelList[i].Close(); err == nil {
					cp.channelList[currentActiveLength-1], cp.channelList[i] = cp.channelList[i], cp.channelList[currentActiveLength-1]
					cp.channelList = cp.channelList[:currentActiveLength-1]
					cp.channelIndex = 0
				}
			}
		}
		if !skipUnlock {
			cp.muLock.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}

func (cp *channelPool) CloseAll() {
	for i := 0; i < len(cp.channelList); i++ {
		cp.muLock.Lock()
		cp.channelList[i].Close()
		cp.channelList = nil
		cp.muLock.Unlock()
	}
	time.Sleep(1 * time.Second)
}

func (cp *channelPool) Done(target *Channel) {
	cp.muLock.Lock()
	defer cp.muLock.Unlock()
	for i, obj := range cp.channelList {
		if obj.ID == target.ID {
			cp.channelList[i].LastUsedNow()
		}
	}
	//fmt.Printf("Mark as done Object with ID: %s\n", target.ID)
	return
}

func (cp *channelPool) destroyChannel() {

}
