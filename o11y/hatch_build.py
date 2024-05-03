from hatchling.builders.hooks.plugin.interface import BuildHookInterface
import shutil
import subprocess
import os

class CustomBuildHook(BuildHookInterface):
    def initialize(self, version: str, build_data: dict):
        build_data['pure_python'] = False
        self.build_go()
    
    def build_go(self):
        go_bin = shutil.which('go')
        if not go_bin:
            self.app.abort('Go is not installed, go is required to build from source')
        
        # self.print_file_tree(os.getcwd())

        subprocess.check_call(
            [
                go_bin,
                'build',
                '-o',
                './dist/go-plugin',
            ],
        cwd='./go-plugin/'
        )
