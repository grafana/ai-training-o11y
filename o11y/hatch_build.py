from hatchling.builders.hooks.plugin.interface import BuildHookInterface
import pathlib
import shutil
import subprocess
import re
import sysconfig

class CustomBuildHook(BuildHookInterface):
    def initialize(self, version: str, build_data: dict):
        build_data['pure_python'] = False
        build_data['artifacts'].extend(["src/o11y"])
        go_bin = self.build_go()
        build_data['force_include'][go_bin] = "/o11y/go-plugin"
        build_data['tag'] = f"py3-none-{self._get_platform_tag()}"
        # TODO: Make it possible to do platform builds based on an env variable

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
