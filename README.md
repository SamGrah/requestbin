## Description
This is a simple tool for capturing and reviewing requests made by a webhook.
The user creates a 'bin' and provides two unique endpoints to the user.

All requests made to one of the unique endpoints will be captured in a MongoDB instance.
The other unique endpoint provides a page on which the user can review all of those captured request.

It's main use is to test that webhooks are configured and working as expected.


## Tech Stack
Frontend  - SvelteKit
Go & gRPC - Service Code
MongoDB   - Requests Datastore
Docker    - Deployment


## Instructions
First, build all necessary docker images. Docker must be installed an running in the environment.
```
> make up_build
```

Run all services (including the frontend)
```
> make up
```

Teardown and destroy all containers
```
> make down
```
