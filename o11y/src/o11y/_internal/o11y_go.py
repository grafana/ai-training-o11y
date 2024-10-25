# SPDX-License-Identifier: Apache-2.0

import subprocess
import os
import pkg_resources
import pathlib


def run():
    if os.environ.get('HATCH_ENV_ACTIVE') or os.environ.get('HATCH'):
        current_dir = pathlib.Path(__file__).parent
        # Go up one level and then to the binary
        binary_path = str(current_dir.parent.parent / "o11y-go" / "dist" / "o11y-go")
    else:
        binary_path = pkg_resources.resource_filename(__name__, 'o11y-go')
    # Duplicate the real stdout/stderr file descriptors
    stdout_copy = os.dup(1)  # 1 is stdout
    stderr_copy = os.dup(2)  # 2 is stderr
    
    # Launch the binary with a pipe for stdin
    process = subprocess.Popen(
        [binary_path],
        stdin=subprocess.PIPE,
        pass_fds=(stdout_copy, stderr_copy)  # Pass our duplicated fds to the child
    )
    
    # Send the descriptor numbers to the Go process
    process.stdin.write(f"{stdout_copy} {stderr_copy}\n".encode())
    process.stdin.flush()
    
    return process
