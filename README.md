# router

## Module

- server service
- agent service
- client

## Dependencies

- [GopherJS](https://github.com/gopherjs/gopherjs) (Go to JavaScript transpiler)

```
go get -u github.com/gopherjs/gopherjs
```

## model

### user

- usename
- password
- detai
- status

### worker

- id
- key
- name
- desc
- status

### task

- id
- user-id
- worker-id
- params
- status
- progress
- detail

