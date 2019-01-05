import os

import testinfra.utils.ansible_runner

testinfra_hosts = testinfra.utils.ansible_runner.AnsibleRunner(
    os.environ['MOLECULE_INVENTORY_FILE']).get_hosts('all')


def test_app_user_is_created(host):
    user = host.user("app")

    assert user.shell == "/bin/bash"
    assert "docker" in user.groups


def test_nginx_is_configured(host):
    f = host.file("/etc/nginx/conf.d/application.local.conf")

    assert f.exists
