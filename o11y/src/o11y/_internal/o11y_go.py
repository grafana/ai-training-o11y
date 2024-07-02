import pkg_resources
import subprocess

resouce_package = __name__
resource_path = 'o11y-go'

def run():
    # run the binary from resource_path
    subprocess.run([pkg_resources.resource_filename(resouce_package, resource_path)])
