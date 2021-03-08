package driver

const (
	//GPIOOUT direction OUT
	GPIOOUT string= "out"
	//GPIOIN direction IN
	GPIOIN string="in" 

	HIGH = 1
    LOW  = 0

	MAXUNCHANGECOUNT = 100

	StateInitPullDown = 1
	StateInitPullUp = 2
	StateDataFirstPullDown = 3
	StateDataPullUp = 4
	StateDataPullDown = 5

	PUDOFF = 0
	PUDDOWN = 1
	PUDUP = 2
)