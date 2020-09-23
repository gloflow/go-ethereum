# Copyright 2019 The go-ethereum Authors
# This file is part of the go-ethereum library.
#
# The go-ethereum library is free software: you can redistribute it and/or modify
# it under the terms of the GNU Lesser General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# The go-ethereum library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Lesser General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public License
# along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

import os, sys
modd_str = os.path.abspath(os.path.dirname(__file__)) # module dir

import argparse
import subprocess

from colored import fg, bg, attr
import delegator

sys.path.append("%s/../utils"%(modd_str))
import gf_core_cli

#--------------------------------------------------
def main():

	args_map = parse_args()

	service_name_str            = "gf_go-ethereum"
	service_bin_path_str        = f"{modd_str}/../../../build/bin/geth"
	service_dir_path_str        = "%s/../../../"%(modd_str)
	service_cont_image_tag_str  = "latest"
	service_cont_image_name_str = f"glofloworg/gf_go-ethereum:{service_cont_image_tag_str}"
	service_cont_dockerfile_path_str = f"{modd_str}/../../../Dockerfile"
	docker_user_str                  = "glofloworg"
	
	#------------------------
	# BUILD
	if args_map["run"] == "build":

		build_go(service_name_str,
			service_bin_path_str,
			service_dir_path_str,
			p_static_bool = True)

	#------------------------
	# BUILD_CONTAINER
	elif args_map["run"] == "build_containers":
		
		build_containers(service_cont_image_name_str,
			service_cont_dockerfile_path_str,
			service_dir_path_str,
			p_docker_sudo_bool=True)

	#------------------------
	# PUBLISH_CONTAINER
	elif args_map["run"] == "publish_containers":
		docker_pass_str = args_map["gf_docker_pass_str"]
		assert not docker_pass_str == None

		publish_containers(service_cont_image_name_str,
			docker_user_str,
			docker_pass_str,
			p_docker_sudo_bool=True)

	#------------------------

#--------------------------------------------------
# BUILD_GO
def build_go(p_name_str,
	p_go_bin_path_str,
	p_go_dir_path_str,
	p_static_bool       = False,
	p_exit_on_fail_bool = True):
	assert isinstance(p_static_bool, bool)
	assert os.path.isdir(p_go_dir_path_str)

	print("")
	if p_static_bool:
		print(" -- %sSTATIC BINARY BUILD%s"%(fg("yellow"), attr(0)))
		
	print(" -- build %s%s%s service"%(fg("green"), p_name_str, attr(0)))
	print(" -- go_dir_path - %s%s%s"%(fg("green"), p_go_dir_path_str, attr(0)))

	cwd_str = os.getcwd()
	os.chdir(p_go_dir_path_str) # change into the target main package dir

	# GO_GET
	_, _, exit_code_int = gf_core_cli.run("go get -u")
	print("")
	print("")

	# STATIC_LINKING - when deploying to containers it is not always guaranteed that all
	#                  required libraries are present. so its safest to compile to a statically
	#                  linked lib.
	#                  build time a few times larger then regular, so slow for dev.
	if p_static_bool:
		
		args_lst = [
			"make geth"
		]
		c_str = " ".join(args_lst)
		print(c_str)

	# DYNAMIC_LINKING - fast build for dev.
	else:
		c_str = "make geth"

	# RUN
	_, _, exit_code_int = gf_core_cli.run(c_str)
	assert os.path.isfile(p_go_bin_path_str)

	# IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
	#               want to fail with a non-zero exit code as well - this way other CI 
	#               programs will flag builds as failed.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	os.chdir(cwd_str) # return to initial dir

#--------------------------------------------------
def build_containers(p_cont_image_name_str,
	p_dockerfile_path_str,
	p_docker_context_dir_str,
	p_docker_sudo_bool=False):
	assert os.path.isdir(p_docker_context_dir_str)

	print("BUILDING CONTAINER -----------=========================")
	print(f"container image name - {p_cont_image_name_str}")
	print(f"dockerfile          - {p_dockerfile_path_str}")
	
	assert os.path.isfile(p_dockerfile_path_str)

	c_lst = []
	if p_docker_sudo_bool:
		c_lst.append("sudo")

	c_lst.extend([
		"docker build",
		f"-f {p_dockerfile_path_str}",
		f"--tag={p_cont_image_name_str}",
		p_docker_context_dir_str
	])

	c_str = " ".join(c_lst)
	print(c_str)

	_, _, exit_code_int = gf_core_cli.run(c_str)

	if not exit_code_int == 0:
		exit()

#--------------------------------------------------
def publish_containers(p_cont_image_name_str,
	p_docker_user_str,
	p_docker_pass_str,
	p_docker_sudo_bool=False):
	print("BUILDING CONTAINER -----------=========================")
	print(f"container image name - {p_cont_image_name_str}")

	#------------------------
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

	if not p.returncode == 0:
		exit()

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

	print(bytes(p_docker_pass_str.encode("utf-8")))
	p = subprocess.Popen(cmd_lst, stdin = subprocess.PIPE, stdout = subprocess.PIPE, stderr = subprocess.PIPE)
	p.stdin.write(bytes(p_docker_pass_str.encode("utf-8"))) # write password on stdin of "docker login" command
	
	stdout_str, stderr_str = p.communicate() # wait for command completion
	print(stdout_str)
	print(stderr_str)

	if not p.returncode == 0:
		exit()

#--------------------------------------------------
def parse_args():
	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)

	#-------------
	# RUN
	arg_parser.add_argument("-run", action = "store", default = "build",
		help = '''
- '''+fg('yellow')+'test'+attr(0)+'''               - run app code tests
- '''+fg('yellow')+'build'+attr(0)+'''              - build app golang/web code
- '''+fg('yellow')+'build_containers'+attr(0)+'''   - build app Docker containers
- '''+fg('yellow')+'publish_containers'+attr(0)+''' - publish app Docker containers
		''')

	#----------------------------
	# RUN_WITH_SUDO - boolean flag
	# in the default Docker setup the daemon is run as root and so docker client commands have to be run with "sudo".
	# newer versions of Docker allow for non-root users to run Docker daemons. 
	# also CI systems might run this command in containers as root-level users in which case "sudo" must not be specified.
	arg_parser.add_argument("-docker_sudo", action = "store_true",
		help = "specify if certain Docker CLI commands are to run with 'sudo'")

	#-------------
	# ENV_VARS
	drone_commit_sha_str         = os.environ.get("DRONE_COMMIT_SHA", None) # Drone defined ENV var
	gf_docker_user_str           = os.environ.get("GF_DOCKER_USER", None)
	gf_docker_pass_str           = os.environ.get("GF_DOCKER_P", None)
	gf_notify_completion_url_str = os.environ.get("GF_NOTIFY_COMPLETION_URL", None)

	#-------------
	cli_args_lst   = sys.argv[1:]
	args_namespace = arg_parser.parse_args(cli_args_lst)
	return {
		"run":                      args_namespace.run,
		"drone_commit_sha":         drone_commit_sha_str,
		"gf_docker_user_str":       gf_docker_user_str,
		"gf_docker_pass_str":       gf_docker_pass_str,
		"gf_notify_completion_url": gf_notify_completion_url_str,
		"docker_sudo":              args_namespace.docker_sudo
	}

#--------------------------------------------------
main()