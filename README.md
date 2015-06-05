# Terraformed Inventory

This is a little Go app which generates an dynamic [Ansible] [ansible] inventory
from a [Terraform] [tf] state file. It allows one to spawn a bunch of VMs with
Terraform, then (re-)provision them with Ansible. It's pretty neat. 

Currently, only **AWS** and **DigitalOcean** are supported.


# Installation

On OSX, install it with Homebrew:

	brew install https://raw.github.com/adammck/terraform-inventory/master/homebrew/terraform-inventory.rb

This is only a tiny tool, so it's not in the main Homebrew repo. Feel free to
add it, if you think that would be useful.


## Usage

Ansible doesn't (seem to) support calling the inventory script with parameters,
so you can specify the path to the state file using the `TF_STATE` environment
variable, like so:

	TF_STATE=deploy/terraform.tfstate ansible-playbook --inventory-file=terraform-inventory deploy/playbook.yml

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

To test against an example statefile, run:

	terraform-inventory --list fixtures/example.tfstate
	terraform-inventory --host=web-aws fixtures/example.tfstate

To update the fixtures, populate `fixtures/secrets.tfvars` with your DO and AWS
account details, and run `fixtures/update`. You probably don't need to do this.


## License

MIT.




[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
