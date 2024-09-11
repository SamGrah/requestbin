## Description

This is a simple tool for capturing and reviewing requests made by a webhook.
The user creates a 'bin' and the app provides two unique endpoints to the user.

All requests made to one of the unique endpoints will be captured in a database.
The other unique endpoint provides a page on which the user can review all of those captured requests.

It's main use is to test that webhooks are configured and working as expected.

## Implementations

There are multiple implementations of this tool, each with its own dedicated branch.

#### MVC-Go-HTMX

> Current Working Branch
The `MVC-Go-HTMX` branch contains a model-view-controller implementation of the
tool. The API server is written in Go and the UI is defined using the Go standard
lib HTML templating in combination with HTMX.

#### MVC-Node

> Archived
The `MVC-Node` branch contains a model-view-controller implementation of the tool.
The API server is written in TS Node and uses Pug to define templates for the served pages.

#### SPA

> Archived
The `SPA` branch contains a single page application implementaion of the tool.
The frontend is written in Svelte and the backend is a single TS Node service.

#### Microservices

> WIP branch
The `Microservices` branch contains a microservices implementation for the backend.
The is composed of multiple services including: api-gateway (Go), bin manager (Go), and pending webhook service (TS Node).
The frontend is implmented in SvelteKit.

## Purpose

This app is my playground. When I want to learn new tooling or a language (or whatever),
I dive into this repo and start mucking about.
It's fullstack, which I kinda enjoy staying on-top of, and satisfies
[the knack](https://www.youtube.com/watch?v=g8vHhgh6oM0) when it arises.
