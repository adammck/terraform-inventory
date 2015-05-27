package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const exampleStateFile = `
{
    "version": 1,
    "serial": 1,
    "modules": [
        {
            "path": [
                "root"
            ],
            "outputs": {},
            "resources": {
                "aws_instance.one": {
                    "type": "aws_instance",
                    "primary": {
                        "id": "i-aaaaaaaa",
                        "attributes": {
                            "ami": "ami-XXXXXXXX",
                            "id": "i-aaaaaaaa"
                        }
                    }
                },
                "aws_instance.two": {
                    "type": "aws_instance",
                    "primary": {
                        "id": "i-bbbbbbbb",
                        "attributes": {
                            "ami": "ami-YYYYYYYY",
                            "id": "i-bbbbbbbb"
                        }
                    }
                },
                "aws_security_group.example": {
                    "type": "aws_security_group",
                    "primary": {
                        "id": "sg-cccccccc",
                        "attributes": {
                            "description": "Whatever",
                            "id": "sg-cccccccc"
                        }
                    }
                }
            }
        }
    ]
}
`

func TestStateRead(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	exp := state{
		Modules: []moduleState{
			moduleState{
				Resources: map[string]resourceState{
					"aws_instance.one": resourceState{
						Type: "aws_instance",
						Primary: instanceState{
							ID: "i-aaaaaaaa",
							Attributes: map[string]string{
								"ami": "ami-XXXXXXXX",
								"id":  "i-aaaaaaaa",
							},
						},
					},
					"aws_instance.two": resourceState{
						Type: "aws_instance",
						Primary: instanceState{
							ID: "i-bbbbbbbb",
							Attributes: map[string]string{
								"ami": "ami-YYYYYYYY",
								"id":  "i-bbbbbbbb",
							},
						},
					},
					"aws_security_group.example": resourceState{
						Type: "aws_security_group",
						Primary: instanceState{
							ID: "sg-cccccccc",
							Attributes: map[string]string{
								"description": "Whatever",
								"id":          "sg-cccccccc",
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, exp, s)
}

func TestInstances(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.instances()
	assert.Equal(t, 2, len(inst))
	assert.Equal(t, "i-aaaaaaaa", inst["one"].ID)
	assert.Equal(t, "i-bbbbbbbb", inst["two"].ID)
}
