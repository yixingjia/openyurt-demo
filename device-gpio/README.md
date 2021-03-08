# GPIO Device Service
## Overview
GPIO Micro Service - device service for connecting GPIO devices to EdgeX.

- Function:
  - This device service use sysfs to control GPIO devices. For a just connected GPIO device, export the needed pin number, set correct GPIO direction, and then start to read or write data from GPIO device.
- Physical interface: system gpio (/sys/class/gpio)
- Driver protocol: IO



## Usage
- This Device Service have to run with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command.
- After the service started, we can use "exportgpio" command to export your GPIO device, then set direction with "gpiodirection" command, and use "gpiovalue" to read or write data.


## Guidance
Here we give two step by step guidance examples of using this device service. In these examples, we use RESTful API to interact with EdgeX.

Before we actually operate GPIO devices, we need to find out RESTful urls of this running device service. By using

`curl http://localhost:48082/api/v1/device/name/gpio`

We can find out the each command url you needed in next steps from `crul` response. Other than simply use `curl`, we can also use tool like `Postman`, which will give better experience on sending RESTful requests. Here, we just provide one possible way to interact with this device service, please feel free to use any tool you like.

### Write value to GPIO
Assume we have a GPIO device connected to pin 134 on current system.

1. Export GPIO pin number

    `curl -H "Content-Type: application/json" -X PUT -d '{"export": "134"}'  http://localhost:48082/api/v1/device/47eb5168-ec63-4be9-bd9a-81d3b2d3ae75/command/57d5cab7-a869-4a4f-a6b4-7bf98f66932`

    By replacing `134` to another value, you can set up your GPIO device service. This number is normally depends on your mother board and which pin the device is conneted to. When this command is successfully executed, you can use `ls -l /sys/class/gpio` to verify, you should see a file named with `gpio134` in the result of `ls` command.

2. Set out direction

    `curl -H "Content-Type: application/json" -X PUT -d '{"direction":"out"}' http://localhost:48082/api/v1/device/47eb5168-ec63-4be9-bd9a-81d3b2d3ae75/command/4c32eb83-e21d-4d09-bcd7-235adc460cba`

    Since we are going to write some value to a exported GPIO, then we just set the direction to `out`.

3. Write value

    `curl -H "Content-Type: application/json" -X PUT -d '{"value":"1"}' http://localhost:48082/api/v1/device/47eb5168-ec63-4be9-bd9a-81d3b2d3ae75/command/5c612311-10c0-4659-890c-4199b47ef988`

    Now if you test pin 134, it is outputing high voltage.


### Read value from GPIO
Assume we have another GPIO device connected to pin 134 on current system.

1. Export GPIO pin number

    `curl -H "Content-Type: application/json" -X PUT -d '{"export": "134"}'  http://localhost:48082/api/v1/device/47eb5168-ec63-4be9-bd9a-81d3b2d3ae75/command/57d5cab7-a869-4a4f-a6b4-7bf98f66932`

2. Set in direction

    `curl -H "Content-Type: application/json" -X PUT -d '{"direction":"in"}' http://localhost:48082/api/v1/device/47eb5168-ec63-4be9-bd9a-81d3b2d3ae75/command/4c32eb83-e21d-4d09-bcd7-235adc460cba`

    Since we are going to read some data from this device, then we just set the direction to `in`.

3. Read value

    `curl http://localhost:48082/api/v1/device/47eb5168-ec63-4be9-bd9a-81d3b2d3ae75/command/5c612311-10c0-4659-890c-4199b47ef988`

    Here, we post some results:

```bash
$ curl http://localhost...890c-4199b47ef988
{"device":"gpio","origin":1611752289806843150,"readings":[{"origin":1611752289806307945,"device":"gpio","name":"value","value":"{\"gpio\":134,\"value\":0}","valueType":"String"}],"EncodedEvent":null}

$ curl http://localhost...890c-4199b47ef988
{"device":"gpio","origin":1611752309686651113,"readings":[{"origin":1611752309686212741,"device":"gpio","name":"value","value":"{\"gpio\":134,\"value\":1}","valueType":"String"}],"EncodedEvent":null}
```


## API Reference

| Method | Core Command  | parameters                | Description                                                  | Response                             |
| ------ | ------------- | ------------------------- | ------------------------------------------------------------ | ------------------------------------ |
| put    | exportgpio    | {"export":<gpionum>}      | Export a gpio from "/sys/class/gpio"<br><gpionum>: int, gpio number | 200 ok                               |
| put    | unexportgpio  | {"unexport":<gpionum>}    | Export a gpio from "/sys/class/gpio"<br/><gpionum>: int, gpio number | 200 ok                               |
| put    | gpiodirection | {"direction":<direction>} | Set direction for the exported gpio<br/><direction>: string, "in" or "out" | 200 ok                               |
| get    | gpiodirection |                           | Get direction of the exported gpio                           | "{\"direction\":\"in\",\"gpio\":65}" |
| put    | gpiovalue     | {"value":<value>}         | Set value for the exported gpio<br/><value>: int, 1 or 0     | 200 ok                               |
| get    | gpiovalue     |                           | Get value of the exported gpio                               | "{\"gpio\":65,\"value\":1}"          |



## License
[Apache-2.0](LICENSE)

