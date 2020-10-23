// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/chat/:id": {
            "post": {
                "description": "This endpoint sends a message to the given thread",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Sends a message to a topic",
                "operationId": "sendMessage",
                "parameters": [
                    {
                        "description": "Message to send",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.SendMessageRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Topic id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/chat/messages": {
            "get": {
                "description": "This endpoint returns the messages for the given topic.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Gets topic messages",
                "operationId": "getMessages",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of messages to take",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Number of messages to skip",
                        "name": "skip",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Topic id",
                        "name": "topic",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.GetTopicMessagesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/chat/threads": {
            "get": {
                "description": "This endpoint returns the latest messaging threads for the currently logged in user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Returns the latest user message threads",
                "operationId": "getLatestThreads",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of threads to take",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Number of threads to skip",
                        "name": "skip",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.GetLatestThreadsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/meta/who-am-i": {
            "get": {
                "description": "Returns information about the currently authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Returns information about myself",
                "operationId": "whoAmI",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.UserAuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/resources": {
            "get": {
                "description": "Search for resources",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resources"
                ],
                "summary": "Searches resources",
                "operationId": "searchResources",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search text",
                        "name": "query",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "0",
                            "1"
                        ],
                        "type": "string",
                        "description": "Resource type",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Created by",
                        "name": "created_by",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of resources to take",
                        "name": "take",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Number of resources to skip",
                        "name": "skip",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.SearchResourcesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates a resource",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resources"
                ],
                "summary": "Updates a resource",
                "operationId": "updateResource",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Resource id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Resource to create",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.UpdateResourceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.UpdateResourceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a resource",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resources"
                ],
                "summary": "Creates a resource",
                "operationId": "createResource",
                "parameters": [
                    {
                        "description": "Resource to create",
                        "name": "resource",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.CreateResourceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.CreateResourceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/resources/:id": {
            "get": {
                "description": "Gets a resource by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "resources"
                ],
                "summary": "Gets a single resource",
                "operationId": "getResource",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "Resource id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.GetResourceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/resources/:id/inquire": {
            "post": {
                "description": "This endpoint sends a message to the resource owner",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "resources"
                ],
                "summary": "Sends a message to the user about a resource",
                "operationId": "inquireAboutResource",
                "parameters": [
                    {
                        "description": "Message to send",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.InquireAboutResourceRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Resource id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {},
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        },
        "/users/:id": {
            "get": {
                "description": "Returns information about the given user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Returns information about a user",
                "operationId": "getUserInfo",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "description": "User id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.UserInfoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "statusCode": {
                    "type": "integer"
                }
            }
        },
        "utils.Error": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "web.CreateResourcePayload": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                },
                "valueInHoursFrom": {
                    "type": "integer"
                },
                "valueInHoursTo": {
                    "type": "integer"
                }
            }
        },
        "web.CreateResourceRequest": {
            "type": "object",
            "properties": {
                "resource": {
                    "type": "object",
                    "$ref": "#/definitions/web.CreateResourcePayload"
                }
            }
        },
        "web.CreateResourceResponse": {
            "type": "object",
            "properties": {
                "resource": {
                    "type": "object",
                    "$ref": "#/definitions/web.Resource"
                }
            }
        },
        "web.GetLatestThreadsResponse": {
            "type": "object",
            "properties": {
                "threads": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Thread"
                    }
                }
            }
        },
        "web.GetResourceResponse": {
            "type": "object",
            "properties": {
                "resource": {
                    "type": "object",
                    "$ref": "#/definitions/web.Resource"
                }
            }
        },
        "web.GetTopicMessagesResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Message"
                    }
                }
            }
        },
        "web.InquireAboutResourceRequest": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "web.Message": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "sentAt": {
                    "type": "string"
                },
                "sentBy": {
                    "type": "string"
                },
                "sentByUsername": {
                    "type": "string"
                },
                "topicId": {
                    "type": "string"
                }
            }
        },
        "web.Resource": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "createdBy": {
                    "type": "string"
                },
                "createdById": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                },
                "valueInHoursFrom": {
                    "type": "integer"
                },
                "valueInHoursTo": {
                    "type": "integer"
                }
            }
        },
        "web.SearchResourcesResponse": {
            "type": "object",
            "properties": {
                "resources": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.Resource"
                    }
                },
                "skip": {
                    "type": "integer"
                },
                "take": {
                    "type": "integer"
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "web.SendMessageRequest": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "web.Thread": {
            "type": "object",
            "properties": {
                "hasUnreadMessages": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "lastChars": {
                    "type": "string"
                },
                "lastMessageAt": {
                    "type": "string"
                },
                "lastMessageUserId": {
                    "type": "string"
                },
                "lastMessageUsername": {
                    "type": "string"
                },
                "recipientId": {
                    "type": "string"
                }
            }
        },
        "web.UpdateResourcePayload": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                },
                "valueInHoursFrom": {
                    "type": "integer"
                },
                "valueInHoursTo": {
                    "type": "integer"
                }
            }
        },
        "web.UpdateResourceRequest": {
            "type": "object",
            "properties": {
                "resource": {
                    "type": "object",
                    "$ref": "#/definitions/web.UpdateResourcePayload"
                }
            }
        },
        "web.UpdateResourceResponse": {
            "type": "object",
            "properties": {
                "resource": {
                    "type": "object",
                    "$ref": "#/definitions/web.Resource"
                }
            }
        },
        "web.UserAuthResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "isAuthenticated": {
                    "type": "boolean"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "web.UserInfoResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "username": {
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
	Host:        "127.0.0.1:8585",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "commonpool api",
	Description: "resources api",
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
