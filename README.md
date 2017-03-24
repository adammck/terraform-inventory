# Terraform Inventory

[![Build Status](https://travis-ci.org/adammck/terraform-inventory.svg?branch=master)](https://travis-ci.org/adammck/terraform-inventory)
[![GitHub release](https://img.shields.io/github/release/adammck/terraform-inventory.svg?maxAge=2592000)](https://github.com/adammck/terraform-inventory/releases)
[![GitHub release](https://img.shields.io/homebrew/v/terraform-inventory.svg?maxAge=2592000)](http://braumeister.org/formula/terraform-inventory)

This is a little Go app which generates a dynamic [Ansible] [ansible] inventory
from a [Terraform] [tf] state file. It allows one to spawn a bunch of instances
with Terraform, then (re-)provision them with Ansible. Currently, only **AWS**,
**DigitalOcean**, **CloudStack**, **VMware**, **OpenStack**, **Google Compute
Engine**, and **SoftLayer** are supported.


# Help Wanted 🙋

This library is stable, but I've been neglecting it somewhat on account of no
longer using Ansible at work. Please drop me a line if you'd be interested in
helping to maintain this tool.


# Installation

On OSX, install it with Homebrew:

	brew install terraform-inventory

Alternatively, you can download a [release][rel] suitable for your platform and
unzip it. Make sure the `terraform-inventory` binary is executable, and you're
ready to go.


## Usage

If you are using [remote state][rs] (or if your state file happens to be named
`terraform.tfstate`), `cd` to it and run:

	ansible-playbook --inventory-file=/path/to/terraform-inventory deploy/playbook.yml

This will provide the resource names and IP addresses of any instances found in
the state file to Ansible, which can then be used as hosts patterns in your
playbooks. For example, given for the following Terraform config:

	resource "digitalocean_droplet" "my_web_server" {
	  image = "centos-7-0-x64"
	  name = "web-1"
	  region = "nyc1"
	  size = "512mb"
	}

The corresponding playbook might look like:

	- hosts: my_web_server
	  tasks:
	    - yum: name=cowsay
	    - command: cowsay hello, world!

Note that the instance was identified by its _resource name_ from the Terraform
config, not its _instance name_ from the provider. On AWS, resources are also
grouped by their tags. For example:

	resource "aws_instance" "my_web_server" {
	  instance_type = "t2.micro"
	  ami = "ami-96a818fe"
	  tags = {
	    Role = "web"
	    Env = "dev"
	  }
	}

	resource "aws_instance" "my_worker" {
	  instance_type = "t2.micro"
	  ami = "ami-96a818fe"
	  tags = {
	    Role = "worker"
	    Env = "dev"
	  }
	}

Can be provisioned separately with:

	- hosts: role_web
	  tasks:
	    - command: cowsay this is a web server!

	- hosts: role_worker
	  tasks:
	    - command: cowsay this is a worker server!

	- hosts: env_dev
	  tasks:
	    - command: cowsay this runs on all dev servers!


## More Usage

Ansible doesn't seem to support calling a dynamic inventory script with params,
so if you need to specify the location of your state file or terraform directory, set the `TF_STATE`
environment variable before running `ansible-playbook`, like:


	TF_STATE=deploy/terraform.tfstate ansible-playbook --inventory-file=/path/to/terraform-inventory deploy/playbook.yml

	or

	TF_STATE=../terraform ansible-playbook --inventory-file=/path/to/terraform-inventory deploy/playbook.yml

If `TF_STATE` is a file, it parses the file as json, if `TF_STATE` is a directory, it runs `terraform state pull` inside the directory, which is supports both local and remote terraform state.

It looks for state config in this order

- `TF_STATE`: environment variable of where to find either a statefile or a terraform project
- `TI_TFSTATE`: another environment variable similar to TF_STATE
- `terraform.tfstate`: it looks in the state file in the current directory.
- `.`: lastly it assumes you are at the root of a terraform project.

Alternately, if you need to do something fancier (like downloading your state
file from S3 before running), you might wrap this tool with a shell script, and
call that instead. Something like:

	#!/bin/bash
	/path/to/terraform-inventory $@ deploy/terraform.tfstate

Then run Ansible with the script as an inventory:

	ansible-playbook --inventory-file=bin/inventory deploy/playbook.yml

This tool returns the public IP of the host by default. If you require the private
IP of the instance to run Ansible, set the `TF_KEY_NAME` environment variable
to `private_ip` before running the playbook, like:

	TF_KEY_NAME=private_ip ansible-playbook --inventory-file=/path/to/terraform-inventory deploy/playbook.yml

## Development

It's just a Go app, so the usual:

	go get github.com/adammck/terraform-inventory

To test against an example statefile, run:

	terraform-inventory --list fixtures/example.tfstate
	terraform-inventory --host=52.7.58.202 fixtures/example.tfstate

To update the fixtures, populate `fixtures/secrets.tfvars` with your DO and AWS
account details, and run `fixtures/update`. To run a tiny Ansible playbook on
the example resourecs, run:

	TF_STATE=fixtures/example.tfstate ansible-playbook --inventory-file=/path/to/terraform-inventory fixtures/playbook.yml

You almost certainly don't need to do any of this. Use the tests instead.


## Acknowledgements

Development of
[#14](https://github.com/adammck/terraform-inventory/issues/14),
[#16](https://github.com/adammck/terraform-inventory/issues/16),
and [#22](https://github.com/adammck/terraform-inventory/issues/22)
was generously sponsored by [Transloadit](https://transloadit.com).


## License

MIT.

[ansible]: https://www.ansible.com
[tf]:      https://www.terraform.io
[rel]:     https://github.com/adammck/terraform-inventory/releases
[rs]:      https://www.terraform.io/docs/state/remote/index.html
