#!/usr/bin/python

from fabric.api import hide, env, settings, abort, run, cd, shell_env
from fabric.colors import magenta, red
from fabric.contrib.files import append
from fabric.contrib.project import rsync_project
import os

env.user = 'root'
env.abort_on_prompts = True
PATH = '/home/dbeliakov/docker/revisor/'

def deploy():
    def rsync(filename):
        rsync_project(PATH, filename, delete=True)

    def generate_conf():
        with open('config.env', 'w') as f:
            f.write('DEBUG=0\n')
            f.write('SECRET_KEY=' + os.getenv('SECRET_KEY') + '\n')
            f.write('CLIENT_FILES_DIR=./client')
            f.write('DATABASE_FILE=/revisor.db')

    def docker_compose(command):
        with cd(PATH):
            run('set -o pipefail; docker-compose %s | tee' % command)

    run('docker login -u %s -p %s %s' % (os.getenv('REGISTRY_USER',
                                                   'gitlab-ci-token'),
                                         os.getenv('CI_BUILD_TOKEN'),
                                         os.getenv('CI_REGISTRY',
                                                   'registry.gitlab.dbeliakov.ru')))

    rsync('docker-compose.yml')
    generate_conf()
    rsync('config.env')
    docker_compose('pull')
    docker_compose('up -d')

