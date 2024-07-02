from hatchling.builders.hooks.plugin.interface import BuildHookInterface
import pathlib
import shutil
import subprocess
import re
import sysconfig
import os

class CustomBuildHook(BuildHookInterface):
    def initialize(self, version: str, build_data: dict):
        build_data['pure_python'] = False
        build_data['artifacts'].extend(["src/o11y"])

        # Check env variable to see what platform we want to build for
        target_os = os.environ['TARGET_OS']
        target_arch = os.environ['TARGET_ARCH']
        if target_os is None:
            build_data['tag'] = f"py3-none-{self._get_platform_tag()}"
            go_os_flag = None
            go_arch_flag = None
        elif target_os == 'windows':
            if target_arch == 'arm64':
                build_data['tag'] = f"py3-none-win_aarch64"
                go_arch_flag = "arm64"
            else:
                build_data['tag'] = f"py3-none-win_amd64"
                go_arch_flag = "amd64"
            go_os_flag = "windows"
        elif target_os == 'mac':
            if target_arch == 'arm64':
                build_data['tag'] = f"py3-none-macosx_13_0_aarch64"
                go_arch_flag = "arm64"
            else:
                build_data['tag'] = f"py3-none-macosx_13_0_x86_64"
                go_arch_flag = "amd64"
            go_os_flag = "darwin"
        elif target_os == 'linux':
            if target_arch == 'arm64':
                build_data['tag'] = f"py3-none-linux_aarch64"
                go_arch_flag = "arm64"
            else:
                build_data['tag'] = f"py3-none-linux_x86_64"
                go_arch_flag = "amd64"
            go_os_flag = "linux"

        go_bin = self.build_go(go_os_flag, go_arch_flag)
        build_data['force_include'][go_bin] = "/o11y/_internal/o11y-go"

    # Substantially identical to the same function from wandb's exporter
    # It turns out this is a pretty tricky problem
    # Used because MIT license allows it
    def _get_platform_tag(self) -> str:
        """Returns the platform tag for the current platform."""
        # Replace dots, spaces and dashes with underscores following
        # https://packaging.python.org/en/latest/specifications/platform-compatibility-tags/#platform-tag
        platform_tag = re.sub("[-. ]", "_", sysconfig.get_platform())

        # On macOS versions >=11, pip expects the minor version to be 0:
        #   https://github.com/pypa/packaging/issues/435
        #
        # You can see the list of tags that pip would support on your machine
        # using `pip debug --verbose`. On my macOS, get_platform() returns
        # 14.1, but `pip debug --verbose` reports only these py3 tags with 14:
        #
        # * py3-none-macosx_14_0_arm64
        # * py3-none-macosx_14_0_universal2
        #
        # We do this remapping here because otherwise, it's possible for `pip wheel`
        # to successfully produce a wheel that you then cannot `pip install` on the
        # same machine.
        macos_match = re.fullmatch(r"macosx_(\d+_\d+)_(\w+)", platform_tag)
        if macos_match:
            major, _ = macos_match.group(1).split("_")
            if int(major) >= 11:
                arch = macos_match.group(2)
                platform_tag = f"macosx_{major}_0_{arch}"

        return platform_tag

    def build_go(self, go_os_flag, go_arch_flag):
        go_bin = shutil.which('go')
        if not go_bin:
            self.app.abort('Go is not installed, go is required to build from source')

        env = {}
        if go_os_flag and go_arch_flag:
            env = {
                'GOOS': go_os_flag,
                'GOARCH': go_arch_flag,
            }
            # Set go cache directory
            go_cache_dir = os.path.join(os.getcwd(), '.cache', 'go-build')
            env['GOCACHE'] = go_cache_dir
            # make the go cache directory if it does not exist
            if not os.path.exists(go_cache_dir):
                os.makedirs(go_cache_dir)
            # Set go mod cache directory
            go_mod_cache_dir = os.path.join(os.getcwd(), '.cache', 'go-mod')
            # make the go mod cache directory if it does not exist
            if not os.path.exists(go_mod_cache_dir):
                os.makedirs(go_mod_cache_dir)
            env['GOMODCACHE'] = os.path.join(os.getcwd(), '.cache', 'go-mod')
            # Set gopath
            go_path_dir = os.path.join(os.getcwd(), '.cache', 'go-path')
            if not os.path.exists(go_path_dir):
                os.makedirs(go_path_dir)
            env['GOPATH'] = go_path_dir


        go_build_cmd = [
            go_bin,
            'build',
            '-o',
            './dist/o11y-go',
            '-v',
        ]

        try:
            result = subprocess.run(
                go_build_cmd,
                cwd='./src/o11y-go/',
                check=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                env=env
            )
            result.check_returncode()
        except subprocess.CalledProcessError as e:
            error_message = f"Error building Go plugin:\n{e.stderr}"
            raise RuntimeError(error_message) from e
        
        return pathlib.Path("src", "o11y-go", "dist", "o11y-go").as_posix()
