variable "do_token" {}
variable "aws_access_key" {}
variable "aws_secret_key" {}
variable "aws_subnet_id" {}

provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "us-east-1"
}

provider "digitalocean" {
  token = "${var.do_token}"
}

resource "aws_instance" "web-aws" {
    ami = "ami-96a818fe"
    instance_type = "t2.micro"
    subnet_id = "${var.aws_subnet_id}"
    root_block_device = {
      delete_on_termination = true
    }
}

resource "digitalocean_droplet" "web-do" {
  image = "centos-7-0-x64"
  name = "terraform-inventory-1"
  region = "nyc1"
  size = "512mb"
  ssh_keys = [55015]
}
