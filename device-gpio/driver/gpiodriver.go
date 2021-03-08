package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type GPIODev struct {
	lc     logger.LoggingClient
	gpio int
	humidity int32
	temperature int32
}

func NewGPIODev(lc logger.LoggingClient) *GPIODev {
	return &GPIODev{lc: lc, gpio: -1}
}


func (dev *GPIODev) ExportGPIO(gpio int) error {
	err := exportgpio(gpio)
	if err == nil {
		dev.gpio = gpio
	}
	return err
}

func (dev *GPIODev) UnexportGPIO(gpio int) error {
	err := unexportgpio(gpio)
	if err == nil {
		dev.gpio = -1
	}
	return err
}

func (dev *GPIODev) SetDirection(direction string) error {
	if dev.gpio == -1 {
		return errors.New("Please export gpio first")
	}
	return setgpiodirection(dev.gpio, direction)
}

func (dev *GPIODev) GetDirection() (string,error) {
	if dev.gpio == -1 {
		return "", errors.New("Please export gpio first")
	}
	direction, err := getgpiodirection(dev.gpio)
	if err != nil {
		return "", err
	}  else {
		res, _ := json.Marshal(map[string]interface{}{"gpio": dev.gpio, "direction": direction})
		return string(res), err
	}
}

func (dev *GPIODev) GetHumidity()(int32, error){
	if dev.gpio == -1 {
		return 0, errors.New("Please export gpio first")
	}
	//TODO
	humidity,temperature, err:= dev.GetTempAndHumidity()
	if(err != nil){
		return 0,err
	}
	dev.humidity = humidity
	dev.temperature = temperature
	return dev.humidity,nil
}
func (dev *GPIODev) GetTemperature()(int32, error){
	if dev.gpio == -1 {
		return 0, errors.New("Please export gpio first")
	}
	//TODO
	return dev.temperature, nil
}
func (dev *GPIODev) GetTempAndHumidity()(int32,int32,error){
	out,err:=exec.Command("./DTH11.py",string(dev.gpio)).Output()
	if(err != nil){
		return 0,0,err
	}
	tempString := string(out)
	tempString = strings.Trim(tempString,"\n")
	ht := strings.Split(tempString,",")
	if len(ht) != 2{
		return 0,0, errors.New("size Invalid temperature and humidity value "+tempString+ string(len(ht)))
	}
	if humidity, err:= strconv.Atoi(ht[0]); err ==nil{
		if temperature, err:= strconv.Atoi(ht[1]); err==nil{
			return int32(humidity),int32(temperature), nil
		}else{
			return 0,0, err
		}
	}else{
		return 0,0, err
	}
	return 0,0,errors.New("Invalid temperature and humidity value"+tempString)
}
/**
func (dev *GPIODev) GetTempAndHumidity()(int32,float64,error){
	//assume the gpio is already set
	// this method should only be called by GetTemperature or GetHumidity
	dev.lc.Info(fmt.Sprintf("set pin port: %v",dev.gpio))
	err :=rpio.Open()
	if err != nil{
		return 0,0,err
	}
	defer rpio.Close()
	pin :=rpio.Pin(17)
	pin.Output()
	pin.High()
	time.Sleep(50 * time.Millisecond)
	pin.Low()
	time.Sleep(20 * time.Millisecond)
	pin.Input()

	unchangedCount :=0
	last :=int8(-1)
	data := make([]int8,1)
	for {
		current := int8(pin.Read())
		data=append(data,current)
		if last != current {
			unchangedCount = 0
			last = current
		}else{
			unchangedCount++
			if unchangedCount > MAXUNCHANGECOUNT{
				break
			}
		}

	}
	lengths := make([]int,10)
	state := StateInitPullDown
	currentLength :=0
	for _, current :=range data{
		currentLength++
		if state == StateInitPullDown {
			if current == LOW{
				state = StateInitPullUp
			}else{
				continue
			}
		}
		if state == StateInitPullUp {
			if current ==HIGH {
				state = StateDataFirstPullDown
			}else{
				continue
			}
		}
		if state == StateDataFirstPullDown {
			if current == LOW {
				state = StateDataPullUp
			}else{
				continue
			}
		}
		if state == StateDataPullUp{
			if current == HIGH {
				currentLength = 0
				state = StateDataPullDown
			}else{
				continue
			}
		}
		if state == StateDataPullDown{
			if current == LOW{
				lengths = append(lengths, currentLength)
				state = StateDataPullUp
			}else{
				continue
			}
		}
	}
	if len(lengths) != 40{
		return 0,0,errors.New("Data not good, skip")
	}
	shortestPullUp,longestPUllUp := MinMax(lengths)
	halfWay := (shortestPullUp +longestPUllUp)/2
	dataBits := make([]byte,1)
	theBytes := make([]byte,1)
	b :=byte(0)

	for _, length := range lengths{
		bit :=byte(0)
		if length > halfWay{
			bit=1
		}
		dataBits = append(dataBits,bit)
	}
	for index,bit := range dataBits{
		b = b << 1
		if bit==1{
			b = b | 1
		}else{
			b = b | 0
		}
		if(index +1)%8 == 0{
			theBytes = append(theBytes, b)
			b=0
		}
	}
	checkSum:=(theBytes[0]+theBytes[1]+theBytes[2]+theBytes[3])&0xFF
	if theBytes[4]!= checkSum{
		dev.lc.Info("Data not good, skip")
		return 0,0,errors.New("Data not good, skip")
	}

	//setgpiodirection(dev.gpio, GPIOOUT)
	//setgpiovalue(dev.gpio,HIGH)
	//time.Sleep(50 * time.Millisecond)
	//setgpiovalue(dev.gpio,LOW)
	//time.Sleep(20 * time.Millisecond)
	
	return int32(theBytes[0]),float64(theBytes[1]),nil
}
*/
func MinMax(array []int) (int, int) {
    var max int = array[0]
    var min int = array[0]
    for _, value := range array {
        if max < value {
            max = value
        }
        if min > value {
            min = value
        }
    }
    return min, max
}
func (dev *GPIODev) SetGPIO(value int) error {
	if dev.gpio == -1 {
		return errors.New("Please export gpio first")
	}
	direction, err := getgpiodirection(dev.gpio)
	if err != nil {
		return err
	}
	if strings.Contains(direction, "in")  {
		return errors.New("Can not set the gpio which is input state")
	}

	return setgpiovalue(dev.gpio, value)
}

