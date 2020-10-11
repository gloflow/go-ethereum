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

import os
import subprocess

#---------------------------------------------------
# RUN
def run(p_cmd_str,
	p_env_map = {}):

	# env map has to contains all the parents ENV vars as well
	p_env_map.update(os.environ)

	p = subprocess.Popen(p_cmd_str,
		env     = p_env_map,
		shell   = True,
		stdout  = subprocess.PIPE,
		stderr  = subprocess.PIPE,
		bufsize = 1)

	for line in iter(p.stdout.readline, b''):	
		line_str = line.strip().decode("utf-8")
		print(line_str)

	for line in iter(p.stderr.readline, b''):	
		line_str = line.strip().decode("utf-8")
		print(line_str)

	p.communicate()
	
	return "", "", p.returncode