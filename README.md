# Terraformed Inventory

This is a little Go app which generates an dynamic [Ansible] [ansible] inventory
from an AWS-based [Terraform] [tf] state file. It allows one to spawn a bunch of EC2 VMs with
Terraform, then (re-)provision them with Ansible. It's pretty neat. 


# Installation

On OSX, install it with Homebrew:

	brew install https://raw.github.com/adammck/terraform-inventory/master/homebrew/terraform-inventory.rb

This is only a tiny tool, so it's not in the main Homebrew repo. Feel free to
add it, if you think that would be useful.


## Usage

Ansible doesn't (seem to) support calling the inventory script with parameters,
so you can specify the path to the state file using the `TI_TFSTATE` environment
variable, like so:

    TI_TFSTATE=deploy/terraform.tfstate ansible-playbook --inventory-file=terraform-inventory

Alternately, you can create a little shell script and call that. Something like:

	#!/bin/bash
	terraform-inventory $@ deploy/terraform.tfstate

Then run Ansible with the script as an inventory:

	ansible-playbook --inventory-file=bin/inventory deploy/playbook.yml


## Development

It's just a Go app, so the usual:

	go get github.com/adammck/terraform-inventory
	cd $GOPATH/adammck/terraform-inventory
	go build


## License

MIT.




[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
