# Terraformed Inventory

This is a little Go app which generates an dynamic [Ansible] [ansible] inventory
from a [Terraform] [tf] state file. It allows one to spawn a bunch of VMs with
Terraform, then (re-)provision them with Ansible. It's pretty neat.


## Usage

Just download the OSX binary and run it.

    curl https://github.com/adammck/terraformedinventory/releases/download/v0.1/terraformedinventory
    ./terraformedinventory --list whatever.tfstate

Ansible doesn't (seem to) support calling the inventory script with parameters,
so I like to wrap this tool up in a little shell script, and call that.
Something like:

	#!/bin/bash
	terraformedinventory $@ $(dirname $0)/deploy/terraform.tfstate

Configuration via environment variables, like most other dynamic inventory
scripts, is coming soon.


## Development

[Install Terraform] [tfdev] from source, then:

	git clone https://github.com/adammck/terraformedinventory.git
	cd terraformedinventory
	go build


## License

MIT.




[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
[tfdev]:   https://github.com/hashicorp/terraform#developing-terraform
