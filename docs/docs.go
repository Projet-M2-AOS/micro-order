// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/orders": {
            "get": {
                "description": "Return all orders.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "micro-orders"
                ],
                "summary": "Get all orders.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search by userid",
                        "name": "user",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Order"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new order with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "micro-orders"
                ],
                "summary": "Create a new order",
                "parameters": [
                    {
                        "description": "Create order",
                        "name": "Order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Order"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Order"
                            }
                        }
                    }
                }
            }
        },
        "/orders/{id}": {
            "get": {
                "description": "Return one orders.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "micro-orders"
                ],
                "summary": "Get one orders.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Order"
                        }
                    }
                }
            },
            "put": {
                "description": "Update order with the input payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "micro-orders"
                ],
                "summary": "Update order",
                "parameters": [
                    {
                        "description": "Update order",
                        "name": "Order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Order"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "Delete order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "micro-orders"
                ],
                "summary": "Delete order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Order": {
            "type": "object",
            "required": [
                "address",
                "paymentState",
                "price",
                "products",
                "user"
            ],
            "properties": {
                "_id": {
                    "type": "string"
                },
                "address": {
                    "type": "string"
                },
                "paymentState": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "user": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Order micro-service",
	Description: "Order micro-service documentation.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
