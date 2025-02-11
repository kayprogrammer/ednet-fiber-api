// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
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
        "/auth/google": {
            "post": {
                "description": "` + "`" + `This endpoint generates new access and refresh tokens for authentication via google` + "`" + `\n` + "`" + `Pass in token gotten from gsi client authentication here in payload to retrieve tokens for authorization` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user via google",
                "parameters": [
                    {
                        "description": "Google auth",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.TokenSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/accounts.LoginResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/base.UnauthorizedErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "` + "`" + `This endpoint generates new access and refresh tokens for authentication` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.LoginSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/accounts.LoginResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/base.UnauthorizedErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "` + "`" + `This endpoint logs a user out from our application from a single device` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Logout a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/base.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/base.UnauthorizedErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/logout/all": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "` + "`" + `This endpoint logs a user out from our application from all devices` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Logout a user from all devices",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/base.ResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/base.UnauthorizedErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "` + "`" + `This endpoint refresh tokens by generating new access and refresh tokens for a user` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Refresh tokens",
                "parameters": [
                    {
                        "description": "Refresh token",
                        "name": "refresh",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.TokenSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/accounts.LoginResponseSchema"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/base.UnauthorizedErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "` + "`" + `This endpoint registers new users into our application.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.RegisterSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/accounts.RegisterResponseSchema"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/resend-verification-email": {
            "post": {
                "description": "` + "`" + `This endpoint resends new otp to the user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Resend Verification Email",
                "parameters": [
                    {
                        "description": "Email object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.EmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/base.ResponseSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/base.NotFoundErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/send-password-reset-otp": {
            "post": {
                "description": "` + "`" + `This endpoint sends new password reset otp to the user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Send Password Reset Otp",
                "parameters": [
                    {
                        "description": "Email object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.EmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/base.ResponseSchema"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/base.NotFoundErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/set-new-password": {
            "post": {
                "description": "` + "`" + `This endpoint verifies the password reset otp.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Set New Password",
                "parameters": [
                    {
                        "description": "Password reset object",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.SetNewPasswordSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/base.ResponseSchema"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/base.InvalidErrorExample"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/base.NotFoundErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/auth/verify-email": {
            "post": {
                "description": "` + "`" + `This endpoint verifies a user's email.` + "`" + `",
                "tags": [
                    "Auth"
                ],
                "summary": "Verify a user's email",
                "parameters": [
                    {
                        "description": "Email object",
                        "name": "email_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/accounts.VerifyEmailRequestSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/base.ResponseSchema"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/base.InvalidErrorExample"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/base.NotFoundErrorExample"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/base.ValidationErrorExample"
                        }
                    }
                }
            }
        },
        "/general/site-detail": {
            "get": {
                "description": "This endpoint retrieves few details of the site/application.",
                "tags": [
                    "General"
                ],
                "summary": "Retrieve site details",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/general.SiteDetailResponseSchema"
                        }
                    }
                }
            }
        },
        "/healthcheck": {
            "get": {
                "description": "This endpoint checks the health of our application.",
                "tags": [
                    "HealthCheck"
                ],
                "summary": "HealthCheck",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.HealthCheckSchema"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "accounts.EmailRequestSchema": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                }
            }
        },
        "accounts.LoginResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/accounts.TokensResponseSchema"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "accounts.LoginSchema": {
            "type": "object",
            "required": [
                "email_or_username",
                "password"
            ],
            "properties": {
                "email_or_username": {
                    "type": "string",
                    "example": "johndoe"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "accounts.RegisterResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/accounts.EmailRequestSchema"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "accounts.RegisterSchema": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@example.com"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "John Doe"
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "strongpassword"
                },
                "username": {
                    "type": "string",
                    "maxLength": 50,
                    "example": "johndoe"
                }
            }
        },
        "accounts.SetNewPasswordSchema": {
            "type": "object",
            "required": [
                "email",
                "otp",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "otp": {
                    "type": "integer",
                    "example": 123456
                },
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8,
                    "example": "newstrongpassword"
                }
            }
        },
        "accounts.TokenSchema": {
            "type": "object",
            "required": [
                "token"
            ],
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"
                }
            }
        },
        "accounts.TokensResponseSchema": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                },
                "refresh": {
                    "type": "string"
                }
            }
        },
        "accounts.VerifyEmailRequestSchema": {
            "type": "object",
            "required": [
                "email",
                "otp"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 5,
                    "example": "johndoe@email.com"
                },
                "otp": {
                    "type": "integer",
                    "example": 123456
                }
            }
        },
        "base.FieldData": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string",
                    "example": "This field is required"
                }
            }
        },
        "base.InvalidErrorExample": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Request was invalid due to ..."
                },
                "status": {
                    "type": "string",
                    "example": "failure"
                }
            }
        },
        "base.NotFoundErrorExample": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "The item was not found"
                },
                "status": {
                    "type": "string",
                    "example": "failure"
                }
            }
        },
        "base.ResponseSchema": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "base.UnauthorizedErrorExample": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Unauthorized user/Invalid credentials/Invalid Token"
                },
                "status": {
                    "type": "string",
                    "example": "failure"
                }
            }
        },
        "base.ValidationErrorExample": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/base.FieldData"
                },
                "message": {
                    "type": "string",
                    "example": "Invalid Entry"
                },
                "status": {
                    "type": "string",
                    "example": "failure"
                }
            }
        },
        "general.SiteDetailResponseSchema": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/general.SiteDetailSchema"
                },
                "message": {
                    "type": "string",
                    "example": "Data fetched/created/updated/deleted"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "general.SiteDetailSchema": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "example": "234, Lagos, Nigeria"
                },
                "email": {
                    "type": "string",
                    "example": "johndoe@email.com"
                },
                "fb": {
                    "type": "string",
                    "example": "https://facebook.com"
                },
                "ig": {
                    "type": "string",
                    "example": "https://instagram.com"
                },
                "name": {
                    "type": "string",
                    "example": "EDNET"
                },
                "phone": {
                    "type": "string",
                    "example": "+2348133831036"
                },
                "tw": {
                    "type": "string",
                    "example": "https://twitter.com"
                },
                "wh": {
                    "type": "string",
                    "example": "https://wa.me/2348133831036"
                }
            }
        },
        "routes.HealthCheckSchema": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "string",
                    "example": "pong"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type 'Bearer jwt_string' to correctly set the API Key",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "security": [
        {
            "BearerAuth": []
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "EDNET API",
	Description:      "## A Full-Featured EDTECH API built with FIBER & ENT ORM.\n\n<!-- ### WEBSOCKETS:\n\n#### Notifications\n\n- URL: `wss://{host}/api/v4/ws/notifications`\n\n- Requires authorization, so pass in the Bearer Authorization header.\n\n- You can only read and not send notification messages into this socket. -->\n\n\n<!-- #### Chats\n\n- URL: `wss://{host}/api/v4/ws/chats/{id}`\n- Requires authorization, so pass in the Bearer Authorization header.\n- Use chat_id as the ID for an existing chat or username if it's the first message in a DM.\n- You cannot read realtime messages from a username that doesn't belong to the authorized user, but you can surely send messages.\n- Only send a message to the socket endpoint after the message has been created or updated, and files have been uploaded.\n- Fields when sending a message through the socket:\n\n  ```json\n  { \"status\": \"CREATED\", \"id\": \"fe4e0235-80fc-4c94-b15e-3da63226f8ab\" }\n  ``` -->",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
