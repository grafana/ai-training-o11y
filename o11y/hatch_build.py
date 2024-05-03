from hatchling.builders.hooks.plugin.interface import BuildHookInterface
import pathlib
import shutil
import subprocess
import os

class CustomBuildHook(BuildHookInterface):
    def initialize(self, version: str, build_data: dict):
        build_data['pure_python'] = False
        build_data['artifacts'].extend(["src/o11y"])
        go_bin = self.build_go()
        build_data['force_include'][go_bin] = "/go-plugin/go-plugin"
        # TODO: Set platform version in the wheel
        # Make it possible to do platform builds based on an env variable

    def build_go(self):
        go_bin = shutil.which('go')
        if not go_bin:
            self.app.abort('Go is not installed, go is required to build from source')

        subprocess.check_call(
            [
                go_bin,
                'build',
                '-o',
                './dist/go-plugin',
                '-v',
            ],
        cwd='./src/go-plugin/'
        )
        
        return pathlib.Path("src", "go-plugin", "dist", "go-plugin").as_posix()
