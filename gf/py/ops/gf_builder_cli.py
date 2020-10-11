# GloFlow application and media management/publishing platform
# Copyright (C) 2020 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import argparse
import subprocess
import urllib.parse
import json
import requests

from colored import fg, bg, attr
import delegator

sys.path.append("%s/../utils"%(modd_str))
import gf_core_cli

#--------------------------------------------------
def main():

	args_map = parse_args()

	service_cont_image_tag_str  = "latest"
	docker_user_str             = "glofloworg"
	service_cont_image_name_str = f"{docker_user_str}/gf_go_ethereum:{service_cont_image_tag_str}"
	service_cont_dockerfile_path_str = f"{modd_str}/../../../Dockerfile"
	
	
	#------------------------
	# BUILD_CONTAINER
	if args_map["run"] == "build_containers":
		
		build_containers(service_cont_image_name_str,
			service_cont_dockerfile_path_str,
			p_docker_sudo_bool=args_map["docker_sudo_bool"])

	#------------------------
	# PUBLISH_CONTAINER
	elif args_map["run"] == "publish_containers":
		docker_pass_str = args_map["gf_docker_pass_str"]
		assert not docker_pass_str == None

		publish_containers(service_cont_image_name_str,
			docker_user_str,
			docker_pass_str,
			p_docker_sudo_bool=args_map["docker_sudo_bool"])

	#------------------------
	# NOTIFY_COMPLETION
	elif args_map["run"] == "notify_completion":

		gf_notify_completion_url_str = args_map["gf_notify_completion_url_str"]
		assert not gf_notify_completion_url_str == None

		# GIT_COMMIT_HASH
		git_commit_hash_str = None
		if "DRONE_COMMIT" in os.environ.keys():
			git_commit_hash_str = os.environ["DRONE_COMMIT"]

		notify_completion(gf_notify_completion_url_str,
			p_git_commit_hash_str = git_commit_hash_str)

	#------------------------

#--------------------------------------------------
# NOTIFY_COMPLETION
def notify_completion(p_gf_notify_completion_url_str,
	p_git_commit_hash_str = None):
	
	url_str = None

	# add git_commit_hash as a querystring argument to the notify_completion URL.
	# the entity thats receiving the completion notification needs to know what the tag
	# is of the newly created container.
	if not p_git_commit_hash_str == None:
		url = urllib.parse.urlparse(p_gf_notify_completion_url_str)
		
		# QUERY_STRING
		qs_lst = urllib.parse.parse_qsl(url.query)
		qs_lst.append(("git_commit", p_git_commit_hash_str)) # .parse_qs() places all values in lists

		qs_str = "&".join(["%s=%s"%(k, v) for k, v in qs_lst])

		# _replace() - "url" is of type ParseResult which is a subclass of namedtuple;
		#              _replace is a namedtuple method that:
		#              "returns a new instance of the named tuple replacing
		#              specified fields with new values".
		url_new = url._replace(query=qs_str)
		url_str = url_new.geturl()
	else:
		url_str = p_gf_notify_completion_url_str

	print("NOTIFY_COMPLETION - HTTP REQUEST - %s"%(url_str))

	# HTTP_GET
	data_map = {
		"app_name": "gf_eth_monitor"
	}
	r = requests.post(url_str, data=json.dumps(data_map))
	print(r.text)

	if not r.status_code == 200:
		print("notify_completion http request failed")
		exit(1)

#--------------------------------------------------
def build_containers(p_cont_image_name_str,
	p_dockerfile_path_str,
	p_docker_sudo_bool=False):
	
	docker_context_dir_str = f"{modd_str}/../../.."

	print("BUILDING CONTAINER -----------=========================")
	print(f"container image name - {p_cont_image_name_str}")
	print(f"dockerfile           - {p_dockerfile_path_str}")
	
	assert os.path.isfile(p_dockerfile_path_str)

	c_lst = []
	if p_docker_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		"docker build",
		f"-f {p_dockerfile_path_str}",
		f"--tag={p_cont_image_name_str}",
		docker_context_dir_str
	])

	c_str = " ".join(c_lst)
	print(c_str)

	_, _, exit_code_int = gf_core_cli.run(c_str)

	if not exit_code_int == 0:
		exit(1)

