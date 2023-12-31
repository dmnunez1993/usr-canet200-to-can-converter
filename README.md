# USR Canet 200 to Can Converter

Converts USR Canet 200 data to standard can packaging format. Only works on Linux Based systems

### Dependencies

- [golang](https://go.dev/) - Version 1.21.5
- [NodeJS](https://nodejs.org) - Version 21.0.0
- [python3](https://www.python.org/)
- [can-utils](https://github.com/linux-can/can-utils)
- [docker](https://www.docker.com/) - Used to create .deb packages

### Start developing

To launch a test USR Canet 200 server, run in the root of the repository:

    ./utils/replay_usr_canet200_data_to_tcp.py test_data/can_data_usr_canet200_3.json

For development purposes, it is possible to use the vcan kernel module:

    sudo modprobe vcan
    sudo ip link add dev vcan0 type vcan

In order to start developing, run the following in the root of the repository:

    cd app
    go mod download

Then, set up the config.yaml file

    cp ./config.yaml.template ./config.yaml

Then, to run the app:

    go run cmd/usr_canet_to_can_converter/usr_canet_to_can_converter.go

To visualize the output on the simulated vcan0 port:

    candump vcan0

To install the dependencies for the admin interface, run the following in the root of the repository:

    cd admin
    yarn install

Then, to start the admin interface:

    yarn start

The admin should be visible at http://localhost:3000

### Create a .deb installer

To create a .deb installer for the project, make sure you have docker installed.

First, build the image:

    cd build_system/docker
    ./build.sh

To create the .deb:

    cd build_system/
    ./build_deb_amd64.sh /path/to/desired/output/folder # Replace with desired output folder path

In the case of this installer, the admin interface is running at port 9402.
