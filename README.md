# Terraformed Inventory

This is a little Go app which generates an dynamic [Ansible] [ansible] inventory
from a [Terraform] [tf] state file. It allows one to spawn a bunch of instances
with Terraform, then (re-)provision them with Ansible. It's pretty neat. 

Currently, only **AWS**, **DigitalOcean**, **CloudStack** and **VMware** are supported.


# Installation

On OSX, install it with Homebrew:

	brew install https://raw.github.com/adammck/terraform-inventory/master/homebrew/terraform-inventory.rb

This is only a tiny tool, so it's not in the main Homebrew repo. Feel free to
add it, if you think that would be useful.


## Usage

If your Terraform state file is named `terraform.tfstate` (the default), `cd` to
it and run:

	ansible-playbook --inventory-file=terraform-inventory deploy/playbook.yml

This will provide the resource names and IP addresses of any instances found in
the state file to Ansible, which can then be used as hosts patterns in your
playbooks. For example, given for the following Terraform config:

	resource "digitalocean_droplet" "my-web-server" {
	  image = "centos-7-0-x64"
	  name = "web-1"
	  region = "nyc1"
	  size = "512mb"
	}

The corresponding playbook might look like:

	- hosts: my-web-server
	  tasks:
	    - yum: name=cowsay
	    - command: cowsay hello, world!

Note that the instance was identified by its _resource name_ from the Terraform
config, not its _instance name_ from the provider.


## More Usage

Ansible doesn't seem to support calling a dynamic inventory script with params,
so if you need to specify the location of your state file, set the `TF_STATE`
environment variable before running `ansible-playbook`, like:

	TF_STATE=deploy/terraform.tfstate ansible-playbook --inventory-file=terraform-inventory deploy/playbook.yml

Alternately, if you need to do something fancier (like downloading your state
file from S3 before running), you might wrap this tool with a shell script, and
call that instead. Something like:

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
account details, and run `fixtures/update`. You almost certainly don't need to
do this.


## License

MIT.




[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