func (dev *GPIODev) GetGPIO() ( string, error ) {

	if dev.gpio == -1 {
		return "", errors.New("Please export gpio first")
	}
	gpiovalue, err := getgpiovalue(dev.gpio)
	if err != nil {
		return "", err
	}  else {
		res, _ := json.Marshal(map[string]interface{}{"gpio": dev.gpio, "value": gpiovalue})
		return string(res), err
	}
}


func exportgpio(gpioNum int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	} else {
		return ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	}
}

func unexportgpio(gpioNum int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ioutil.WriteFile("/sys/class/gpio/unexport", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	} else {
		return nil
	}
}

func setgpiodirection(gpioNum int, direction string) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var way string
		if direction == "in" {
			way = "in"
		} else {
			way = "out"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum), []byte(way), 0644)
	} else {
		return errors.New("Please export gpio first")
	}
}

func getgpiodirection(gpioNum int) (string, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		direction, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum))
		if err != nil {
			return "", err
		} else {
			return strings.Replace(string(direction), "\n", "", -1), err
		}
	} else {
		return "", errors.New("Please export gpio first")
	}
}

func setgpiovalue(gpioNum int, value int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var tmp string
		if value == 0 {
			tmp = "0"
		} else {
			tmp = "1"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum), []byte(tmp), 0644)
	} else {
		return errors.New("Please export gpio first")
	}
}

func getgpiovalue(gpioNum int) (int, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ret, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum))
		if err != nil {
			return 0, err
		} else {
			value, _ := strconv.Atoi(strings.Replace(string(ret), "\n", "", -1))
			return value, err
		}
	} else {
		return -1, errors.New("Please export gpio first")
	}
}
