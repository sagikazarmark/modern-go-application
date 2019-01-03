import os

import testinfra.utils.ansible_runner

testinfra_hosts = testinfra.utils.ansible_runner.AnsibleRunner(
    os.environ['MOLECULE_INVENTORY_FILE']).get_hosts('all')


def test_mga_user_is_created(host):
    user = host.user("mga")

    assert user.shell == "/bin/bash"


def test_mga_user_is_lingering(host):
    f = host.file("/var/lib/systemd/linger/mga")

    assert f.exists


def test_mga_service_unit_is_created(host):
    user = host.user("mga")
    f = host.file(user.home + "/.config/systemd/user/mga.service")

    assert f.exists
    assert f.user == user.name
    assert f.group == user.group
    assert f.mode == 0o600


def test_nginx_is_configured(host):
    f = host.file("/etc/nginx/conf.d/modern-go-application.conf")

    assert f.exists
