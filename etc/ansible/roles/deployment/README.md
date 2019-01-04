Deploy Modern Go Application
============================

This role is used to deploy Modern Go Application from a CI/CD pipeline.

Requirements
------------

`modern-go-application` role should be applied to the host (either in a separate or in the same playbook).

Role Variables
--------------

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `binary_source` | *none* | Local source of the binaries to copy |
| `binary_name` | `modern-go-application` | Binary to copy |
| `mga_service_name` | `mga` | Service to be restarted |

Dependencies
------------

- `modern-go-application` role

Example Playbook
----------------

    - hosts: servers
      roles:
         - { role: deploy-modern-go-application, binary_source: build/ }

License
-------

MIT
