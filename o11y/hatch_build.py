from hatchling.builders.hooks.plugin.interface import BuildHookInterface
import pathlib
import shutil
import subprocess
import os

# function to print pwd recursively
def list_files_recursive(directory):
    file_list = []
    # List all files and directories in the current directory
    for item in os.listdir(directory):
        # Construct full path
        full_path = os.path.join(directory, item)
        # If it's a file, add it to the list
        if os.path.isfile(full_path):
            file_list.append(full_path)
        # If it's a directory, recursively call this function
        elif os.path.isdir(full_path):
            file_list.extend(list_files_recursive(full_path))
    return file_list


class CustomBuildHook(BuildHookInterface):
    def initialize(self, version: str, build_data: dict):
        print("Buildhookran")
        build_data['pure_python'] = False
        build_data['artifacts'].extend(["src/o11y"])
        print("Buildhookrun2")
        cwd = os.getcwd()
        _ = [print(i) for i in list_files_recursive(cwd)]
        go_bin = self.build_go()
        build_data['force_include'][go_bin] = "/go-plugin/go-plugin"
    
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
