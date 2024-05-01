from setuptools import setup, find_packages

setup(
    name="o11y",
    version="0.1.0",
    packages=find_packages(),
    install_requires=["requests", "pydantic", "python-logging-loki"],
    extras_require={"dev": ["pytest"]}
)
