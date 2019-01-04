import os

import testinfra.utils.ansible_runner

testinfra_hosts = testinfra.utils.ansible_runner.AnsibleRunner(
    os.environ['MOLECULE_INVENTORY_FILE']).get_hosts('all')


def test_app_user_is_created(host):
    user = host.user("mga")

    assert user.shell == "/bin/bash"
    assert "docker" in user.groups


def test_app_user_facts_are_registered(host):
    f = host.file("/etc/ansible/facts.d/app_user.fact")

    assert f.exists
    assert f.mode == 0o644
    assert f.contains('"name": "mga"')
    assert f.contains('"home": "/home/mga"')
    assert f.contains('"uid": "1000"')


def test_nginx_is_configured(host):
    f = host.file("/etc/nginx/conf.d/modern-go-application.conf")

    assert f.exists