#--------------------------------------------------
def publish_containers(p_cont_image_name_str,
	p_docker_user_str,
	p_docker_pass_str,
	p_docker_sudo_bool=False):

	print("PUBLISHING CONTAINER -----------=========================")
	print(f"container image name - {p_cont_image_name_str}")

	# LOGIN
	docker_login(p_docker_user_str,
		p_docker_pass_str,
		p_docker_sudo_bool = p_docker_sudo_bool)

	#------------------------
	c_lst = []
	if p_docker_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		f"docker push {p_cont_image_name_str}"
	])

	c_str = " ".join(c_lst)
	print(c_str)

	_, _, exit_code_int = gf_core_cli.run(c_str)

	if not exit_code_int == 0:
		exit(1)

	#------------------------

#-------------------------------------------------------------
# DOCKER_LOGIN
def docker_login(p_docker_user_str,
	p_docker_pass_str,
	p_docker_sudo_bool = False):
	assert isinstance(p_docker_user_str, str)
	assert isinstance(p_docker_pass_str, str)

	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker", "login",
		"-u", p_docker_user_str,
		"--password-stdin"
	])
	print(" ".join(cmd_lst))

	p = subprocess.Popen(cmd_lst, stdin = subprocess.PIPE, stdout = subprocess.PIPE, stderr = subprocess.PIPE)
	p.stdin.write(bytes(p_docker_pass_str.encode("utf-8"))) # write password on stdin of "docker login" command
	
	stdout, stderr = p.communicate() # wait for command completion
	stdout_str = stdout.decode("ascii")
	stderr_str = stderr.decode("ascii")

	if not stdout_str == "":
		print(stdout_str)
	if not stderr_str == "":
		print(stderr_str)

	if not p.returncode == 0:
		exit(1)

	# ERROR
	if "Error" in stderr_str or "unauthorized" in stderr_str:
		print("failed to Docker login")
		exit(1)

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#-------------
	# RUN
	arg_parser.add_argument("-run", action = "store", default = "build_containers",
		help = '''
- '''+fg('yellow')+'build_containers'+attr(0)+'''   - build app Docker containers
- '''+fg('yellow')+'publish_containers'+attr(0)+''' - publish app Docker containers
- '''+fg('yellow')+'notify_completion'+attr(0)+'''  - notify remote HTTP endpoint of build completion
		''')

	#----------------------------
	# RUN_WITH_SUDO - boolean flag
	# in the default Docker setup the daemon is run as root and so docker client commands have to be run with "sudo".
	# newer versions of Docker allow for non-root users to run Docker daemons. 
	# also CI systems might run this command in containers as root-level users in which case "sudo" must not be specified.
	arg_parser.add_argument("-docker_sudo", action = "store_true", default=False,
		help = "specify if certain Docker CLI commands are to run with 'sudo'")

	#----------------------------
	# STATIC - boolean flag
	arg_parser.add_argument("-static", action = "store_true", default=False,
		help = "compile binaries with static linking")

	#-------------
	# ENV_VARS
	drone_commit_sha_str         = os.environ.get("DRONE_COMMIT_SHA", None) # Drone defined ENV var
	gf_docker_user_str           = os.environ.get("GF_DOCKER_USER", None)
	gf_docker_pass_str           = os.environ.get("GF_DOCKER_PASS", None)
	gf_notify_completion_url_str = os.environ.get("GF_NOTIFY_COMPLETION_URL", None)

	#-------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)

	return {
		"run":                      args_namespace.run,
		"drone_commit_sha":         drone_commit_sha_str,
		"gf_docker_user_str":       gf_docker_user_str,
		"gf_docker_pass_str":       gf_docker_pass_str,
		"gf_notify_completion_url_str": gf_notify_completion_url_str,
		"docker_sudo_bool":             args_namespace.docker_sudo,
		"static_bool":                  args_namespace.static
	}

#--------------------------------------------------
main()