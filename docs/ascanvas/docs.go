// Package ascanvas GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package ascanvas

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
        "/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "\"List all canvas items\"",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ascanvas.Canvas"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "\"Create all canvas items\"",
                "parameters": [
                    {
                        "description": "Canvas creation details",
                        "name": "CreateArgs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ascanvas.CreateArgs"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ascanvas.Canvas"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    }
                }
            }
        },
        "/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get a specific canvas by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identifier of canvas to fetch",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ascanvas.Canvas"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "summary": "\"Delete a specific canvas item by id\"",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identifier of canvas to delete",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    }
                }
            }
        },
        "/{id}/events": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/event-stream",
                    "application/json"
                ],
                "summary": "\"Obtain an SSE live stream of canvas events for a specific canvas id\"",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identifier of canvas to observe",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/{id}/floodfill": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "\"Apply flood fill on a specific canvas\"",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identifier of canvas to modify",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Flood fill transformation details",
                        "name": "Transformation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ascanvas.TransformFloodfillArgs"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    },
                    "400": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/{id}/rectangle": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "\"Draw a rectangle on a specific canvas\"",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identifier of canvas to modify",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Rectangle transformation details",
                        "name": "Transformation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ascanvas.TransformRectangleArgs"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.Response"
                        }
                    },
                    "400": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "ascanvas.Canvas": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "height": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "ascanvas.Coordinates": {
            "type": "object",
            "properties": {
                "x": {
                    "type": "integer"
                },
                "y": {
                    "type": "integer"
                }
            }
        },
        "ascanvas.CreateArgs": {
            "type": "object",
            "properties": {
                "fill": {
                    "type": "string"
                },
                "height": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "ascanvas.TransformFloodfillArgs": {
            "type": "object",
            "properties": {
                "fill": {
                    "type": "string"
                },
                "start": {
                    "$ref": "#/definitions/ascanvas.Coordinates"
                }
            }
        },
        "ascanvas.TransformRectangleArgs": {
            "type": "object",
            "properties": {
                "fill": {
                    "type": "string"
                },
                "height": {
                    "type": "integer"
                },
                "outline": {
                    "type": "string"
                },
                "top_left": {
                    "$ref": "#/definitions/ascanvas.Coordinates"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "web.Response": {
            "type": "object",
            "properties": {
                "message": {
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
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
	swag.Register(swag.Name, &s{})
}
