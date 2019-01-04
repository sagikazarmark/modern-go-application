Application Setup
=================

Prepares a host for running the application.
Normally this role is part of a bigger playbook executed outside of the application.
Alternatively this role can become a whole playbook as running an application might require configuration
across various hosts (eg. setting up database users, firewall rules, etc)   .

Requirements
------------

Requires nginx to be installed.

Example Playbook
----------------

    - hosts: servers
      roles:
         - { role: application-setup }

License
-------

MIT
