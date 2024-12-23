// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Cypodsolutions",
            "url": "http://www.cypod.com/",
            "email": "naders@cypodsolutions.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/categories": {
            "get": {
                "description": "get all categories names, pagination options available",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Categories"
                ],
                "summary": "Get all categories",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page ID",
                        "name": "page_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page Size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Category"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new category specifying its name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Categories"
                ],
                "summary": "Create new category",
                "parameters": [
                    {
                        "description": "Create category body",
                        "name": "arg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.createCategorytRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Category"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/tickets": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "List all tickets for a specific user, Admin can get all tickets and can add query param to filter by category ID, assigned engineer and ticket owner normal user only can only get all tickets assigned to him",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "List tickets",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter Ticket owner",
                        "name": "user_assigned",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page ID",
                        "name": "page_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page Size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Filter Category ID",
                        "name": "category_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter Assigned engineer",
                        "name": "assigned_to",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Is admin",
                        "name": "is_admin",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User sending the request",
                        "name": "requester",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Ticket"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "post": {
                "description": "Create a support ticket for an end user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Create ticket",
                "parameters": [
                    {
                        "description": "Create Ticket body",
                        "name": "arg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.createTicketRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Ticket"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/tickets/{ticket_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Admins get any ticket, normal user only get a ticket he owns",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Get ticket by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Ticket ID",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Is admin",
                        "name": "is_admin",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User sending the request",
                        "name": "requester",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Ticket"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update ticket by a ticket ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Update ticket",
                "parameters": [
                    {
                        "description": "Update ticket body",
                        "name": "arg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.updateTicketRequestJSON"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "ticket ID for update",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Ticket"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Delete ticket by a ticket ID",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Delete ticket",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Ticket ID",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "true"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/tickets/{ticket_id}/comments": {
            "get": {
                "description": "List all comments from a ticket",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "List comments",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Ticket ID",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page ID",
                        "name": "page_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page Size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Comment"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "post": {
                "description": "Add a new comment to a ticket",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Add comment",
                "parameters": [
                    {
                        "description": "Create comment body",
                        "name": "arg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.createCommentRequestJSON"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Ticket ID",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Comment"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/tickets/{ticket_id}/comments/{comment_id}": {
            "put": {
                "description": "Update a comment from a ticket",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Update comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Comment ID",
                        "name": "comment_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Ticket ID",
                        "name": "ticket_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Comment text",
                        "name": "arg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.updateCommentRequestJSON"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Comment"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Delete a comment from a ticket",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Comments"
                ],
                "summary": "Delete comment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Comment ID",
                        "name": "comment_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "true"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "api.createCategorytRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "minLength": 1
                }
            }
        },
        "api.createCommentRequestJSON": {
            "type": "object",
            "required": [
                "comment_text",
                "user_commented"
            ],
            "properties": {
                "comment_text": {
                    "type": "string",
                    "minLength": 1
                },
                "user_commented": {
                    "type": "string",
                    "minLength": 1
                }
            }
        },
        "api.createTicketRequest": {
            "type": "object",
            "required": [
                "category_id",
                "description",
                "status",
                "title",
                "user_assigned"
            ],
            "properties": {
                "category_id": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "inprogress",
                        "closed",
                        "open"
                    ]
                },
                "title": {
                    "type": "string"
                },
                "user_assigned": {
                    "type": "string"
                }
            }
        },
        "api.updateCommentRequestJSON": {
            "type": "object",
            "required": [
                "comment_text"
            ],
            "properties": {
                "comment_text": {
                    "type": "string",
                    "minLength": 0
                }
            }
        },
        "api.updateTicketRequestJSON": {
            "type": "object",
            "required": [
                "status"
            ],
            "properties": {
                "assigned_to": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "inprogress",
                        "closed",
                        "open"
                    ]
                }
            }
        },
        "db.Category": {
            "type": "object",
            "properties": {
                "category_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "db.Comment": {
            "type": "object",
            "properties": {
                "comment_id": {
                    "type": "integer"
                },
                "comment_text": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "ticket_id": {
                    "type": "integer"
                },
                "user_commented": {
                    "type": "string"
                }
            }
        },
        "db.Ticket": {
            "type": "object",
            "properties": {
                "assigned_to": {
                    "$ref": "#/definitions/sql.NullString"
                },
                "category_id": {
                    "type": "integer"
                },
                "closed_at": {
                    "$ref": "#/definitions/sql.NullTime"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "ticket_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_assigned": {
                    "type": "string"
                }
            }
        },
        "sql.NullString": {
            "type": "object",
            "properties": {
                "string": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if String is not NULL",
                    "type": "boolean"
                }
            }
        },
        "sql.NullTime": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Description for what is this security definition being used",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Gin Swagger Example API",
	Description:      "Ticketing support microservice",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
